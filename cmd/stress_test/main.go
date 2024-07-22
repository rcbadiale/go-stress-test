package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"flag"

	"github.com/rcbadiale/go-stress-test/internal/manager"
	"github.com/rcbadiale/go-stress-test/internal/tasks"
)

type Config struct {
	Concurrency int
	Requests    int
	Timeout     int
	Method      string
	URL         string
}

func ParseConfig() *Config {
	concurrency := flag.Int("concurrency", 3, "number of workers")
	requests := flag.Int("requests", 10, "number of requests")
	timeout := flag.Int("timeout", 5, "timeout in seconds")
	method := flag.String("method", "GET", "HTTP method")
	url := flag.String("url", "", "URL")
	flag.Parse()

	if *url == "" {
		log.Fatal("error: URL is required")
	}

	return &Config{
		Concurrency: *concurrency,
		Requests:    *requests,
		Timeout:     *timeout,
		Method:      *method,
		URL:         *url,
	}
}

func main() {
	config := ParseConfig()

	fmt.Println("Configuration:")
	fmt.Printf("- Concurrency: %d\n", config.Concurrency)
	fmt.Printf("- Requests: %d\n", config.Requests)
	fmt.Printf("- Timeout: %ds\n", config.Timeout)
	fmt.Printf("- Method: %s\n", config.Method)
	fmt.Printf("- URL: %s\n\n", config.URL)

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	client.Timeout = time.Duration(config.Timeout) * time.Second

	tasksList := []tasks.Task{}
	for i := range config.Requests {
		task, err := tasks.NewRequestTask(i, client, config.Method, config.URL, nil)
		if err != nil {
			log.Fatalf("error: %s", err)
		}
		tasksList = append(tasksList, task)
	}

	manager := manager.NewManager(config.Concurrency, tasksList)
	duration := manager.Execute()
	fmt.Println(manager.Report())
	fmt.Printf("\nTotal time: %.2fs\n", duration.Seconds())
}
