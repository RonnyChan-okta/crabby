package main

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"
)

// startInternalMetrics collects metrics from Crabby's Go runtime and reports them as metrics
func startInternalMetrics(ctx context.Context, wg *sync.WaitGroup, storage *Storage, interval uint) {
	var memstats runtime.MemStats
	var metricsTicker *time.Ticker

	if interval > 0 {
		metricsTicker = time.NewTicker(time.Duration(interval) * time.Second)
	} else {
		metricsTicker = time.NewTicker(15 * time.Second)
	}

	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-metricsTicker.C:
			runtime.ReadMemStats(&memstats)
			storage.MetricDistributor <- makeInternalMetric("mem.alloc", float64(memstats.Alloc))
			storage.MetricDistributor <- makeInternalMetric("heap.alloc", float64(memstats.HeapAlloc))
			storage.MetricDistributor <- makeInternalMetric("heap.in_use", float64(memstats.HeapInuse))
			storage.MetricDistributor <- makeInternalMetric("num_goroutines", float64(runtime.NumGoroutine()))

		case <-ctx.Done():
			log.Println("Cancellation request received.  Cancelling job runner.")
			return
		}
	}

}

func makeInternalMetric(name string, value float64) Metric {
	m := Metric{
		Job:       "internal_metrics",
		Timing:    name,
		Value:     value,
		Timestamp: time.Now(),
	}
	return m
}
