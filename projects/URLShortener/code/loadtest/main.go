package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	targetURL      = "http://localhost:8080/shorten"
	rampUpInterval = 5 * time.Second
	stepSize       = 50   // Increase concurrent workers by this amount
	maxErrorRate   = 0.05 // 5% error rate is considered "breaking"
)

func main() {
	fmt.Println("Starting load test...")
	fmt.Printf("Target: %s\n", targetURL)

	concurrency := 10
	for {
		fmt.Printf("\nTesting with %d concurrent workers...\n", concurrency)
		metrics := runLoad(concurrency, rampUpInterval)

		fmt.Printf("RPS: %.2f, Avg Latency: %.2fms, Error Rate: %.2f%%\n", metrics.RPS, metrics.AvgLatency, metrics.ErrorRate*100)

		if metrics.ErrorRate > maxErrorRate {
			fmt.Printf("Breaking point reached at ~%d concurrent user! Error rate: %.2f%%\n", concurrency, metrics.ErrorRate*100)
			break
		}

		concurrency += stepSize
	}
}

type Metrics struct {
	RPS        float64
	AvgLatency float64 // ms
	ErrorRate  float64
}

func runLoad(concurrency int, duration time.Duration) Metrics {
	var wg sync.WaitGroup
	var totalReqs int64
	var totalErrors int64
	var totalLatency int64 // microseconds

	stop := make(chan struct{})

	// Create workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{Timeout: 2 * time.Second}
			for {
				select {
				case <-stop:
					return
				default:
					start := time.Now()
					reqBody := `{"url": "http://example.com"}`
					resp, err := client.Post(targetURL, "application/json", bytes.NewBufferString(reqBody))
					latency := time.Since(start).Microseconds()

					atomic.AddInt64(&totalReqs, 1)
					atomic.AddInt64(&totalLatency, latency)

					if err != nil || resp.StatusCode >= 500 {
						atomic.AddInt64(&totalErrors, 1)
					}
					if resp != nil {
						resp.Body.Close()
					}
				}
			}
		}()
	}

	time.Sleep(duration)
	close(stop)
	wg.Wait()

	rps := float64(totalReqs) / duration.Seconds()
	avgLatency := float64(totalLatency) / float64(totalReqs) / 1000.0 // convert micro to milli
	errorRate := float64(totalErrors) / float64(totalReqs)

	if totalReqs == 0 {
		avgLatency = 0
		errorRate = 0
	}

	return Metrics{
		RPS:        rps,
		AvgLatency: avgLatency,
		ErrorRate:  errorRate,
	}
}
