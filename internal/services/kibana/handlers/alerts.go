// Package handlers provides HTTP handlers for Kibana MCP operations.
// This file contains alert-related handlers.
package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/mahmut-Abi/cloud-native-mcp-server/internal/services/kibana/client"
)

// HandleGetKibanaAlerts handles Kibana alerting rules retrieval requests.
func HandleGetKibanaAlerts(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get alerts handler")

		alerts, err := c.GetAlerts(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get Kibana alerts: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(alerts)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alerts: %v", err)),
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

// HandleGetAlertRules handles listing alert rules.
func HandleGetAlertRules(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		page := getOptionalIntParam(req, "page", 1)
		perPage := getOptionalIntParam(req, "per_page", 20)
		filter := getOptionalStringParam(req, "filter")
		enabled := getOptionalBoolParam(req, "enabled")

		logrus.WithFields(logrus.Fields{
			"page":    page,
			"perPage": perPage,
			"filter":  filter,
		}).Debug("Executing Kibana get alert rules handler")

		rules, err := c.GetAlertRules(ctx, page, perPage, filter, enabled)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get alert rules: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(rules)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alert rules: %v", err)),
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

// HandleGetAlertRule handles getting a specific alert rule.
func HandleGetAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana get alert rule handler")

		rule, err := c.GetAlertRule(ctx, ruleID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get alert rule: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(rule)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alert rule: %v", err)),
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

// HandleCreateAlertRule handles creating a new alert rule.
func HandleCreateAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := getOptionalStringParam(req, "name")
		alertTypeID := getOptionalStringParam(req, "alertTypeId")
		schedule, err := getOptionalObjectParam(req, "schedule")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		params, err := getOptionalObjectParam(req, "params")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		actions, err := getOptionalObjectArrayParam(req, "actions")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		tags, err := getOptionalStringArrayParam(req, "tags")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if name == "" || alertTypeID == "" || schedule == nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("name, alertTypeId, and schedule are required"),
				},
			}, nil
		}

		logrus.WithField("name", name).Debug("Executing Kibana create alert rule handler")

		rule, err := c.CreateAlertRule(ctx, name, alertTypeID, schedule, params, actions, tags)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to create alert rule: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(rule)
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

// HandleUpdateAlertRule handles updating an existing alert rule.
func HandleUpdateAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		name := getOptionalStringParam(req, "name")
		schedule, err := getOptionalObjectParam(req, "schedule")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		params, err := getOptionalObjectParam(req, "params")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		actions, err := getOptionalObjectArrayParam(req, "actions")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		tags, err := getOptionalStringArrayParam(req, "tags")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana update alert rule handler")

		rule, err := c.UpdateAlertRule(ctx, ruleID, name, schedule, params, actions, tags)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to update alert rule: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(rule)
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

// HandleDeleteAlertRule handles deleting an alert rule.
func HandleDeleteAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana delete alert rule handler")

		err = c.DeleteAlertRule(ctx, ruleID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to delete alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully deleted alert rule: %s", ruleID)),
			},
		}, nil
	}
}

// HandleEnableAlertRule handles enabling an alert rule.
func HandleEnableAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana enable alert rule handler")

		err = c.EnableAlertRule(ctx, ruleID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to enable alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully enabled alert rule: %s", ruleID)),
			},
		}, nil
	}
}

// HandleDisableAlertRule handles disabling an alert rule.
func HandleDisableAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana disable alert rule handler")

		err = c.DisableAlertRule(ctx, ruleID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to disable alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully disabled alert rule: %s", ruleID)),
			},
		}, nil
	}
}

// HandleMuteAlertRule handles muting an alert rule.
func HandleMuteAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		duration := getOptionalStringParam(req, "duration")
		if duration == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent("duration is required"),
				},
			}, nil
		}

		logrus.WithFields(logrus.Fields{
			"rule_id":  ruleID,
			"duration": duration,
		}).Debug("Executing Kibana mute alert rule handler")

		err = c.MuteAlertRule(ctx, ruleID, duration)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to mute alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully muted alert rule: %s for %s", ruleID, duration)),
			},
		}, nil
	}
}

// HandleUnmuteAlertRule handles unmuting an alert rule.
func HandleUnmuteAlertRule(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		logrus.WithField("rule_id", ruleID).Debug("Executing Kibana unmute alert rule handler")

		err = c.UnmuteAlertRule(ctx, ruleID)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to unmute alert rule: %v", err)),
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Successfully unmuted alert rule: %s", ruleID)),
			},
		}, nil
	}
}

// HandleGetAlertRuleTypes handles listing available alert rule types.
func HandleGetAlertRuleTypes(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logrus.Debug("Executing Kibana get alert rule types handler")

		ruleTypes, err := c.GetAlertRuleTypes(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get alert rule types: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(ruleTypes)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alert rule types: %v", err)),
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

// HandleGetAlertRuleHistory handles getting alert rule execution history.
func HandleGetAlertRuleHistory(c *client.Client) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ruleID, err := requireStringParam(req, "rule_id")
		if err != nil {
			return nil, err
		}

		page := getOptionalIntParam(req, "page", 1)
		perPage := getOptionalIntParam(req, "per_page", 20)

		logrus.WithFields(logrus.Fields{
			"rule_id": ruleID,
			"page":    page,
			"perPage": perPage,
		}).Debug("Executing Kibana get alert rule history handler")

		history, err := c.GetAlertRuleHistory(ctx, ruleID, page, perPage)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to get alert rule history: %v", err)),
				},
			}, nil
		}

		resultJSON, err := marshalIndentJSON(history)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					mcp.NewTextContent(fmt.Sprintf("Failed to format alert rule history: %v", err)),
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
