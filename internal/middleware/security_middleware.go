package middleware

import (
	"net/http"
	"strings"
)

// SecurityConfig holds security middleware configuration
type SecurityConfig struct {
	Enabled bool

	// Content Security Policy
	CSPEnabled        bool
	CSPDefaultSrc     string
	CSPScriptSrc      string
	CSPStyleSrc       string
	CSPImgSrc         string
	CSPConnectSrc     string
	CSPFontSrc        string
	CSPObjectSrc      string
	CSPMediaSrc       string
	CSPFrameSrc       string
	CSPBaseURI        string
	CSPFormAction     string
	CSPFrameAncestors string
	CSPReportURI      string

	// Other security headers
	XFrameOptions           string // DENY, SAMEORIGIN, or ALLOW-FROM uri
	XContentTypeOptions     string // nosniff
	XSSProtection           string // 1; mode=block
	ReferrerPolicy          string // no-referrer, strict-origin-when-cross-origin, etc.
	PermissionsPolicy       string // geolocation=(), camera=(), microphone=(), etc.
	StrictTransportSecurity string // max-age=31536000; includeSubDomains; preload

	// Feature policies
	CrossOriginEmbedderPolicy string // require-corp, credentialless
	CrossOriginOpenerPolicy   string // same-origin, same-origin-allow-popups
	CrossOriginResourcePolicy string // same-site, same-origin, cross-origin
}

// DefaultSecurityConfig returns a secure default configuration
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		Enabled: true,

		CSPEnabled:        true,
		CSPDefaultSrc:     "'self'",
		CSPScriptSrc:      "'self' 'unsafe-inline' 'unsafe-eval'",
		CSPStyleSrc:       "'self' 'unsafe-inline'",
		CSPImgSrc:         "'self' data: https:",
		CSPConnectSrc:     "'self'",
		CSPFontSrc:        "'self' data:",
		CSPObjectSrc:      "'none'",
		CSPMediaSrc:       "'self'",
		CSPFrameSrc:       "'none'",
		CSPBaseURI:        "'self'",
		CSPFormAction:     "'self'",
		CSPFrameAncestors: "'none'",

		XFrameOptions:           "DENY",
		XContentTypeOptions:     "nosniff",
		XSSProtection:           "1; mode=block",
		ReferrerPolicy:          "strict-origin-when-cross-origin",
		PermissionsPolicy:       "geolocation=(), camera=(), microphone=(), payment=()",
		StrictTransportSecurity: "max-age=31536000; includeSubDomains",

		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-site",
	}
}

// SecurityMiddleware creates a middleware that adds security headers to all responses
func SecurityMiddleware(config SecurityConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Wrap the response writer to add headers before sending
			sw := &securityResponseWriter{
				ResponseWriter: w,
				config:         config,
			}

			next.ServeHTTP(sw, r)
		})
	}
}

// securityResponseWriter wraps http.ResponseWriter to add security headers
type securityResponseWriter struct {
	http.ResponseWriter
	config     SecurityConfig
	headersSet bool
}

// WriteHeader implements http.ResponseWriter and adds security headers
func (sw *securityResponseWriter) WriteHeader(statusCode int) {
	if !sw.headersSet {
		sw.setSecurityHeaders()
		sw.headersSet = true
	}
	sw.ResponseWriter.WriteHeader(statusCode)
}

// Write implements http.ResponseWriter and ensures headers are set
func (sw *securityResponseWriter) Write(b []byte) (int, error) {
	if !sw.headersSet {
		sw.setSecurityHeaders()
		sw.headersSet = true
	}
	return sw.ResponseWriter.Write(b)
}

// setSecurityHeaders sets all configured security headers
func (sw *securityResponseWriter) setSecurityHeaders() {
	config := sw.config

	// Content Security Policy
	if config.CSPEnabled {
		var cspBuilder strings.Builder
		if config.CSPDefaultSrc != "" {
			cspBuilder.WriteString("default-src ")
			cspBuilder.WriteString(config.CSPDefaultSrc)
			cspBuilder.WriteString("; ")
		}
		if config.CSPScriptSrc != "" {
			cspBuilder.WriteString("script-src ")
			cspBuilder.WriteString(config.CSPScriptSrc)
			cspBuilder.WriteString("; ")
		}
		if config.CSPStyleSrc != "" {
			cspBuilder.WriteString("style-src ")
			cspBuilder.WriteString(config.CSPStyleSrc)
			cspBuilder.WriteString("; ")
		}
		if config.CSPImgSrc != "" {
			cspBuilder.WriteString("img-src ")
			cspBuilder.WriteString(config.CSPImgSrc)
			cspBuilder.WriteString("; ")
		}
		if config.CSPConnectSrc != "" {
			cspBuilder.WriteString("connect-src ")
			cspBuilder.WriteString(config.CSPConnectSrc)
			cspBuilder.WriteString("; ")
		}
		if config.CSPFontSrc != "" {
			cspBuilder.WriteString("font-src ")
			cspBuilder.WriteString(config.CSPFontSrc)
			cspBuilder.WriteString("; ")
		}
		if config.CSPObjectSrc != "" {
			cspBuilder.WriteString("object-src ")
			cspBuilder.WriteString(config.CSPObjectSrc)
			cspBuilder.WriteString("; ")
		}
		if config.CSPMediaSrc != "" {
			cspBuilder.WriteString("media-src ")
			cspBuilder.WriteString(config.CSPMediaSrc)
			cspBuilder.WriteString("; ")
		}
		if config.CSPFrameSrc != "" {
			cspBuilder.WriteString("frame-src ")
			cspBuilder.WriteString(config.CSPFrameSrc)
			cspBuilder.WriteString("; ")
		}
		if config.CSPBaseURI != "" {
			cspBuilder.WriteString("base-uri ")
			cspBuilder.WriteString(config.CSPBaseURI)
			cspBuilder.WriteString("; ")
		}
		if config.CSPFormAction != "" {
			cspBuilder.WriteString("form-action ")
			cspBuilder.WriteString(config.CSPFormAction)
			cspBuilder.WriteString("; ")
		}
		if config.CSPFrameAncestors != "" {
			cspBuilder.WriteString("frame-ancestors ")
			cspBuilder.WriteString(config.CSPFrameAncestors)
			cspBuilder.WriteString("; ")
		}
		if config.CSPReportURI != "" {
			cspBuilder.WriteString("report-uri ")
			cspBuilder.WriteString(config.CSPReportURI)
			cspBuilder.WriteString("; ")
		}
		if cspBuilder.Len() > 0 {
			// Remove trailing "; "
			csp := cspBuilder.String()
			sw.Header().Set("Content-Security-Policy", csp[:len(csp)-2])
		}
	}

	// X-Frame-Options
	if config.XFrameOptions != "" {
		sw.Header().Set("X-Frame-Options", config.XFrameOptions)
	}

	// X-Content-Type-Options
	if config.XContentTypeOptions != "" {
		sw.Header().Set("X-Content-Type-Options", config.XContentTypeOptions)
	}

	// X-XSS-Protection
	if config.XSSProtection != "" {
		sw.Header().Set("X-XSS-Protection", config.XSSProtection)
	}

	// Referrer-Policy
	if config.ReferrerPolicy != "" {
		sw.Header().Set("Referrer-Policy", config.ReferrerPolicy)
	}

	// Permissions-Policy
	if config.PermissionsPolicy != "" {
		sw.Header().Set("Permissions-Policy", config.PermissionsPolicy)
	}

	// Strict-Transport-Security (only for HTTPS)
	if config.StrictTransportSecurity != "" {
		sw.Header().Set("Strict-Transport-Security", config.StrictTransportSecurity)
	}

	// Cross-Origin-Embedder-Policy
	if config.CrossOriginEmbedderPolicy != "" {
		sw.Header().Set("Cross-Origin-Embedder-Policy", config.CrossOriginEmbedderPolicy)
	}

	// Cross-Origin-Opener-Policy
	if config.CrossOriginOpenerPolicy != "" {
		sw.Header().Set("Cross-Origin-Opener-Policy", config.CrossOriginOpenerPolicy)
	}

	// Cross-Origin-Resource-Policy
	if config.CrossOriginResourcePolicy != "" {
		sw.Header().Set("Cross-Origin-Resource-Policy", config.CrossOriginResourcePolicy)
	}

	// Additional security headers
	sw.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
	sw.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
}
