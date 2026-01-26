// Package client provides Helm client operations for the MCP server.
package client

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/repo"
)

// RepoHealthCache tracks the health status of Helm repositories
type RepoHealthCache struct {
	status  map[string]bool      // URL -> healthy status
	expiry  map[string]time.Time // URL -> expiry time
	mu      sync.RWMutex         // protect maps
	timeout time.Duration        // timeout for health check
	ttl     time.Duration        // time to live for cache entry
}

// NewRepoHealthCache creates a new repository health cache
func NewRepoHealthCache(timeout time.Duration, ttl time.Duration) *RepoHealthCache {
	return &RepoHealthCache{
		status:  make(map[string]bool),
		expiry:  make(map[string]time.Time),
		timeout: timeout,
		ttl:     ttl,
	}
}

// IsHealthy checks if a repository is healthy, using cache if available
func (rhc *RepoHealthCache) IsHealthy(repoURL string) bool {
	rhc.mu.RLock()
	if status, exists := rhc.status[repoURL]; exists {
		if time.Now().Before(rhc.expiry[repoURL]) {
			rhc.mu.RUnlock()
			return status
		}
	}
	rhc.mu.RUnlock()

	// Not in cache or expired, check health
	return rhc.checkAndCache(repoURL)
}

// checkAndCache performs an actual health check and caches the result
func (rhc *RepoHealthCache) checkAndCache(repoURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), rhc.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", repoURL+"/index.yaml", nil)
	if err != nil {
		logrus.WithField("repo_url", repoURL).WithError(err).Debug("Failed to create health check request")
		rhc.setCached(repoURL, false)
		return false
	}

	client := &http.Client{Timeout: rhc.timeout}
	resp, err := client.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if err != nil {
		logrus.WithField("repo_url", repoURL).WithError(err).Debug("Repository health check failed")
		rhc.setCached(repoURL, false)
		return false
	}

	healthy := resp.StatusCode >= 200 && resp.StatusCode < 300
	logrus.WithField("repo_url", repoURL).WithField("status_code", resp.StatusCode).Debug("Repository health check result")
	rhc.setCached(repoURL, healthy)
	return healthy
}

// setCached stores the health check result in cache
func (rhc *RepoHealthCache) setCached(repoURL string, healthy bool) {
	rhc.mu.Lock()
	defer rhc.mu.Unlock()
	rhc.status[repoURL] = healthy
	rhc.expiry[repoURL] = time.Now().Add(rhc.ttl)
}

// FilterHealthyRepos returns only healthy repositories from the given list
func (rhc *RepoHealthCache) FilterHealthyRepos(repos []*repo.Entry) []*repo.Entry {
	var healthy []*repo.Entry
	for _, r := range repos {
		if r != nil && rhc.IsHealthy(r.URL) {
			healthy = append(healthy, r)
		}
	}
	return healthy
}

// InvalidateRepo removes a repository from the cache
func (rhc *RepoHealthCache) InvalidateRepo(repoURL string) {
	rhc.mu.Lock()
	defer rhc.mu.Unlock()
	delete(rhc.status, repoURL)
	delete(rhc.expiry, repoURL)
}

// InvalidateAll clears all cached entries
func (rhc *RepoHealthCache) InvalidateAll() {
	rhc.mu.Lock()
	defer rhc.mu.Unlock()
	rhc.status = make(map[string]bool)
	rhc.expiry = make(map[string]time.Time)
}

// GetStats returns cache statistics
func (rhc *RepoHealthCache) GetStats() map[string]interface{} {
	rhc.mu.RLock()
	defer rhc.mu.RUnlock()

	healthyCount := 0
	for _, healthy := range rhc.status {
		if healthy {
			healthyCount++
		}
	}

	return map[string]interface{}{
		"total_repos":   len(rhc.status),
		"healthy_repos": healthyCount,
		"failed_repos":  len(rhc.status) - healthyCount,
		"cache_ttl":     rhc.ttl.String(),
	}
}
