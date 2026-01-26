package handlers

import (
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// GrafanaToolMetrics Grafana tool execution metrics
type GrafanaToolMetrics struct {
	ToolName         string
	CallCount        int64
	TotalDuration    time.Duration
	MinDuration      time.Duration
	MaxDuration      time.Duration
	LastCalled       time.Time
	ErrorCount       int64
	SuccessCount     int64
	AverageSize      int64 // Average response size
	LastResponseSize int64
}

// GrafanaMetricsCollector Grafana metrics collector
type GrafanaMetricsCollector struct {
	metrics map[string]*GrafanaToolMetrics
	mutex   sync.RWMutex
}

// NewGrafanaMetricsCollector creates new Grafana metrics collector
func NewGrafanaMetricsCollector() *GrafanaMetricsCollector {
	return &GrafanaMetricsCollector{
		metrics: make(map[string]*GrafanaToolMetrics),
	}
}

// RecordGrafanaToolCall records Grafana tool call
func (collector *GrafanaMetricsCollector) RecordGrafanaToolCall(toolName string, duration time.Duration, responseSize int64, isError bool) {
	collector.mutex.Lock()
	defer collector.mutex.Unlock()

	metrics, exists := collector.metrics[toolName]
	if !exists {
		metrics = &GrafanaToolMetrics{
			ToolName:    toolName,
			MinDuration: duration,
			MaxDuration: duration,
		}
		collector.metrics[toolName] = metrics
	}

	metrics.CallCount++
	metrics.TotalDuration += duration
	metrics.LastCalled = time.Now()
	metrics.LastResponseSize = responseSize

	if isError {
		metrics.ErrorCount++
	} else {
		metrics.SuccessCount++
	}

	// Update min/max duration
	if duration < metrics.MinDuration || metrics.MinDuration == 0 {
		metrics.MinDuration = duration
	}
	if duration > metrics.MaxDuration {
		metrics.MaxDuration = duration
	}

	// Update average size
	if metrics.CallCount > 0 {
		totalSize := metrics.AverageSize*(metrics.CallCount-1)/metrics.CallCount + responseSize/metrics.CallCount
		metrics.AverageSize = totalSize
	}

	logrus.WithFields(logrus.Fields{
		"toolName":     toolName,
		"duration":     duration,
		"responseSize": responseSize,
		"isError":      isError,
		"totalCalls":   metrics.CallCount,
	}).Debug("Grafana tool call metrics recorded")
}

// GetGrafanaMetrics gets Grafana tool metrics
func (collector *GrafanaMetricsCollector) GetGrafanaMetrics(toolName string) *GrafanaToolMetrics {
	collector.mutex.RLock()
	defer collector.mutex.RUnlock()

	if metrics, exists := collector.metrics[toolName]; exists {
		// Return copy to avoid concurrency issues
		return &GrafanaToolMetrics{
			ToolName:         metrics.ToolName,
			CallCount:        metrics.CallCount,
			TotalDuration:    metrics.TotalDuration,
			MinDuration:      metrics.MinDuration,
			MaxDuration:      metrics.MaxDuration,
			LastCalled:       metrics.LastCalled,
			ErrorCount:       metrics.ErrorCount,
			SuccessCount:     metrics.SuccessCount,
			AverageSize:      metrics.AverageSize,
			LastResponseSize: metrics.LastResponseSize,
		}
	}
	return nil
}

// GetAllGrafanaMetrics gets all Grafana tool metrics
func (collector *GrafanaMetricsCollector) GetAllGrafanaMetrics() map[string]*GrafanaToolMetrics {
	collector.mutex.RLock()
	defer collector.mutex.RUnlock()

	result := make(map[string]*GrafanaToolMetrics)
	for name, metrics := range collector.metrics {
		result[name] = &GrafanaToolMetrics{
			ToolName:         metrics.ToolName,
			CallCount:        metrics.CallCount,
			TotalDuration:    metrics.TotalDuration,
			MinDuration:      metrics.MinDuration,
			MaxDuration:      metrics.MaxDuration,
			LastCalled:       metrics.LastCalled,
			ErrorCount:       metrics.ErrorCount,
			SuccessCount:     metrics.SuccessCount,
			AverageSize:      metrics.AverageSize,
			LastResponseSize: metrics.LastResponseSize,
		}
	}
	return result
}

// GetAverageDuration gets average execution time
func (metrics *GrafanaToolMetrics) GetAverageDuration() time.Duration {
	if metrics.CallCount == 0 {
		return 0
	}
	return metrics.TotalDuration / time.Duration(metrics.CallCount)
}

// GetSuccessRate gets success rate
func (metrics *GrafanaToolMetrics) GetSuccessRate() float64 {
	if metrics.CallCount == 0 {
		return 0
	}
	return float64(metrics.SuccessCount) / float64(metrics.CallCount)
}

// Global Grafana metrics collector
var DefaultGrafanaMetricsCollector = NewGrafanaMetricsCollector()

// MeasureGrafanaToolExecution decorator to measure Grafana tool execution time
func MeasureGrafanaToolExecution(toolName string, execFunc func() (*interface{}, error)) (*interface{}, error) {
	start := time.Now()

	result, err := execFunc()
	duration := time.Since(start)

	responseSize := int64(0)
	if result != nil {
		if data, err2 := json.Marshal(*result); err2 == nil {
			responseSize = int64(len(data))
		}
	}

	DefaultGrafanaMetricsCollector.RecordGrafanaToolCall(toolName, duration, responseSize, err != nil)

	logrus.WithFields(logrus.Fields{
		"toolName":     toolName,
		"duration":     duration,
		"responseSize": responseSize,
		"success":      err == nil,
	}).Debug("Grafana tool execution measured")

	return result, err
}

// GrafanaPerformanceReport Grafana performance report
type GrafanaPerformanceReport struct {
	TopUsedTools        []GrafanaToolSummary `json:"topUsedTools"`
	SlowestTools        []GrafanaToolSummary `json:"slowestTools"`
	LargestResponses    []GrafanaToolSummary `json:"largestResponses"`
	TotalCalls          int64                `json:"totalCalls"`
	TotalErrors         int64                `json:"totalErrors"`
	AverageResponseSize int64                `json:"averageResponseSize"`
}

// GrafanaToolSummary Grafana tool summary
type GrafanaToolSummary struct {
	ToolName    string        `json:"toolName"`
	CallCount   int64         `json:"callCount"`
	AverageTime time.Duration `json:"averageTime"`
	MinTime     time.Duration `json:"minTime"`
	MaxTime     time.Duration `json:"maxTime"`
	SuccessRate float64       `json:"successRate"`
	AverageSize int64         `json:"averageSize"`
	LastCalled  time.Time     `json:"lastCalled"`
}

// GenerateGrafanaPerformanceReport generates Grafana performance report
func GenerateGrafanaPerformanceReport() GrafanaPerformanceReport {
	allMetrics := DefaultGrafanaMetricsCollector.GetAllGrafanaMetrics()

	var totalCalls int64
	var totalErrors int64
	var totalSize int64

	var toolSummaries []GrafanaToolSummary

	for _, metrics := range allMetrics {
		totalCalls += metrics.CallCount
		totalErrors += metrics.ErrorCount
		totalSize += metrics.AverageSize

		summary := GrafanaToolSummary{
			ToolName:    metrics.ToolName,
			CallCount:   metrics.CallCount,
			AverageTime: metrics.GetAverageDuration(),
			MinTime:     metrics.MinDuration,
			MaxTime:     metrics.MaxDuration,
			SuccessRate: metrics.GetSuccessRate(),
			AverageSize: metrics.AverageSize,
			LastCalled:  metrics.LastCalled,
		}

		toolSummaries = append(toolSummaries, summary)
	}

	// Sort to get most used tools
	topUsed := make([]GrafanaToolSummary, len(toolSummaries))
	copy(topUsed, toolSummaries)
	// Sort by call count - descending
	sort.Slice(topUsed, func(i, j int) bool {
		return topUsed[i].CallCount > topUsed[j].CallCount
	})

	// Sort to get slowest tools
	slowest := make([]GrafanaToolSummary, len(toolSummaries))
	copy(slowest, toolSummaries)
	// Sort by average time - descending
	sort.Slice(slowest, func(i, j int) bool {
		return slowest[i].AverageTime > slowest[j].AverageTime
	})

	// Sort to get largest response tools
	largest := make([]GrafanaToolSummary, len(toolSummaries))
	copy(largest, toolSummaries)
	// Sort by response size - descending
	sort.Slice(largest, func(i, j int) bool {
		return largest[i].AverageSize > largest[j].AverageSize
	})

	var avgSize int64
	if totalCalls > 0 {
		avgSize = totalSize / int64(len(toolSummaries))
	}

	// Simplified return - should implement complete sorting logic
	return GrafanaPerformanceReport{
		TopUsedTools:        topUsed[:min(3, len(topUsed))], // Show top 3
		SlowestTools:        slowest[:min(3, len(slowest))], // Show top 3
		LargestResponses:    largest[:min(3, len(largest))], // Show top 3
		TotalCalls:          totalCalls,
		TotalErrors:         totalErrors,
		AverageResponseSize: avgSize,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
