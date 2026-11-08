
package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CheckResult struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Error      error
}

type URLChecker interface {
	Check(url string) CheckResult
}

type HTTPClient struct {
	Client *http.Client
}

type Stats struct {
	Total       int
	Success200  int
	Errors      int
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (h *HTTPClient) Check(url string) CheckResult {
	result := CheckResult{
		URL: url,
	}

	start := time.Now()

	response, err := h.Client.Get(url)

	result.Duration = time.Since(start)

	if err != nil {
		result.Error = err
		return result
	}

	defer response.Body.Close()

	result.StatusCode = response.StatusCode

	return result
}

func main() {
	client := NewHTTPClient()

	var checker URLChecker = client

	urls := []string{
		"https://go.dev",
		"https://github.com",
		"https://google.com",
		"https://example.com",
		"https://this-domain-does-not-exist-12345.com",
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	stats := Stats{}

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			result := checker.Check(url)

			fmt.Printf(
				"URL: %s | Status: %d | Duration: %v | Error: %v\n",
				result.URL,
				result.StatusCode,
				result.Duration,
				result.Error,
			)

			mutex.Lock()

			stats.Total++

			if result.Error != nil {
				stats.Errors++
			} else if result.StatusCode == http.StatusOK {
				stats.Success200++
			}

			mutex.Unlock()
		}(url)
	}

	wg.Wait()

	fmt.Println("\nSTATISTICS")
	fmt.Println("Total checks:", stats.Total)
	fmt.Println("Successful (200):", stats.Success200)
	fmt.Println("Errors:", stats.Errors)
}
