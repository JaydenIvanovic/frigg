package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type FriggConfig struct {
	Healthchecks []Healthcheck `yaml:"healthchecks"`
}

type Healthcheck struct {
	Name          string   `yaml:"name"`
	Interval      int      `yaml:"interval"`
	Url           string   `yaml:"url"`
	RawAssertions []string `yaml:"assertions"`
	Asserters     []Asserter
}

func (h Healthcheck) DebugInfo() string {
	return fmt.Sprintf("%s : %d : %s", h.Name, h.Interval, h.Url)
}

func (h Healthcheck) PrintDebugInfo() {
	fmt.Printf("%s - %d \n", h.Name, h.Interval)
}

func (h Healthcheck) Do() {
	resp, err := http.Get(h.Url)
	if err != nil {
		log.Printf("request for %s failed with: %v", h.Url, err)
	}

	pass := true
	for _, a := range h.Asserters {
		pass = pass && a.Do(resp)
	}

	if pass {
		log.Printf("%s pass", h.DebugInfo())
	} else {
		log.Printf("%s fail", h.DebugInfo())
	}
}

func main() {
	var wg sync.WaitGroup

	rawData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	friggConfig := parseConfig(rawData)
	for _, h := range friggConfig.Healthchecks {
		h.PrintDebugInfo()

		// Always spawn a goroutine for each healthcheck and
		// perform them straight away on startup. That way new
		// deploys won't cause a delay / gap between intervals (missed checks)
		go func(h Healthcheck) {
			h.Do()
		}(h)

		wg.Add(1)
		go func(h Healthcheck) {
			ticker := time.NewTicker(time.Duration(h.Interval) * time.Second)
			for range ticker.C {
				h.Do()
			}
		}(h)
	}

	// Waits indefinitely...
	wg.Wait()
}

func parseConfig(config []byte) FriggConfig {
	var friggConfig FriggConfig
	err := yaml.Unmarshal(config, &friggConfig)
	if err != nil {
		panic(err)
	}

	for i, h := range friggConfig.Healthchecks {
		for _, a := range h.RawAssertions {
			h.Asserters = append(h.Asserters, NewAsserter(a))
		}
		friggConfig.Healthchecks[i].Asserters = h.Asserters
	}

	return friggConfig
}

// Assertion Parser and Logic
type AssertionParser struct {
	raw string
}

func NewAsserter(dsl string) Asserter {
	a := AssertionParser{raw: dsl}
	return a.ParseDsl()
}

func (a AssertionParser) ParseDsl() Asserter {
	parts := strings.Split(a.raw, "(")
	fn := parts[0]
	val := strings.Split(parts[1], ")")[0]

	switch fn {
	case "text":
		return TextAsserter{
			Value: val,
		}
	case "status_code":
		valAsInt, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		return StatusCodeAsserter{
			Value: valAsInt,
		}
	}

	panic("should be unreachable")
}

type Asserter interface {
	Do(resp *http.Response) bool
}

type TextAsserter struct {
	Value string
}

func (a TextAsserter) Do(resp *http.Response) bool {
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	body := string(rawBody)
	return strings.Contains(body, a.Value)
}

type StatusCodeAsserter struct {
	Value int
}

func (a StatusCodeAsserter) Do(resp *http.Response) bool {
	return resp.StatusCode == a.Value
}
