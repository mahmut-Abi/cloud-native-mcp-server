package middleware

import (
	"regexp"
	"strings"
)

// SensitiveDataMasker handles masking of sensitive data in audit logs
type SensitiveDataMasker struct {
	enabled    bool
	maskValue  string
	fieldNames []string
	patterns   []*regexp.Regexp
}

// DefaultMaskValue is the default value used to mask sensitive data
const DefaultMaskValue = "********"

// DefaultSensitiveFields are default field names that should be masked
var DefaultSensitiveFields = []string{
	"password",
	"passwd",
	"pwd",
	"secret",
	"token",
	"apikey",
	"api_key",
	"api-key",
	"apikey",
	"authorization",
	"auth",
	"credentials",
	"credential",
	"private_key",
	"private-key",
	"privatekey",
	"access_token",
	"access-token",
	"accesstoken",
	"refresh_token",
	"refresh-token",
	"refreshtoken",
	"session_token",
	"session-token",
	"sessiontoken",
	"bearer_token",
	"bearer-token",
	"bearertoken",
	"jwt",
	"oauth",
	"oauth_token",
	"oauth-token",
	"oauthtoken",
	"client_secret",
	"client-secret",
	"clientsecret",
	"api_secret",
	"api-secret",
	"apisecret",
	"signing_key",
	"signing-key",
	"signingkey",
	"encryption_key",
	"encryption-key",
	"encryptionkey",
	"decrypt_key",
	"decrypt-key",
	"decryptkey",
	"ssh_key",
	"ssh-key",
	"sshkey",
	"rsa_key",
	"rsa-key",
	"rsakey",
	"credit_card",
	"credit-card",
	"creditcard",
	"ssn",
	"social_security_number",
	"social-security-number",
	"pin",
	"cvv",
	"cvc",
	"account_number",
	"account-number",
	"accountnumber",
	"routing_number",
	"routing-number",
	"routingnumber",
}

// DefaultSensitivePatterns are regex patterns for detecting sensitive data
var DefaultSensitivePatterns = []*regexp.Regexp{
	// Bearer tokens (e.g., "Bearer eyJhbGciOiJIUzI1NiIs...")
	regexp.MustCompile(`(?i)bearer\s+[a-zA-Z0-9\-_=]+\.[a-zA-Z0-9\-_=]+\.[a-zA-Z0-9\-_=]+`),
	// API keys in format: key=xxx
	regexp.MustCompile(`(?i)(?:api[_-]?key|apikey|token|secret)[\s=:]+["']?[a-zA-Z0-9\-_]{16,}["']?`),
	// Basic auth credentials
	regexp.MustCompile(`(?i)basi[c]\s+[a-zA-Z0-9+/=]+`),
	// JWT tokens
	regexp.MustCompile(`[a-zA-Z0-9\-_=]+\.[a-zA-Z0-9\-_=]+\.[a-zA-Z0-9\-_=]+`),
	// Credit card numbers (basic pattern)
	regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`),
	// AWS access keys
	regexp.MustCompile(`(?i)AKIA[0-9A-Z]{16}`),
	// Generic sensitive strings (16+ alphanumeric chars)
	regexp.MustCompile(`[a-zA-Z0-9]{32,}`),
}

// NewSensitiveDataMasker creates a new masker with default settings
func NewSensitiveDataMasker() *SensitiveDataMasker {
	patterns := make([]*regexp.Regexp, len(DefaultSensitivePatterns))
	copy(patterns, DefaultSensitivePatterns)

	return &SensitiveDataMasker{
		enabled:    true,
		maskValue:  DefaultMaskValue,
		fieldNames: DefaultSensitiveFields,
		patterns:   patterns,
	}
}

// NewSensitiveDataMaskerWithConfig creates a new masker with custom configuration
func NewSensitiveDataMaskerWithConfig(enabled bool, maskValue string, fieldNames []string, patterns []string) (*SensitiveDataMasker, error) {
	compiledPatterns := make([]*regexp.Regexp, 0)
	for _, pattern := range patterns {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		compiledPatterns = append(compiledPatterns, regex)
	}

	return &SensitiveDataMasker{
		enabled:    enabled,
		maskValue:  maskValue,
		fieldNames: fieldNames,
		patterns:   compiledPatterns,
	}, nil
}

// MaskAuditEntry masks sensitive data in an audit log entry
func (m *SensitiveDataMasker) MaskAuditEntry(entry *AuditLogEntry) {
	if !m.enabled {
		return
	}

	// Mask input parameters
	if entry.InputParams != nil {
		entry.InputParams = m.maskMap(entry.InputParams)
	}

	// Mask output
	entry.Output = m.maskInterface(entry.Output)

	// Mask error message
	entry.ErrorMsg = m.maskString(entry.ErrorMsg)
}

// maskMap masks sensitive data in a map
func (m *SensitiveDataMasker) maskMap(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		if m.isSensitiveField(key) {
			result[key] = m.maskValue
		} else {
			result[key] = m.maskInterface(value)
		}
	}
	return result
}

// maskInterface masks sensitive data in an interface
func (m *SensitiveDataMasker) maskInterface(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return m.maskString(v)
	case map[string]interface{}:
		return m.maskMap(v)
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = m.maskInterface(item)
		}
		return result
	default:
		return value
	}
}

// maskString masks sensitive patterns in a string
func (m *SensitiveDataMasker) maskString(s string) string {
	if s == "" {
		return s
	}

	result := s
	for _, pattern := range m.patterns {
		result = pattern.ReplaceAllString(result, m.maskValue)
	}
	return result
}

// isSensitiveField checks if a field name is sensitive
func (m *SensitiveDataMasker) isSensitiveField(fieldName string) bool {
	lowerField := strings.ToLower(fieldName)
	for _, sensitive := range m.fieldNames {
		if strings.Contains(lowerField, strings.ToLower(sensitive)) {
			return true
		}
	}
	return false
}

// SetEnabled enables or disables masking
func (m *SensitiveDataMasker) SetEnabled(enabled bool) {
	m.enabled = enabled
}

// SetMaskValue sets the mask value
func (m *SensitiveDataMasker) SetMaskValue(maskValue string) {
	m.maskValue = maskValue
}

// AddSensitiveField adds a field name to the sensitive list
func (m *SensitiveDataMasker) AddSensitiveField(fieldName string) {
	m.fieldNames = append(m.fieldNames, fieldName)
}

// AddSensitivePattern adds a regex pattern to the sensitive patterns list
func (m *SensitiveDataMasker) AddSensitivePattern(pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	m.patterns = append(m.patterns, regex)
	return nil
}
