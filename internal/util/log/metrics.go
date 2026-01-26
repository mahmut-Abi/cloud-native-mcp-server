package log

import (
	"sync"
	"time"
)

type Metrics struct {
	mu       sync.RWMutex
	requests map[string]*RequestMetric
}

type RequestMetric struct {
	Count      int64
	Total      time.Duration
	Min        time.Duration
	Max        time.Duration
	LastError  error
	ErrorCount int64
}

var metrics = &Metrics{
	requests: make(map[string]*RequestMetric),
}

func RecordRequest(name string, duration time.Duration, err error) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	m := metrics.requests[name]
	if m == nil {
		m = &RequestMetric{
			Min: duration,
			Max: duration,
		}
		metrics.requests[name] = m
	}

	m.Count++
	m.Total += duration
	if duration < m.Min {
		m.Min = duration
	}
	if duration > m.Max {
		m.Max = duration
	}

	if err != nil {
		m.ErrorCount++
		m.LastError = err
	}
}

func GetMetrics() map[string]*RequestMetric {
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	result := make(map[string]*RequestMetric)
	for k, v := range metrics.requests {
		result[k] = v
	}
	return result
}
