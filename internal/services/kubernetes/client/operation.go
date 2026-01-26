package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/scale"
	"k8s.io/client-go/tools/remotecommand"
)

// PaginationInfo represents pagination metadata for API responses
type PaginationInfo struct {
	ContinueToken   string `json:"continueToken"`
	RemainingCount  int64  `json:"remainingCount"`
	CurrentPageSize int64  `json:"currentPageSize"`
	HasMore         bool   `json:"hasMore"`
}

var kindAliasMap = map[string]string{
	"pods":        "Pod",
	"deployments": "Deployment",
	"services":    "Service",
	"configmaps":  "ConfigMap",
	"secrets":     "Secret",
	"namespaces":  "Namespace",
	"nodes":       "Node",
	"ingresses":   "Ingress",
	"jobs":        "Job",
	"cronjobs":    "CronJob",
}

// APIResource represents a structured API resource information
type APIResource struct {
	Name         string   `json:"name"`
	SingularName string   `json:"singularName"`
	Namespaced   bool     `json:"namespaced"`
	Kind         string   `json:"kind"`
	Group        string   `json:"group"`
	Version      string   `json:"version"`
	Verbs        []string `json:"verbs"`
}

// GetAPIResources retrieves all API resource types in the cluster
func (c *Client) GetAPIResources(ctx context.Context, apiGroup string, namespaced *bool) ([]APIResource, error) {
	logrus.WithFields(logrus.Fields{"apiGroup": apiGroup, "namespaced": namespaced}).Debug("GetAPIResources called")
	resourceLists, err := c.cacheDiscovery.GetAPIResources()
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return nil, fmt.Errorf("failed to get API resources: %w", err)
	}

	var resources []APIResource

	for _, resourceList := range resourceLists {
		gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
		if err != nil {
			continue
		}

		// Filter by API group if specified
		if apiGroup != "" && gv.Group != apiGroup {
			continue
		}

		for _, resource := range resourceList.APIResources {
			// Skip subresources
			if strings.Contains(resource.Name, "/") {
				continue
			}

			// Filter resources based on their namespace scope if specified
			if namespaced != nil && resource.Namespaced != *namespaced {
				continue
			}

			resources = append(resources, APIResource{
				Name:         resource.Name,
				SingularName: resource.SingularName,
				Namespaced:   resource.Namespaced,
				Kind:         resource.Kind,
				Group:        gv.Group,
				Version:      gv.Version,
				Verbs:        resource.Verbs,
			})
		}
	}

	logrus.WithField("count", len(resources)).Debug("GetAPIResources returned")
	return resources, nil
}

// GetAPIVersions retrieves all available API versions in the cluster
func (c *Client) GetAPIVersions(ctx context.Context) ([]string, error) {
	logrus.Debug("GetAPIVersions called")

	serverGroups, err := c.cacheDiscovery.GetServerGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get server groups: %w", err)
	}

	var versions []string

	// Add core API version
	versions = append(versions, "v1")

	// Add all group versions
	for _, group := range serverGroups.Groups {
		for _, version := range group.Versions {
			if group.Name == "" {
				// Core group
				versions = append(versions, version.Version)
			} else {
				// Named group
				versions = append(versions, version.GroupVersion)
			}
		}
	}

	logrus.WithField("count", len(versions)).Debug("GetAPIVersions returned")
	return versions, nil
}

// GetResource retrieves detailed information about a specific resource
func (c *Client) GetResource(ctx context.Context, kind, name, namespace string) (map[string]any, error) {
	logrus.WithFields(logrus.Fields{"kind": kind, "name": name, "namespace": namespace}).Debug("GetResource called")
	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return nil, err
	}

	var resourceClient dynamic.ResourceInterface
	if namespace != "" {
		resourceClient = c.dynamicClient.Resource(*gvr).Namespace(namespace)
	} else {
		resourceClient = c.dynamicClient.Resource(*gvr)
	}

	obj, err := resourceClient.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get resource %s/%s: %w", kind, name, err)
	}

	logrus.Debug("GetResource succeeded")
	return obj.UnstructuredContent(), nil
}

// ListResources lists all instances of a specific resource type
func (c *Client) ListResources(ctx context.Context, kind, namespace string, labelSelector, fieldSelector string) ([]map[string]any, error) {
	return c.ListResourcesWithPagination(ctx, kind, namespace, labelSelector, fieldSelector, "", 0)
}

// ListResourcesWithPagination lists instances of a specific resource type with pagination support
func (c *Client) ListResourcesWithPagination(ctx context.Context, kind, namespace string, labelSelector, fieldSelector, continueToken string, limit int64) ([]map[string]any, error) {
	logrus.WithFields(logrus.Fields{
		"kind":      kind,
		"namespace": namespace,
		"labels":    labelSelector,
		"fields":    fieldSelector,
		"continue":  continueToken,
		"limit":     limit,
	}).Debug("ListResourcesWithPagination called")

	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return nil, err
	}

	options := metav1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
		Limit:         limit,
		Continue:      continueToken,
	}

	var resourceClient dynamic.ResourceInterface
	if namespace != "" {
		resourceClient = c.dynamicClient.Resource(*gvr).Namespace(namespace)
	} else {
		resourceClient = c.dynamicClient.Resource(*gvr)
	}

	list, err := resourceClient.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources of kind %s: %w", kind, err)
	}

	if list == nil {
		return []map[string]any{}, nil
	}

	resources := make([]map[string]any, 0, len(list.Items))
	for _, item := range list.Items {
		resources = append(resources, item.UnstructuredContent())
	}

	// Extract pagination metadata from list options
	var extractedContinueToken string
	var extractedRemainingCount int64
	listObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(list)
	if err == nil {
		if continueVal, ok := listObj["continue"].(string); ok {
			extractedContinueToken = continueVal
		}
		if remaining, ok := listObj["remainingItemCount"].(int64); ok {
			extractedRemainingCount = remaining
		}
	}

	logrus.WithFields(logrus.Fields{
		"count":     len(resources),
		"continue":  extractedContinueToken,
		"remaining": extractedRemainingCount,
	}).Debug("ListResourcesWithPagination succeeded")

	return resources, nil
}

// GetPaginationInfo returns pagination metadata for resource listings
func (c *Client) GetPaginationInfo(ctx context.Context, kind, namespace string, labelSelector, fieldSelector, continueToken string, limit int64) (*PaginationInfo, error) {
	logrus.WithFields(logrus.Fields{
		"kind":      kind,
		"namespace": namespace,
		"labels":    labelSelector,
		"fields":    fieldSelector,
		"continue":  continueToken,
		"limit":     limit,
	}).Debug("GetPaginationInfo called")

	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return nil, err
	}

	options := metav1.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
		Limit:         limit,
		Continue:      continueToken,
	}

	var resourceClient dynamic.ResourceInterface
	if namespace != "" {
		resourceClient = c.dynamicClient.Resource(*gvr).Namespace(namespace)
	} else {
		resourceClient = c.dynamicClient.Resource(*gvr)
	}

	list, err := resourceClient.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources of kind %s: %w", kind, err)
	}

	// Extract pagination metadata from unstructured list
	var extractedContinueToken string
	var extractedRemainingCount int64
	if list != nil {
		listObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(list)
		if err == nil {
			if continueVal, ok := listObj["continue"].(string); ok {
				extractedContinueToken = continueVal
			}
			if remaining, ok := listObj["remainingItemCount"].(int64); ok {
				extractedRemainingCount = remaining
			}
		}
	}

	paginationInfo := &PaginationInfo{
		ContinueToken:   extractedContinueToken,
		RemainingCount:  extractedRemainingCount,
		CurrentPageSize: int64(len(list.Items)),
		HasMore:         extractedContinueToken != "",
	}

	logrus.WithFields(logrus.Fields{
		"continue":  paginationInfo.ContinueToken,
		"remaining": paginationInfo.RemainingCount,
		"pageSize":  paginationInfo.CurrentPageSize,
		"hasMore":   paginationInfo.HasMore,
	}).Debug("GetPaginationInfo succeeded")

	return paginationInfo, nil
}

// GetResourcesDetail retrieves detailed information for multiple resources efficiently
func (c *Client) GetResourcesDetail(ctx context.Context, kind string, names []string, namespace string, includeEvents, includeStatus bool) (map[string]map[string]any, error) {
	logrus.WithFields(logrus.Fields{
		"kind":      kind,
		"names":     len(names),
		"namespace": namespace,
		"events":    includeEvents,
		"status":    includeStatus,
	}).Debug("GetResourcesDetail called")

	if len(names) == 0 {
		return map[string]map[string]any{}, nil
	}

	// Limit the number of resources to prevent context overflow
	if len(names) > 20 {
		logrus.WithField("requested", len(names)).Warn("Limiting resource detail request to 20 resources to prevent context overflow")
		names = names[:20]
	}

	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return nil, err
	}

	var resourceClient dynamic.ResourceInterface
	if namespace != "" {
		resourceClient = c.dynamicClient.Resource(*gvr).Namespace(namespace)
	} else {
		resourceClient = c.dynamicClient.Resource(*gvr)
	}

	results := make(map[string]map[string]any)
	var errors []string

	// Batch gather resources with concurrency control
	semaphore := make(chan struct{}, 5) // Limit concurrent requests
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, name := range names {
		if name == "" {
			continue
		}

		wg.Add(1)
		go func(resourceName string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Get the main resource
			obj, err := resourceClient.Get(ctx, resourceName, metav1.GetOptions{})
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Sprintf("Failed to get %s: %v", resourceName, err))
				mu.Unlock()
				return
			}

			resource := obj.UnstructuredContent()

			// Optionally remove status to reduce response size
			if !includeStatus {
				delete(resource, "status")
			}

			// Optionally include related events
			if includeEvents {
				events, err := c.getResourceEvents(ctx, resourceName, namespace, kind)
				if err == nil {
					resource["relatedEvents"] = events
				} else {
					logrus.WithError(err).WithField("resource", resourceName).Warn("Failed to get resource events")
				}
			}

			mu.Lock()
			results[resourceName] = resource
			mu.Unlock()

		}(name)
	}

	wg.Wait()

	// If we have errors but some results, return partial results with error information
	if len(errors) > 0 && len(results) > 0 {
		results["_errors"] = map[string]any{
			"message": "Some resources could not be retrieved",
			"errors":  errors,
		}
		logrus.WithFields(logrus.Fields{
			"successful": len(results),
			"errors":     len(errors),
		}).Warn("Partial success in GetResourcesDetail")
	} else if len(errors) > 0 {
		return nil, fmt.Errorf("failed to retrieve any resources: %s", strings.Join(errors, "; "))
	}

	logrus.WithField("count", len(results)).Debug("GetResourcesDetail succeeded")
	return results, nil
}

// getResourceEvents retrieves events related to a specific resource
func (c *Client) getResourceEvents(ctx context.Context, resourceName, namespace, kind string) ([]map[string]any, error) {
	// Create field selector for resource-specific events
	fieldSelector := fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=%s", resourceName, kind)

	options := metav1.ListOptions{
		FieldSelector: fieldSelector,
		Limit:         10, // Limit events to prevent excessive output
	}

	eventInterface := c.clientset.CoreV1().Events(namespace)
	eventList, err := eventInterface.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	var events []map[string]any
	for _, event := range eventList.Items {
		events = append(events, map[string]any{
			"type":    event.Type,
			"reason":  event.Reason,
			"message": event.Message,
			"time":    event.LastTimestamp.Format(time.RFC3339),
		})
	}

	return events, nil
}

// CreateResource creates a new resource from the given metadata and spec
func (c *Client) CreateResource(ctx context.Context, kind, apiVersion, metadataJSON, specJSON string) (map[string]any, error) {
	logrus.WithFields(logrus.Fields{"kind": kind, "apiVersion": apiVersion}).Debug("CreateResource called")
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion(apiVersion)
	obj.SetKind(kind)

	// Parse and set metadata
	var metadata map[string]any
	if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
	}

	if name, ok := metadata["name"].(string); ok {
		obj.SetName(name)
	}
	if namespace, ok := metadata["namespace"].(string); ok {
		obj.SetNamespace(namespace)
	}
	obj.Object["metadata"] = metadata

	// Parse and set spec if provided
	if specJSON != "" {
		var spec map[string]any
		if err := json.Unmarshal([]byte(specJSON), &spec); err != nil {
			return nil, fmt.Errorf("failed to parse spec JSON: %w", err)
		}
		obj.Object["spec"] = spec
	}

	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return nil, err
	}

	var resourceClient dynamic.ResourceInterface
	if obj.GetNamespace() != "" {
		resourceClient = c.dynamicClient.Resource(*gvr).Namespace(obj.GetNamespace())
	} else {
		resourceClient = c.dynamicClient.Resource(*gvr)
	}

	created, err := resourceClient.Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create resource %s: %w", kind, err)
	}

	logrus.Debug("CreateResource succeeded")
	return created.UnstructuredContent(), nil
}

// UpdateResource updates an existing resource with the provided manifest
func (c *Client) UpdateResource(ctx context.Context, kind, name, namespace string, manifest string) (map[string]any, error) {
	logrus.WithFields(logrus.Fields{"kind": kind, "name": name, "namespace": namespace}).Debug("UpdateResource called")
	obj := &unstructured.Unstructured{}
	if err := json.Unmarshal([]byte(manifest), &obj.Object); err != nil {
		return nil, fmt.Errorf("failed to parse resource manifest: %w", err)
	}

	if obj.GetName() != name {
		return nil, fmt.Errorf("name mismatch: manifest has %q, expected %q", obj.GetName(), name)
	}

	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return nil, err
	}

	var resourceClient dynamic.ResourceInterface
	if namespace != "" {
		resourceClient = c.dynamicClient.Resource(*gvr).Namespace(namespace)
	} else {
		resourceClient = c.dynamicClient.Resource(*gvr)
	}

	result, err := resourceClient.Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update resource %s/%s: %w", kind, name, err)
	}

	logrus.Debug("UpdateResource succeeded")
	return result.UnstructuredContent(), nil
}

// DeleteResource deletes a resource by its kind, name, and namespace
func (c *Client) DeleteResource(ctx context.Context, kind, name, namespace string) error {
	logrus.WithFields(logrus.Fields{"kind": kind, "name": name, "namespace": namespace}).Debug("DeleteResource called")
	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return err
	}

	var resourceClient dynamic.ResourceInterface
	if namespace != "" {
		resourceClient = c.dynamicClient.Resource(*gvr).Namespace(namespace)
	} else {
		resourceClient = c.dynamicClient.Resource(*gvr)
	}

	if err := resourceClient.Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to delete resource %s/%s: %w", kind, name, err)
	}

	logrus.Debug("DeleteResource succeeded")
	return nil
}

// normalizeKind normalizes the kind string using alias map or title case
func normalizeKind(kind string) string {
	normalized := strings.ToLower(kind)
	if realKind, ok := kindAliasMap[normalized]; ok {
		return realKind
	}
	return cases.Title(language.English).String(normalized)
}

// findGroupVersionResource finds the corresponding GroupVersionResource by Kind with improved caching
func (c *Client) findGroupVersionResource(kind string) (*schema.GroupVersionResource, error) {
	// Normalize the kind using the existing normalizeKind function
	kind = normalizeKind(kind)

	// Try cache first - should handle most cases
	if gvr, err := c.getCachedGVR(kind); err == nil {
		logrus.WithField("kind", kind).Debug("Resolved GVR via cache")
		return gvr, nil
	}

	logrus.WithField("kind", kind).Debug("Cache resolution failed, using improved discovery")

	// Default to core API version first
	gvr, err := c.findGVRUsingDiscovery(kind, "v1")
	if err == nil {
		c.updateGVRCache(kind, gvr)
		return &gvr, nil
	}

	// If not found, try without version (will use preferred version)
	gvr, err = c.findGVRUsingDiscovery(kind, "")
	if err == nil {
		c.updateGVRCache(kind, gvr)
		return &gvr, nil
	}

	// If still not found, fall back to full discovery
	return c.discoverAndCacheGVR(kind)
}

// findGVRUsingDiscovery finds GVR using the discovery client
func (c *Client) findGVRUsingDiscovery(kind, version string) (schema.GroupVersionResource, error) {
	// Force refresh cache by invalidating it first
	c.cacheDiscovery.Invalidate()

	resourceLists, err := c.cacheDiscovery.GetAPIResources()
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	for _, resourceList := range resourceLists {
		gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
		if err != nil {
			continue
		}

		if version != "" && gv.Version != version {
			continue
		}

		for _, resource := range resourceList.APIResources {
			if strings.EqualFold(resource.Kind, kind) {
				return schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: resource.Name,
				}, nil
			}
		}
	}

	return schema.GroupVersionResource{}, fmt.Errorf("resource kind %s not found", kind)
}

// discoverAndCacheGVR discovers GVR via API and updates cache
func (c *Client) discoverAndCacheGVR(kind string) (*schema.GroupVersionResource, error) {
	// Force refresh cache by invalidating it first
	c.cacheDiscovery.Invalidate()

	resourceLists, err := c.cacheDiscovery.GetAPIResources()
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return nil, fmt.Errorf("failed to get API resources: %w", err)
	}

	c.gvrCacheMux.Lock()
	defer c.gvrCacheMux.Unlock()

	// Update cache expiry time
	c.cacheExpiry = time.Now().Add(c.cacheTTL)

	for _, resourceList := range resourceLists {
		gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
		if err != nil {
			continue
		}

		for _, resource := range resourceList.APIResources {
			// Skip subresources like pods/status
			if strings.Contains(resource.Name, "/") {
				continue
			}

			// Cache all discovered resources for future lookups
			gvr := schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: resource.Name,
			}
			c.gvrCache[toLower(resource.Kind)] = gvr

			// Check if this is our target kind
			if strings.EqualFold(resource.Kind, kind) {
				logrus.WithField("kind", kind).Debug("Resolved GVR via discovery and cached")
				return &gvr, nil
			}
		}
	}

	return nil, fmt.Errorf("resource kind %q not found", kind)
}

// ExecCommand executes a command in a container
func (c *Client) ExecCommand(ctx context.Context, podName, namespace, container string, command []string) (string, error) {
	logrus.WithFields(logrus.Fields{"pod": podName, "ns": namespace, "container": container, "cmd": strings.Join(command, " ")}).Debug("ExecCommand called")
	req := c.clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.restConfig, "POST", req.URL())
	if err != nil {
		return "", fmt.Errorf("failed to create executor: %w", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w, stderr: %s", err, stderr.String())
	}

	logrus.Debug("ExecCommand succeeded")
	return stdout.String(), nil
}

// GetContainerLog retrieves logs for a specific container in a pod
func (c *Client) GetContainerLog(ctx context.Context, podName, namespace, container string, tailLines int64) (string, error) {
	logrus.WithFields(logrus.Fields{"pod": podName, "ns": namespace, "container": container, "tail": tailLines}).Debug("GetContainerLog called")
	req := c.clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: container,
		Follow:    false,
		TailLines: &tailLines,
	})

	stream, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get log stream: %w", err)
	}
	defer func() { _ = stream.Close() }()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, stream); err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	logrus.Debug("GetContainerLog succeeded")
	return buf.String(), nil
}

// CheckPermissions checks if the current user has the specified permissions
func (c *Client) CheckPermissions(ctx context.Context, verb, resourceName, resourceGroup, resourceResource, subresource, namespace string) (bool, error) {
	logrus.WithFields(logrus.Fields{"verb": verb, "group": resourceGroup, "resource": resourceResource, "subresource": subresource, "name": resourceName, "ns": namespace}).Debug("CheckPermissions called")
	sar := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace:   namespace,
				Verb:        verb,
				Group:       resourceGroup,
				Resource:    resourceResource,
				Subresource: subresource,
				Name:        resourceName,
			},
		},
	}

	response, err := c.authClient.SelfSubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to check permissions: %w", err)
	}

	logrus.WithField("allowed", response.Status.Allowed).Debug("CheckPermissions result")
	return response.Status.Allowed, nil
}

// ScaleResource scales a resource (placeholder implementation)
func (c *Client) ScaleResource(ctx context.Context, gvr schema.GroupVersionResource, name, namespace string, replicas int32) error {
	logrus.WithFields(logrus.Fields{"group": gvr.Group, "resource": gvr.Resource, "name": name, "ns": namespace, "replicas": replicas}).Debug("ScaleResource called")
	restClient := c.clientset.CoreV1().RESTClient()
	mapper, err := genericclioptions.NewTestConfigFlags().ToRESTMapper()
	if err != nil {
		return fmt.Errorf("failed to create REST mapper: %w", err)
	}

	resolver := scale.NewDiscoveryScaleKindResolver(c.discoveryClient)
	scaleClient := scale.New(restClient, mapper, dynamic.LegacyAPIPathResolverFunc, resolver)

	gr := schema.GroupResource{Group: gvr.Group, Resource: gvr.Resource}
	scaleObj, err := scaleClient.Scales(namespace).Get(ctx, gr, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get scale for resource: %w", err)
	}

	scaleObj.Spec.Replicas = replicas
	_, err = scaleClient.Scales(namespace).Update(ctx, gr, scaleObj, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update scale for resource: %w", err)
	}

	logrus.Debug("ScaleResource succeeded")
	return nil
}

// ScaleResourceByKind resolves GVR by kind and scales the resource
func (c *Client) ScaleResourceByKind(ctx context.Context, kind, name, namespace string, replicas int32) error {
	logrus.WithFields(logrus.Fields{"kind": kind, "name": name, "ns": namespace, "replicas": replicas}).Debug("ScaleResourceByKind called")
	gvr, err := c.findGroupVersionResource(kind)
	if err != nil {
		return err
	}
	return c.ScaleResource(ctx, *gvr, name, namespace, replicas)
}
