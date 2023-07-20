package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type FriggConfig struct {
	Healthchecks []Healthcheck `yaml:"healthchecks"`
}

type Healthcheck struct {
	Name     string `yaml:"name"`
	Interval int    `yaml:"interval"`
	Url      string `yaml:"url"`
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

	if resp.StatusCode == 200 {
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
	return friggConfig
}
