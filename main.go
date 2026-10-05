package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	a := gin.Default()
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://www.medium.com",
		"https://www.linkedin.com",
		"https://www.twitter.com",
		"https://www.facebook.com",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := &http.Client{Timeout: 5 * time.Second}

	jobs := make(chan string)
	results := make(chan Result)

	var wg sync.WaitGroup
	for range urls {
		wg.Add(1)
		go worker(ctx, client, jobs, results, &wg)
	}
	go func() {
		for _, u := range urls {
			jobs <- u
		}
	}()
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.Err != nil {
			fmt.Printf("DOWN %s (%v)\n", r.URL, r.Err)
			continue
		}
		fmt.Printf("%d %s (%v)\n", r.Status, r.URL, r.Latency)
	}

	a.Run()
}
