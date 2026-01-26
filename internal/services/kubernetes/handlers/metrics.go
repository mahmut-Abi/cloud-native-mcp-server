package handlers

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ToolMetrics tool execution metrics
type ToolMetrics struct {
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

// MetricsCollector metrics collector
type MetricsCollector struct {
	metrics map[string]*ToolMetrics
	mutex   sync.RWMutex
}

// NewMetricsCollector creates new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make(map[string]*ToolMetrics),
	}
}

// RecordToolCall records tool call
func (collector *MetricsCollector) RecordToolCall(toolName string, duration time.Duration, responseSize int64, isError bool) {
	collector.mutex.Lock()
	defer collector.mutex.Unlock()

	metrics, exists := collector.metrics[toolName]
	if !exists {
		metrics = &ToolMetrics{
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
	}).Debug("Tool call metrics recorded")
}

// GetMetrics gets tool metrics
func (collector *MetricsCollector) GetMetrics(toolName string) *ToolMetrics {
	collector.mutex.RLock()
	defer collector.mutex.RUnlock()

	if metrics, exists := collector.metrics[toolName]; exists {
		// Return copy to avoid concurrency issues
		return &ToolMetrics{
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

// GetAllMetrics gets all tool metrics
func (collector *MetricsCollector) GetAllMetrics() map[string]*ToolMetrics {
	collector.mutex.RLock()
	defer collector.mutex.RUnlock()

	result := make(map[string]*ToolMetrics)
	for name, metrics := range collector.metrics {
		result[name] = &ToolMetrics{
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
func (metrics *ToolMetrics) GetAverageDuration() time.Duration {
	if metrics.CallCount == 0 {
		return 0
	}
	return metrics.TotalDuration / time.Duration(metrics.CallCount)
}

// GetSuccessRate gets success rate
func (metrics *ToolMetrics) GetSuccessRate() float64 {
	if metrics.CallCount == 0 {
		return 0
	}
	return float64(metrics.SuccessCount) / float64(metrics.CallCount)
}

// Global metrics collector
var DefaultMetricsCollector = NewMetricsCollector()

// MeasureToolExecution decorator to measure tool execution time
func MeasureToolExecution(toolName string, execFunc func() (*interface{}, error)) (*interface{}, error) {
	start := time.Now()

	result, err := execFunc()
	duration := time.Since(start)

	responseSize := int64(0)
	if result != nil {
		if data, err2 := json.Marshal(*result); err2 == nil {
			responseSize = int64(len(data))
		}
	}

	DefaultMetricsCollector.RecordToolCall(toolName, duration, responseSize, err != nil)

	logrus.WithFields(logrus.Fields{
		"toolName":     toolName,
		"duration":     duration,
		"responseSize": responseSize,
		"success":      err == nil,
	}).Debug("Tool execution measured")

	return result, err
}

// PerformanceReport performance report
type PerformanceReport struct {
	TopUsedTools        []ToolSummary `json:"topUsedTools"`
	SlowestTools        []ToolSummary `json:"slowestTools"`
	LargestResponses    []ToolSummary `json:"largestResponses"`
	TotalCalls          int64         `json:"totalCalls"`
	TotalErrors         int64         `json:"totalErrors"`
	AverageResponseSize int64         `json:"averageResponseSize"`
}

// ToolSummary tool summary
type ToolSummary struct {
	ToolName    string        `json:"toolName"`
	CallCount   int64         `json:"callCount"`
	AverageTime time.Duration `json:"averageTime"`
	MinTime     time.Duration `json:"minTime"`
	MaxTime     time.Duration `json:"maxTime"`
	SuccessRate float64       `json:"successRate"`
	AverageSize int64         `json:"averageSize"`
	LastCalled  time.Time     `json:"lastCalled"`
}

// GeneratePerformanceReport generates performance report
func GeneratePerformanceReport() PerformanceReport {
	allMetrics := DefaultMetricsCollector.GetAllMetrics()

	var totalCalls int64
	var totalErrors int64
	var totalSize int64

	var toolSummaries []ToolSummary

	for _, metrics := range allMetrics {
		totalCalls += metrics.CallCount
		totalErrors += metrics.ErrorCount
		totalSize += metrics.AverageSize

		summary := ToolSummary{
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
	topUsed := make([]ToolSummary, len(toolSummaries))
	copy(topUsed, toolSummaries)
	// Sort by call count

	// Sort to get slowest tools
	slowest := make([]ToolSummary, len(toolSummaries))
	copy(slowest, toolSummaries)
	// Sort by average time

	// Sort to get largest response tools
	largest := make([]ToolSummary, len(toolSummaries))
	copy(largest, toolSummaries)
	// Sort by response size

	var avgSize int64
	if totalCalls > 0 {
		avgSize = totalSize / int64(len(toolSummaries))
	}

	return PerformanceReport{
		TopUsedTools:        topUsed[:min(5, len(topUsed))],
		SlowestTools:        slowest[:min(5, len(slowest))],
		LargestResponses:    largest[:min(5, len(largest))],
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
