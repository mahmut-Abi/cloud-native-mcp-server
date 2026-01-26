package errors

import (
	"fmt"
)

// ToServiceError converts a standard error to a ServiceError
// If the error is already a ServiceError, it returns it as-is
// Otherwise, it wraps it in an InternalError
func ToServiceError(err error) *ServiceError {
	if err == nil {
		return nil
	}

	if se, ok := AsServiceError(err); ok {
		return se
	}

	return InternalError(err)
}

// ToServiceErrorWithCode converts a standard error to a ServiceError with a specific code
func ToServiceErrorWithCode(err error, code, message string) *ServiceError {
	if err == nil {
		return nil
	}

	if se, ok := AsServiceError(err); ok {
		return se
	}

	return Wrap(err, code, message)
}

// MustConvert converts an error to ServiceError, returns error if error is nil
// This function is deprecated in favor of ToServiceError which handles nil errors gracefully
func MustConvert(err error) (*ServiceError, error) {
	if err == nil {
		return nil, fmt.Errorf("cannot convert nil error to ServiceError")
	}
	return ToServiceError(err), nil
}

// CheckAndConvert checks if error is nil, if not converts to ServiceError
func CheckAndConvert(err error) error {
	if err == nil {
		return nil
	}
	return ToServiceError(err)
}

// FormatError formats an error message with context
func FormatError(code, message string, args ...interface{}) *ServiceError {
	msg := message
	if len(args) > 0 {
		msg = fmt.Sprintf(message, args...)
	}
	return New(code, msg)
}
