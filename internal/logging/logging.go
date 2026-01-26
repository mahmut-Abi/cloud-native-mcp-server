package logging

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// LogLevel log level type
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// LogOperation standard operation type
const (
	OpInit       = "init"
	OpStart      = "start"
	OpStop       = "stop"
	OpCreate     = "create"
	OpUpdate     = "update"
	OpDelete     = "delete"
	OpGet        = "get"
	OpList       = "list"
	OpExecute    = "execute"
	OpQuery      = "query"
	OpConnect    = "connect"
	OpDisconnect = "disconnect"
	OpTimeout    = "timeout"
	OpRetry      = "retry"
	OpFail       = "fail"
	OpSuccess    = "success"
)

// LogConfig log configuration structure
type LogConfig struct {
	Output          io.Writer
	Level           logrus.Level
	UseJSONFormat   bool
	TimestampFormat string
	EnableColors    bool
}

// ResourceContext represents resource operation context
type ResourceContext struct {
	Component string
	Kind      string
	Name      string
	Namespace string
	Action    string
}

// DefaultConfig returns default log configuration
func DefaultConfig() *LogConfig {
	return &LogConfig{
		Output:          os.Stdout,
		Level:           logrus.InfoLevel,
		UseJSONFormat:   false,
		TimestampFormat: time.RFC3339,
		EnableColors:    true,
	}
}

// InitStdoutLogger initializes logrus log configuration
func InitStdoutLogger() {
	config := DefaultConfig()
	InitLoggerWithConfig(config)
}

// InitLoggerWithConfig initializes logger with configuration
func InitLoggerWithConfig(config *LogConfig) {
	logrus.SetOutput(config.Output)
	logrus.SetLevel(config.Level)

	if config.UseJSONFormat {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: config.TimestampFormat,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: config.TimestampFormat,
			FullTimestamp:   true,
			ForceColors:     config.EnableColors,
			PadLevelText:    true,
		})
	}
}

// SetLogLevel sets log level
func SetLogLevel(level logrus.Level) {
	logrus.SetLevel(level)
}

// EnableJSONFormat enables JSON format output
func EnableJSONFormat() {
	config := &LogConfig{
		Output:          os.Stdout,
		Level:           logrus.GetLevel(),
		UseJSONFormat:   true,
		TimestampFormat: time.RFC3339,
	}
	InitLoggerWithConfig(config)
}

// LogResourceError logs resource operation errors
func LogResourceError(ctx ResourceContext, err error) {
	if err == nil {
		return
	}

	fields := logrus.Fields{
		"component": ctx.Component,
		"kind":      ctx.Kind,
		"action":    ctx.Action,
	}

	if ctx.Name != "" {
		fields["name"] = ctx.Name
	}
	if ctx.Namespace != "" {
		fields["namespace"] = ctx.Namespace
	}

	logrus.WithFields(fields).WithError(err).Error("Resource operation failed")
}

// WrapResourceError wraps resource operation errors
func WrapResourceError(ctx ResourceContext, err error) error {
	if err == nil {
		return nil
	}

	msg := buildErrorMessage(ctx, err)
	return errors.New(msg)
}

// buildErrorMessage constructs error message
func buildErrorMessage(ctx ResourceContext, err error) string {
	if ctx.Namespace != "" && ctx.Name != "" {
		return fmt.Sprintf("failed to %s %s %s/%s: %v",
			ctx.Action, ctx.Kind, ctx.Namespace, ctx.Name, err)
	}
	if ctx.Name != "" {
		return fmt.Sprintf("failed to %s %s %s: %v",
			ctx.Action, ctx.Kind, ctx.Name, err)
	}
	return fmt.Sprintf("failed to %s %s: %v",
		ctx.Action, ctx.Kind, err)
}

// LogAndWrapResourceError logs and wraps errors
func LogAndWrapResourceError(ctx ResourceContext, err error) error {
	if err == nil {
		return nil
	}
	LogResourceError(ctx, err)
	return WrapResourceError(ctx, err)
}

// LogFieldsBuilder log field builder
type LogFieldsBuilder struct {
	fields logrus.Fields
}

// NewLogFieldsBuilder creates field builder
func NewLogFieldsBuilder() *LogFieldsBuilder {
	return &LogFieldsBuilder{
		fields: make(logrus.Fields),
	}
}

// Service adds service information
func (b *LogFieldsBuilder) Service(name, op string) *LogFieldsBuilder {
	b.fields["service"] = name
	b.fields["operation"] = op
	return b
}

// Resource adds resource information
func (b *LogFieldsBuilder) Resource(kind, name, ns string) *LogFieldsBuilder {
	b.fields["resource_kind"] = kind
	b.fields["resource_name"] = name
	if ns != "" {
		b.fields["namespace"] = ns
	}
	return b
}

// HTTP adds HTTP information
func (b *LogFieldsBuilder) HTTP(method, path, status string) *LogFieldsBuilder {
	b.fields["http_method"] = method
	b.fields["http_path"] = path
	b.fields["http_status"] = status
	return b
}

// Duration adds execution time
func (b *LogFieldsBuilder) Duration(ms int64) *LogFieldsBuilder {
	b.fields["duration_ms"] = ms
	return b
}

// Error adds error information
func (b *LogFieldsBuilder) Error(err error) *LogFieldsBuilder {
	if err != nil {
		b.fields["error"] = err.Error()
	}
	return b
}

// Status adds status
func (b *LogFieldsBuilder) Status(status string) *LogFieldsBuilder {
	b.fields["status"] = status
	return b
}

// Custom adds custom field
func (b *LogFieldsBuilder) Custom(key string, value interface{}) *LogFieldsBuilder {
	b.fields[key] = value
	return b
}

// Build builds log fields
func (b *LogFieldsBuilder) Build() logrus.Fields {
	return b.fields
}

// LogInfo logs Info level log
func (b *LogFieldsBuilder) LogInfo(msg string) {
	logrus.WithFields(b.fields).Info(msg)
}

// LogDebug logs Debug level log
func (b *LogFieldsBuilder) LogDebug(msg string) {
	logrus.WithFields(b.fields).Debug(msg)
}

// LogWarn logs Warn level log
func (b *LogFieldsBuilder) LogWarn(msg string) {
	logrus.WithFields(b.fields).Warn(msg)
}

// LogError logs Error level log
func (b *LogFieldsBuilder) LogError(msg string) {
	logrus.WithFields(b.fields).Error(msg)
}

// LogErrorWithErr logs error level log (with error information)
func (b *LogFieldsBuilder) LogErrorWithErr(msg string, err error) {
	if err != nil {
		b.fields["error"] = err.Error()
	}
	logrus.WithFields(b.fields).Error(msg)
}
