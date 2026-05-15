package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func checkEndpoint(client *http.Client, baseURL, endpoint string) error {
	url := strings.TrimRight(baseURL, "/") + endpoint
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

func runHealthCheck(baseURL string, endpoints []string, client *http.Client) []string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var failures []string

	for _, ep := range endpoints {
		ep = strings.TrimSpace(ep)
		if ep == "" {
			continue
		}
		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()
			err := checkEndpoint(client, baseURL, endpoint)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				failure := fmt.Sprintf("FAIL %s: %v", endpoint, err)
				failures = append(failures, failure)
				fmt.Println(failure)
			} else {
				fmt.Printf("PASS %s: HTTP 200\n", endpoint)
			}
		}(ep)
	}

	wg.Wait()
	return failures
}

func writeSummary(baseURL string, failures []string, isRollback bool) {
	summaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
	if summaryFile == "" {
		return
	}
	f, err := os.OpenFile(summaryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to GITHUB_STEP_SUMMARY: %v\n", err)
		return
	}
	defer f.Close()

	title := "## Post-deploy Health Verification"
	if isRollback {
		title = "## Post-rollback Health Verification"
	}

	fmt.Fprintln(f, title)
	fmt.Fprintf(f, "Base URL: %s\n\n", baseURL)

	if len(failures) == 0 {
		fmt.Fprintln(f, "- ✅ All deterministic checks passed.")
	} else {
		fmt.Fprintln(f, "- ❌ One or more checks failed:")
		for _, failure := range failures {
			fmt.Fprintf(f, "  - %s\n", failure)
		}
	}
}

func main() {
	var baseURL string
	var endpointsArg string
	var timeout int
	var retries int
	var retryDelay int
	var isRollback bool

	flag.StringVar(&baseURL, "base-url", "", "Base URL to test")
	flag.StringVar(&endpointsArg, "endpoints", "/", "Comma-separated list of endpoints")
	flag.IntVar(&timeout, "timeout", 10, "HTTP client timeout in seconds")
	flag.IntVar(&retries, "retries", 3, "Number of retries for failures")
	flag.IntVar(&retryDelay, "retry-delay", 5, "Seconds to wait between retries")
	flag.BoolVar(&isRollback, "rollback", false, "Set to true if this is a rollback healthcheck")
	flag.Parse()

	if baseURL == "" {
		fmt.Fprintln(os.Stderr, "Error: --base-url is required")
		os.Exit(1)
	}

	endpoints := strings.Split(endpointsArg, ",")
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}

	var failures []string
	for attempt := 1; attempt <= retries; attempt++ {
		fmt.Printf("Health check attempt %d of %d...\n", attempt, retries)
		failures = runHealthCheck(baseURL, endpoints, client)
		if len(failures) == 0 {
			fmt.Println("All health checks passed.")
			writeSummary(baseURL, nil, isRollback)
			os.Exit(0)
		}
		if attempt < retries {
			fmt.Printf("Waiting %d seconds before retrying...\n", retryDelay)
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}
	}

	fmt.Fprintln(os.Stderr, "Health check failed after all retries.")
	writeSummary(baseURL, failures, isRollback)
	os.Exit(1)
}
