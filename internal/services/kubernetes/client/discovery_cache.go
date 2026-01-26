package client

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

// CacheableDiscovery wraps the standard discovery client with caching capabilities
type CacheableDiscovery struct {
	discoveryClient discovery.DiscoveryInterface
	cacheMutex      sync.RWMutex
	apiResources    map[string]*metav1.APIResourceList // cached API resources
	serverGroups    *metav1.APIGroupList               // cached server groups
	cacheExpiry     time.Time                          // when the cache expires
	cacheTTL        time.Duration                      // cache time-to-live
}

// NewCacheableDiscovery creates a new CacheableDiscovery instance
func NewCacheableDiscovery(client discovery.DiscoveryInterface, ttl time.Duration) *CacheableDiscovery {
	return &CacheableDiscovery{
		discoveryClient: client,
		apiResources:    make(map[string]*metav1.APIResourceList),
		cacheTTL:        ttl,
	}
}

// GetAPIResources returns cached API resources or fetches them if cache is expired
func (c *CacheableDiscovery) GetAPIResources() ([]*metav1.APIResourceList, error) {
	c.cacheMutex.RLock()
	// Check if cache is valid
	if time.Now().Before(c.cacheExpiry) && len(c.apiResources) > 0 {
		// Convert map to slice
		resources := make([]*metav1.APIResourceList, 0, len(c.apiResources))
		for _, resource := range c.apiResources {
			resources = append(resources, resource)
		}
		c.cacheMutex.RUnlock()
		logrus.Debug("Returning cached API resources")
		return resources, nil
	}
	c.cacheMutex.RUnlock()

	// Cache miss or expired, fetch fresh data
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	// Double-check after acquiring write lock
	if time.Now().Before(c.cacheExpiry) && len(c.apiResources) > 0 {
		// Convert map to slice
		resources := make([]*metav1.APIResourceList, 0, len(c.apiResources))
		for _, resource := range c.apiResources {
			resources = append(resources, resource)
		}
		logrus.Debug("Returning cached API resources (double-checked)")
		return resources, nil
	}

	// Fetch fresh data
	logrus.Debug("Fetching fresh API resources from server")
	resourceLists, err := c.discoveryClient.ServerPreferredResources()
	if err != nil {
		// If we have cached data, return it even if expired
		if len(c.apiResources) > 0 {
			logrus.WithError(err).Warn("Failed to fetch API resources, returning cached data")
			resources := make([]*metav1.APIResourceList, 0, len(c.apiResources))
			for _, resource := range c.apiResources {
				resources = append(resources, resource)
			}
			return resources, nil
		}
		return nil, err
	}

	// Update cache
	c.apiResources = make(map[string]*metav1.APIResourceList)
	for _, rl := range resourceLists {
		c.apiResources[rl.GroupVersion] = rl
	}
	c.cacheExpiry = time.Now().Add(c.cacheTTL)

	// Return fresh data
	resources := make([]*metav1.APIResourceList, len(resourceLists))
	copy(resources, resourceLists)

	return resources, nil
}

// GetServerGroups returns cached server groups or fetches them if cache is expired
func (c *CacheableDiscovery) GetServerGroups() (*metav1.APIGroupList, error) {
	c.cacheMutex.RLock()
	// Check if cache is valid
	if time.Now().Before(c.cacheExpiry) && c.serverGroups != nil {
		groups := c.serverGroups
		c.cacheMutex.RUnlock()
		logrus.Debug("Returning cached server groups")
		return groups, nil
	}
	c.cacheMutex.RUnlock()

	// Cache miss or expired, fetch fresh data
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	// Double-check after acquiring write lock
	if time.Now().Before(c.cacheExpiry) && c.serverGroups != nil {
		logrus.Debug("Returning cached server groups (double-checked)")
		return c.serverGroups, nil
	}

	// Fetch fresh data
	logrus.Debug("Fetching fresh server groups from server")
	groups, err := c.discoveryClient.ServerGroups()
	if err != nil {
		// If we have cached data, return it even if expired
		if c.serverGroups != nil {
			logrus.WithError(err).Warn("Failed to fetch server groups, returning cached data")
			return c.serverGroups, nil
		}
		return nil, err
	}

	// Update cache
	c.serverGroups = groups
	c.cacheExpiry = time.Now().Add(c.cacheTTL)

	return groups, nil
}

// FindGVR finds the GroupVersionResource for a given kind, with caching
func (c *CacheableDiscovery) FindGVR(kind, version string) (schema.GroupVersionResource, error) {
	// Get API resources (cached)
	resourceLists, err := c.GetAPIResources()
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	// Search for the resource
	for _, resourceList := range resourceLists {
		gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
		if err != nil {
			continue
		}

		// If version is specified, check if it matches
		if version != "" && gv.Version != version {
			continue
		}

		for _, resource := range resourceList.APIResources {
			// Check if kind matches (case-insensitive)
			if resource.Kind == kind {
				return schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: resource.Name,
				}, nil
			}
		}
	}

	return schema.GroupVersionResource{}, errors.NewNotFound(schema.GroupResource{}, kind)
}

// Refresh forces a cache refresh
func (c *CacheableDiscovery) Refresh() error {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	// Clear cache
	c.apiResources = make(map[string]*metav1.APIResourceList)
	c.serverGroups = nil
	c.cacheExpiry = time.Time{}

	// Fetch fresh data
	resourceLists, err := c.discoveryClient.ServerPreferredResources()
	if err != nil {
		return err
	}

	groups, err := c.discoveryClient.ServerGroups()
	if err != nil {
		return err
	}

	// Update cache
	for _, rl := range resourceLists {
		c.apiResources[rl.GroupVersion] = rl
	}
	c.serverGroups = groups
	c.cacheExpiry = time.Now().Add(c.cacheTTL)

	return nil
}

// Invalidate clears the cache
func (c *CacheableDiscovery) Invalidate() {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	c.apiResources = make(map[string]*metav1.APIResourceList)
	c.serverGroups = nil
	c.cacheExpiry = time.Time{}
}
