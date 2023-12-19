# frigg

A proof of concept (POC) of a simple healthcheck tool written in Go, that can be self hosted on your own infrastructure.
Healthcheck routes and assertions are stored as code (config file), rather than added in some SaaS UI. 
The idea would be to programmatically generate this config (or provide an SDK), to apply more robust health checks
across different geographic locations, and as new use cases emerge.
Once again, this is just a POC, and I haven't touched it since an initial urge to hack it together.

# Running

Example 1: `go run main.go example1.yml`
Output:
```bash
➜  frigg git:(main) ✗ go run main.go example1.yml
google - 5
github - 10
gitlab - 20
2023/12/19 14:01:10 github : 10 : https://www.github.com/ pass
2023/12/19 14:01:10 google : 5 : https://www.google.com/ pass
2023/12/19 14:01:11 gitlab : 20 : https://www.gitlab.com/ pass
2023/12/19 14:01:15 google : 5 : https://www.google.com/ pass
2023/12/19 14:01:20 github : 10 : https://www.github.com/ pass
2023/12/19 14:01:20 google : 5 : https://www.google.com/ pass
2023/12/19 14:01:25 google : 5 : https://www.google.com/ pass
2023/12/19 14:01:30 github : 10 : https://www.github.com/ pass
2023/12/19 14:01:30 google : 5 : https://www.google.com/ pass
2023/12/19 14:01:31 gitlab : 20 : https://www.gitlab.com/ pass
```


Example 2: `go run main.go example2.yml`
Output:
```bash
➜  frigg git:(main) ✗ go run main.go example2.yml
google - 5
2023/12/19 14:02:06 google : 5 : https://www.google.com/ pass
2023/12/19 14:02:11 google : 5 : https://www.google.com/ pass
2023/12/19 14:02:16 google : 5 : https://www.google.com/ pass
2023/12/19 14:02:21 google : 5 : https://www.google.com/ pass
```
