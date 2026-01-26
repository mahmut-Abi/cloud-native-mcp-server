// Package client provides summary extraction utilities for Elasticsearch resources.
package client

type IndexSummary struct {
	Name      string
	Health    string
	Status    string
	DocsCount int64
}

type NodeSummary struct {
	NodeID string
	Name   string
	Role   string
}

type ClusterHealth struct {
	Status        string
	NumberOfNodes int
	ActiveShards  int
}
