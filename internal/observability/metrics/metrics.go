package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/sirupsen/logrus"
)

var (
	// Registry is the global Prometheus registry
	Registry = prometheus.NewRegistry()

	// BuildInfo provides build and version information
	BuildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "build_info",
			Help: "Build and version information",
		},
		[]string{"version", "commit", "go_version"},
	)

	// ServerInfo provides server runtime information
	ServerInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_info",
			Help: "Server runtime information",
		},
		[]string{"mode", "addr"},
	)

	// ServiceStatus tracks the status of each service
	ServiceStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_status",
			Help: "Status of each service (1=enabled, 0=disabled)",
		},
		[]string{"service_name"},
	)
)

// Init initializes the metrics system
func Init(version, commit, goVersion, mode, addr string) {
	// Register metrics (safe to call multiple times)
	// Use Register instead of MustRegister to avoid panic on duplicate registration
	// Duplicate registration errors are expected in tests and can be ignored
	_ = Registry.Register(BuildInfo)
	_ = Registry.Register(ServerInfo)
	_ = Registry.Register(ServiceStatus)
	_ = Registry.Register(HTTPRequestsTotal)
	_ = Registry.Register(HTTPRequestDuration)
	_ = Registry.Register(HTTPRequestSize)
	_ = Registry.Register(HTTPResponseSize)
	_ = Registry.Register(HTTPConnectionsActive)
	_ = Registry.Register(ToolCallsTotal)
	_ = Registry.Register(ToolCallDuration)
	_ = Registry.Register(ExternalAPICallsTotal)
	_ = Registry.Register(ExternalAPICallDuration)
	_ = Registry.Register(CacheHitsTotal)
	_ = Registry.Register(CacheMissesTotal)
	_ = Registry.Register(CircuitBreakerState)
	_ = Registry.Register(CircuitBreakerFailures)

	// Register default metrics (safe to call multiple times)
	_ = Registry.Register(collectors.NewGoCollector())
	_ = Registry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// Set build info (always 1 as it's a gauge)
	BuildInfo.WithLabelValues(version, commit, goVersion).Set(1)

	// Set server info (always 1 as it's gauge)
	ServerInfo.WithLabelValues(mode, addr).Set(1)

	logrus.WithFields(logrus.Fields{
		"component": "metrics",
		"version":   version,
		"commit":    commit,
	}).Info("Metrics system initialized")
}

// SetServiceStatus sets the status of a service
func SetServiceStatus(serviceName string, enabled bool) {
	value := 0.0
	if enabled {
		value = 1.0
	}
	ServiceStatus.WithLabelValues(serviceName).Set(value)
}
