package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ConfigLoader handles loading configuration from files
type ConfigLoader struct {
	validator *ConfigValidator
	parser    *EnvParser
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		validator: NewConfigValidator(),
		parser:    NewEnvParser(),
	}
}

// Load loads configuration from YAML file (if provided) and merges environment overrides.
// It also validates the configuration before returning it.
func (l *ConfigLoader) Load(path string) (*AppConfig, error) {
	cfg := &AppConfig{}

	// Load from file if path is provided
	if path != "" {
		if err := l.loadFromFile(path, cfg); err != nil {
			return nil, err
		}
	}

	// Apply environment variable overrides
	l.parser.Parse(cfg)

	// Validate configuration
	if err := l.validator.Validate(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// loadFromFile loads configuration from a YAML file
func (l *ConfigLoader) loadFromFile(path string, cfg *AppConfig) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file failed: %w", err)
	}

	if err := yaml.Unmarshal(b, cfg); err != nil {
		return fmt.Errorf("yaml unmarshal failed: %w", err)
	}

	return nil
}
