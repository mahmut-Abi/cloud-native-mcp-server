// Package client provides summary extraction utilities for Prometheus resources.
package client

type TargetSummary struct {
	Job      string
	Instance string
	Health   string
}

type AlertSummary struct {
	AlertName string
	State     string
}

type RuleSummary struct {
	Name string
	Type string
}

type SeriesSummary struct {
	Metric string
	Points int
}
