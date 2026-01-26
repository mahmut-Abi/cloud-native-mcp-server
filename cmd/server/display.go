package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/services"
	optimize "github.com/mahmut-Abi/k8s-mcp-server/internal/util/performance"
)

// ListDisplayOptions holds options for listing
type ListDisplayOptions struct {
	Format        string // output format: text, json, table, csv
	ServiceFilter string // Filter by service name
	Verbose       bool   // Include full descriptions
}

// DisplayServices displays all services
func DisplayServices(enabledServices map[string]services.Service, opts ListDisplayOptions) error {
	serviceList := make([]map[string]any, 0)

	for name, service := range enabledServices {
		if opts.ServiceFilter != "" && name != opts.ServiceFilter {
			continue
		}
		serviceList = append(serviceList, map[string]any{
			"name":    name,
			"enabled": service.IsEnabled(),
			"tools":   len(service.GetTools()),
		})
	}

	// Sort by name
	sort.Slice(serviceList, func(i, j int) bool {
		return serviceList[i]["name"].(string) < serviceList[j]["name"].(string)
	})

	switch opts.Format {
	case "json":
		return displayServicesJSON(serviceList)
	case "csv":
		return displayServicesCSV(serviceList)
	case "table":
		return displayServicesTable(serviceList)
	default:
		return displayServicesText(serviceList)
	}
}

// DisplayTools displays all tools
func DisplayTools(tools []mcp.Tool, opts ListDisplayOptions) error {
	// Filter tools by service if specified
	filteredTools := tools
	if opts.ServiceFilter != "" {
		filteredTools = filterToolsByService(tools, opts.ServiceFilter)
	}

	// Sort by name
	sort.Slice(filteredTools, func(i, j int) bool {
		return filteredTools[i].Name < filteredTools[j].Name
	})

	switch opts.Format {
	case "json":
		return displayToolsJSON(filteredTools)
	case "csv":
		return displayToolsCSV(filteredTools)
	case "table":
		return displayToolsTable(filteredTools)
	default:
		return displayToolsText(filteredTools, opts.Verbose)
	}
}

// displayServicesText displays services in text format
func displayServicesText(services []map[string]any) error {
	fmt.Println("Available Services:")
	fmt.Println(strings.Repeat("=", 50))
	for _, svc := range services {
		status := "enabled"
		if !svc["enabled"].(bool) {
			status = "disabled"
		}
		fmt.Printf("  %s (%s) - %d tools\n", svc["name"], status, svc["tools"])
	}
	return nil
}

// displayServicesJSON displays services in JSON format
func displayServicesJSON(services []map[string]any) error {
	data := map[string]any{
		"total":    len(services),
		"services": services,
	}
	byte, err := marshalIndentJSON(data)
	if err != nil {
		return err
	}
	fmt.Println(string(byte))
	return nil
}

// displayServicesCSV displays services in CSV format
func displayServicesCSV(services []map[string]any) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	if err := writer.Write([]string{"Service Name", "Status", "Tools Count"}); err != nil {
		return err
	}
	for _, svc := range services {
		status := "enabled"
		if !svc["enabled"].(bool) {
			status = "disabled"
		}
		if err := writer.Write([]string{
			svc["name"].(string),
			status,
			fmt.Sprintf("%d", svc["tools"]),
		}); err != nil {
			return err
		}
	}
	return nil
}

// displayServicesTable displays services in table format
func displayServicesTable(services []map[string]any) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "SERVICE\tSTATUS\tTOOLS"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, strings.Repeat("-", 40)); err != nil {
		return err
	}

	for _, svc := range services {
		status := "enabled"
		if !svc["enabled"].(bool) {
			status = "disabled"
		}
		if _, err := fmt.Fprintf(w, "%s\t%s\t%d\n", svc["name"], status, svc["tools"]); err != nil {
			return err
		}
	}
	return w.Flush()
}

// displayToolsText displays tools in text format
func displayToolsText(tools []mcp.Tool, verbose bool) error {
	fmt.Println("Available Tools:")
	fmt.Println(strings.Repeat("=", 80))
	for _, tool := range tools {
		fmt.Printf("  %s\n", tool.Name)
		if verbose && tool.Description != "" {
			desc := tool.Description
			if len(desc) > 75 {
				desc = desc[:72] + "..."
			}
			fmt.Printf("    Description: %s\n", desc)
		}
	}
	return nil
}

// displayToolsJSON displays tools in JSON format
func displayToolsJSON(tools []mcp.Tool) error {
	type ToolInfo struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	toolList := make([]ToolInfo, 0)
	for _, tool := range tools {
		toolList = append(toolList, ToolInfo{
			Name:        tool.Name,
			Description: tool.Description,
		})
	}

	data := map[string]any{
		"total": len(toolList),
		"tools": toolList,
	}

	byte, err := marshalIndentJSON(data)
	if err != nil {
		return err
	}
	fmt.Println(string(byte))
	return nil
}

// displayToolsCSV displays tools in CSV format
func displayToolsCSV(tools []mcp.Tool) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	if err := writer.Write([]string{"Tool Name", "Description"}); err != nil {
		return err
	}
	for _, tool := range tools {
		desc := tool.Description
		if len(desc) > 200 {
			desc = desc[:197] + "..."
		}
		desc = strings.ReplaceAll(desc, "\"", "\\")
		if err := writer.Write([]string{tool.Name, desc}); err != nil {
			return err
		}
	}
	return nil
}

// displayToolsTable displays tools in table format
func displayToolsTable(tools []mcp.Tool) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "TOOL NAME\tDESCRIPTION"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, strings.Repeat("-", 100)); err != nil {
		return err
	}

	for _, tool := range tools {
		desc := tool.Description
		if len(desc) > 70 {
			desc = desc[:67] + "..."
		}
		if _, err := fmt.Fprintf(w, "%s\t%s\n", tool.Name, desc); err != nil {
			return err
		}
	}
	return w.Flush()
}

// filterToolsByService filters tools by service name prefix
func filterToolsByService(tools []mcp.Tool, serviceName string) []mcp.Tool {
	var filtered []mcp.Tool
	prefix := serviceName + "_"
	for _, tool := range tools {
		if strings.HasPrefix(tool.Name, prefix) {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

// marshalIndentJSON performs indented JSON encoding using object pool
func marshalIndentJSON(data interface{}) ([]byte, error) {
	// First encode to compact format using object pool
	compactBytes, err := optimize.GlobalJSONPool.MarshalToBytes(data)
	if err != nil {
		return nil, err
	}

	// For scenarios requiring indented display, still use standard library but reduce allocations
	// This is a trade-off between performance and readability
	var result bytes.Buffer
	err = json.Indent(&result, compactBytes, "", "  ")
	return result.Bytes(), err
}
