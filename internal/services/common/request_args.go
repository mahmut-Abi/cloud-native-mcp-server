package common

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	optimize "github.com/mahmut-Abi/cloud-native-mcp-server/internal/util/performance"
)

// LookupArg returns the first matching argument value, checking common snake_case and camelCase aliases.
func LookupArg(args map[string]interface{}, keys ...string) (interface{}, bool) {
	if args == nil {
		return nil, false
	}

	seen := make(map[string]struct{})
	for _, key := range keys {
		for _, variant := range argKeyVariants(key) {
			if _, ok := seen[variant]; ok {
				continue
			}
			seen[variant] = struct{}{}

			if value, ok := args[variant]; ok && value != nil {
				return value, true
			}
		}
	}

	for _, key := range keys {
		normalizedKey := normalizeArgKey(key)
		if normalizedKey == "" {
			continue
		}

		for existingKey, value := range args {
			if value == nil {
				continue
			}
			if normalizeArgKey(existingKey) == normalizedKey {
				return value, true
			}
		}
	}

	return nil, false
}

// GetStringArg returns a string argument, coercing simple numeric and boolean values when needed.
func GetStringArg(args map[string]interface{}, keys ...string) (string, bool) {
	raw, ok := LookupArg(args, keys...)
	if !ok {
		return "", false
	}

	switch typed := raw.(type) {
	case string:
		value := strings.TrimSpace(typed)
		if value == "" {
			return "", false
		}
		return value, true
	case fmt.Stringer:
		value := strings.TrimSpace(typed.String())
		if value == "" {
			return "", false
		}
		return value, true
	case float64, float32, int, int32, int64, uint, uint32, uint64, bool:
		value := strings.TrimSpace(fmt.Sprintf("%v", typed))
		if value == "" {
			return "", false
		}
		return value, true
	default:
		return "", false
	}
}

// RequireStringArg returns a required string argument using the first key as the canonical name in errors.
func RequireStringArg(args map[string]interface{}, keys ...string) (string, error) {
	value, ok := GetStringArg(args, keys...)
	if ok {
		return value, nil
	}

	name := "argument"
	if len(keys) > 0 && keys[0] != "" {
		name = keys[0]
	}
	return "", fmt.Errorf("missing required parameter: %s", name)
}

// GetIntArg returns an integer argument or the provided default value.
func GetIntArg(args map[string]interface{}, defaultValue int, keys ...string) int {
	raw, ok := LookupArg(args, keys...)
	if !ok {
		return defaultValue
	}

	switch typed := raw.(type) {
	case float64:
		return int(typed)
	case float32:
		return int(typed)
	case int:
		return typed
	case int32:
		return int(typed)
	case int64:
		return int(typed)
	case uint:
		return int(typed)
	case uint32:
		return int(typed)
	case uint64:
		return int(typed)
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		if err == nil {
			return parsed
		}
	}

	return defaultValue
}

// GetBoolArg returns a boolean argument, accepting native booleans and common string forms.
func GetBoolArg(args map[string]interface{}, keys ...string) (*bool, error) {
	raw, ok := LookupArg(args, keys...)
	if !ok {
		return nil, nil
	}

	switch typed := raw.(type) {
	case bool:
		return &typed, nil
	case string:
		switch strings.ToLower(strings.TrimSpace(typed)) {
		case "true", "1", "yes", "y", "on":
			value := true
			return &value, nil
		case "false", "0", "no", "n", "off":
			value := false
			return &value, nil
		default:
			return nil, fmt.Errorf("invalid boolean value %q", typed)
		}
	default:
		return nil, fmt.Errorf("invalid boolean value type %T", raw)
	}
}

// GetObjectArg returns an object argument, accepting native objects and JSON strings.
func GetObjectArg(args map[string]interface{}, keys ...string) (map[string]interface{}, bool, error) {
	raw, ok := LookupArg(args, keys...)
	if !ok {
		return nil, false, nil
	}

	switch typed := raw.(type) {
	case map[string]interface{}:
		return typed, true, nil
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil, false, nil
		}
		var value map[string]interface{}
		if err := json.Unmarshal([]byte(typed), &value); err != nil {
			return nil, true, fmt.Errorf("failed to parse JSON object: %w", err)
		}
		return value, true, nil
	default:
		return nil, true, fmt.Errorf("expected JSON object, got %T", raw)
	}
}

// GetObjectSliceArg returns an array of objects, accepting native arrays and JSON strings.
func GetObjectSliceArg(args map[string]interface{}, keys ...string) ([]map[string]interface{}, bool, error) {
	raw, ok := LookupArg(args, keys...)
	if !ok {
		return nil, false, nil
	}

	switch typed := raw.(type) {
	case []map[string]interface{}:
		return typed, true, nil
	case []interface{}:
		result := make([]map[string]interface{}, 0, len(typed))
		for _, item := range typed {
			object, ok := item.(map[string]interface{})
			if !ok {
				return nil, true, fmt.Errorf("expected array of objects, got item %T", item)
			}
			result = append(result, object)
		}
		return result, true, nil
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil, false, nil
		}
		var value []map[string]interface{}
		if err := json.Unmarshal([]byte(typed), &value); err != nil {
			return nil, true, fmt.Errorf("failed to parse JSON object array: %w", err)
		}
		return value, true, nil
	default:
		return nil, true, fmt.Errorf("expected array of objects, got %T", raw)
	}
}

// GetStringSliceArg returns an array of strings, accepting arrays, JSON strings, CSV strings, and single strings.
func GetStringSliceArg(args map[string]interface{}, keys ...string) ([]string, bool, error) {
	raw, ok := LookupArg(args, keys...)
	if !ok {
		return nil, false, nil
	}

	switch typed := raw.(type) {
	case []string:
		return typed, true, nil
	case []interface{}:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			value, ok := item.(string)
			if !ok {
				return nil, true, fmt.Errorf("expected array of strings, got item %T", item)
			}
			value = strings.TrimSpace(value)
			if value != "" {
				result = append(result, value)
			}
		}
		return result, true, nil
	case string:
		value := strings.TrimSpace(typed)
		if value == "" {
			return nil, false, nil
		}

		if strings.HasPrefix(value, "[") {
			var parsed []string
			if err := json.Unmarshal([]byte(value), &parsed); err != nil {
				return nil, true, fmt.Errorf("failed to parse JSON string array: %w", err)
			}
			return parsed, true, nil
		}

		if strings.Contains(value, ",") {
			parts := strings.Split(value, ",")
			result := make([]string, 0, len(parts))
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part != "" {
					result = append(result, part)
				}
			}
			return result, true, nil
		}

		return []string{value}, true, nil
	default:
		return nil, true, fmt.Errorf("expected array of strings, got %T", raw)
	}
}

// GetJSONStringArg returns a JSON string argument, accepting native strings, objects, and arrays.
func GetJSONStringArg(args map[string]interface{}, keys ...string) (string, bool, error) {
	raw, ok := LookupArg(args, keys...)
	if !ok {
		return "", false, nil
	}

	switch typed := raw.(type) {
	case string:
		value := strings.TrimSpace(typed)
		if value == "" {
			return "", false, nil
		}
		return value, true, nil
	case map[string]interface{}, []interface{}:
		data, err := optimize.GlobalJSONPool.MarshalToBytes(typed)
		if err != nil {
			return "", true, fmt.Errorf("failed to serialize JSON argument: %w", err)
		}
		return string(data), true, nil
	default:
		return "", true, fmt.Errorf("expected JSON string, object, or array, got %T", raw)
	}
}

// GetRFC3339TimeArg returns an RFC3339 timestamp argument if present.
func GetRFC3339TimeArg(args map[string]interface{}, keys ...string) (*time.Time, error) {
	value, ok := GetStringArg(args, keys...)
	if !ok {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, fmt.Errorf("invalid RFC3339 timestamp %q: %w", value, err)
	}

	return &parsed, nil
}

func argKeyVariants(key string) []string {
	if key == "" {
		return nil
	}

	variants := []string{key}
	snake := toSnakeCase(key)
	camel := toLowerCamelCase(key)

	if snake != key {
		variants = append(variants, snake)
	}
	if camel != key {
		variants = append(variants, camel)
	}

	return variants
}

func toSnakeCase(input string) string {
	if input == "" {
		return input
	}

	var out []rune
	for i, r := range input {
		if r == '-' {
			out = append(out, '_')
			continue
		}

		if unicode.IsUpper(r) {
			if i > 0 && out[len(out)-1] != '_' {
				out = append(out, '_')
			}
			out = append(out, unicode.ToLower(r))
			continue
		}

		out = append(out, r)
	}

	return string(out)
}

func toLowerCamelCase(input string) string {
	if input == "" {
		return input
	}

	input = strings.ReplaceAll(input, "-", "_")
	if !strings.ContainsRune(input, '_') {
		return input
	}

	parts := strings.Split(input, "_")
	if len(parts) == 0 {
		return input
	}

	first := strings.ToLower(parts[0])
	var out strings.Builder
	out.WriteString(first)

	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		runes := []rune(strings.ToLower(part))
		runes[0] = unicode.ToUpper(runes[0])
		out.WriteString(string(runes))
	}

	return out.String()
}

func normalizeArgKey(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	input = strings.ReplaceAll(input, "_", "")
	input = strings.ReplaceAll(input, "-", "")
	return input
}
