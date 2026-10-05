package main

import "time"

type Result struct {
	URL     string
	Status  int
	Latency time.Duration
	Err     error
}
