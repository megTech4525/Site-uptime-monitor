package main

import (
	"context"
	"net/http"
	"sync"
	"time"
)

func check(ctx context.Context, client *http.Client, url string) Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{URL: url, Err: err}
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{URL: url, Err: err, Latency: time.Since(start)}

	}
	defer resp.Body.Close()
	return Result{URL: url, Status: resp.StatusCode, Latency: time.Since(start)}
}

func worker(ctx context.Context, client *http.Client, jobs <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range jobs {
		results <- check(ctx, client, url)
	}
}
