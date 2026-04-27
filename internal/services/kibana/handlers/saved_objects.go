// Package handlers provides HTTP handlers for Kibana MCP operations.
// This file contains saved object-related handlers.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
)

// HandleSearchSavedObjects handles Kibana saved objects search requests.
func HandleSearchSavedObjects(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana search saved objects handler")

		// Get optional parameters
		objectType := getOptionalStringParam(req, "type")
		search := getOptionalStringParam(req, "search")
		page := getOptionalIntParam(req, "page", 1)
		perPage := getOptionalIntParam(req, "per_page", 20)

		// Search saved objects
		result, err := c.SearchSavedObjects(ctx, objectType, search, page, perPage)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to search saved objects: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(result)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format search results: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetSavedSearches handles Kibana saved searches retrieval requests.
func HandleGetSavedSearches(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get saved searches handler")

		// Get saved searches
		savedSearches, err := c.GetSavedSearches(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get saved searches: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(savedSearches)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format saved searches: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleGetSavedSearch handles specific Kibana saved search retrieval requests.
func HandleGetSavedSearch(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get saved search handler")

		// Extract parameters
		args := req.Params.Arguments
		if args == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("No arguments provided"),
				},
			}, nil
		}

		// Get search ID parameter
		searchID, err := requireStringParam(req, "search_id")
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("Search ID parameter is required"),
				},
			}, nil
		}

		// Get saved search
		savedSearch, err := c.GetSavedSearch(ctx, searchID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get saved search: %v", err)),
				},
			}, nil
		}

		// Format result
		resultJSON, err := marshalIndentJSON(savedSearch)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format saved search: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleCreateSavedObject handles creating a generic saved object
func HandleCreateSavedObject(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana create saved object handler")

		objectType := getOptionalStringParam(req, "type")
		initialObjectType := getOptionalStringParam(req, "initialObjectType")
		attributes, err := getOptionalObjectParam(req, "attributes")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		referenceObjects, err := getOptionalObjectArrayParam(req, "references")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var references []client.Reference
		for _, refMap := range referenceObjects {
			references = append(references, client.Reference{
				Name: getStringFieldFromMap(refMap, "name"),
				Type: getStringFieldFromMap(refMap, "type"),
				ID:   getStringFieldFromMap(refMap, "id"),
			})
		}

		if objectType == "" || len(attributes) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("type and attributes are required"),
				},
			}, nil
		}

		obj, err := c.CreateSavedObject(ctx, objectType, attributes, references, initialObjectType)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create saved object: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(obj)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleUpdateSavedObject handles updating a saved object
func HandleUpdateSavedObject(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana update saved object handler")

		objectType := getOptionalStringParam(req, "type")
		objectID := getOptionalStringParam(req, "object_id")
		version := getOptionalStringParam(req, "version")
		attributes, err := getOptionalObjectParam(req, "attributes")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		referenceObjects, err := getOptionalObjectArrayParam(req, "references")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var references []client.Reference
		for _, refMap := range referenceObjects {
			references = append(references, client.Reference{
				Name: getStringFieldFromMap(refMap, "name"),
				Type: getStringFieldFromMap(refMap, "type"),
				ID:   getStringFieldFromMap(refMap, "id"),
			})
		}

		if objectType == "" || objectID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("type and object_id are required"),
				},
			}, nil
		}

		obj, err := c.UpdateSavedObject(ctx, objectType, objectID, attributes, references, version)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update saved object: %v", err)),
				},
			}, nil
		}

		resultJSON, err := json.Marshal(obj)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format response: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// HandleDeleteSavedObject handles deleting a saved object
func HandleDeleteSavedObject(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana delete saved object handler")

		objectType := getOptionalStringParam(req, "type")
		objectID := getOptionalStringParam(req, "object_id")

		force := false
		if f := getOptionalBoolParam(req, "force"); f != nil {
			force = *f
		}

		if objectType == "" || objectID == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("type and object_id are required"),
				},
			}, nil
		}

		err := c.DeleteSavedObject(ctx, objectType, objectID, force)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete saved object: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted saved object: %s/%s", objectType, objectID)),
			},
		}, nil
	}
}

// HandleBulkDeleteSavedObjects handles bulk deleting saved objects
func HandleBulkDeleteSavedObjects(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana bulk delete saved objects handler")

		objectMaps, err := getOptionalObjectArrayParam(req, "objects")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var objects []client.SavedObject
		for _, objMap := range objectMaps {
			objects = append(objects, client.SavedObject{
				Type: getStringFieldFromMap(objMap, "type"),
				ID:   getStringFieldFromMap(objMap, "id"),
			})
		}

		if len(objects) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("objects array is required"),
				},
			}, nil
		}

		err = c.BulkDeleteSavedObjects(ctx, objects)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to bulk delete saved objects: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted %d saved objects", len(objects))),
			},
		}, nil
	}
}

// HandleExportSavedObjects handles exporting saved objects
func HandleExportSavedObjects(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana export saved objects handler")

		objectMaps, err := getOptionalObjectArrayParam(req, "objects")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var objects []client.SavedObject
		for _, objMap := range objectMaps {
			objects = append(objects, client.SavedObject{
				Type: getStringFieldFromMap(objMap, "type"),
				ID:   getStringFieldFromMap(objMap, "id"),
			})
		}

		includeReferences := true
		if ir := getOptionalBoolParam(req, "includeReferences"); ir != nil {
			includeReferences = *ir
		}

		if len(objects) == 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("objects array is required"),
				},
			}, nil
		}

		data, err := c.ExportSavedObjects(ctx, objects, includeReferences)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to export saved objects: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(data)),
			},
		}, nil
	}
}

// HandleImportSavedObjects handles importing saved objects
func HandleImportSavedObjects(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana import saved objects handler")

		fileContent := getOptionalStringParam(req, "file")

		createNewCopies := false
		if cnc := getOptionalBoolParam(req, "createNewCopies"); cnc != nil {
			createNewCopies = *cnc
		}

		if fileContent == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("file is required"),
				},
			}, nil
		}

		err := c.ImportSavedObjects(ctx, fileContent, createNewCopies)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to import saved objects: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Successfully imported saved objects"),
			},
		}, nil
	}
}

// HandleSearchSavedObjectsAdvanced handles advanced saved objects search with enhanced filters
func HandleSearchSavedObjectsAdvanced(c *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		objectType := getOptionalStringParam(request, "type")
		search := getOptionalStringParam(request, "search")
		page := getOptionalIntParam(request, "page", 1)
		perPage := getOptionalIntParam(request, "per_page", 30)
		sortField := getOptionalStringParam(request, "sort_field")
		sortOrder := getOptionalStringParam(request, "sort_order")
		hasReference := getOptionalStringParam(request, "has_reference")

		fields, err := getOptionalStringArrayParam(request, "fields")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		logrus.WithFields(logrus.Fields{
			"tool":         "kibana_search_saved_objects_advanced",
			"objectType":   objectType,
			"search":       search,
			"page":         page,
			"perPage":      perPage,
			"sortField":    sortField,
			"sortOrder":    sortOrder,
			"hasReference": hasReference,
			"fields":       fields,
		}).Debug("Handler invoked")

		result, err := c.SearchSavedObjectsAdvanced(ctx, objectType, search, page, perPage, sortField, sortOrder, hasReference, fields)
		if err != nil {
			return nil, fmt.Errorf("failed to search saved objects advanced: %w", err)
		}

		response := map[string]interface{}{
			"savedObjects": result.SavedObjects,
			"count":        len(result.SavedObjects),
			"searchCriteria": map[string]interface{}{
				"objectType":   objectType,
				"search":       search,
				"sortField":    sortField,
				"sortOrder":    sortOrder,
				"hasReference": hasReference,
				"fields":       fields,
			},
			"pagination": map[string]interface{}{
				"currentPage": result.Page,
				"perPage":     result.PerPage,
				"total":       result.Total,
			},
			"metadata": map[string]interface{}{
				"tool":         "kibana_search_saved_objects_advanced",
				"optimizedFor": "finding specific objects",
			},
		}

		return marshalOptimizedResponse(response, "kibana_search_saved_objects_advanced")
	}
}
