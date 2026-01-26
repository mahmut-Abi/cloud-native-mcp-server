package errors

import (
	"errors"
	"fmt"
)

type ServiceError struct {
	code       string
	message    string
	orig       error
	HTTPStatus int
	context    map[string]interface{}
}

const (
	ErrCodeMissingParam = "MISSING_PARAM"
	ErrCodeInvalidParam = "INVALID_PARAM"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeTimeout      = "TIMEOUT"
)

func New(code, message string) *ServiceError {
	return &ServiceError{
		code:       code,
		message:    message,
		HTTPStatus: 500,
		context:    make(map[string]interface{}),
	}
}

func Wrap(err error, code, message string) *ServiceError {
	return &ServiceError{
		orig:       err,
		code:       code,
		message:    message,
		HTTPStatus: 500,
		context:    make(map[string]interface{}),
	}
}

func (e *ServiceError) WithHTTPStatus(status int) *ServiceError {
	e.HTTPStatus = status
	return e
}

func (e *ServiceError) WithContext(key string, value interface{}) *ServiceError {
	e.context[key] = value
	return e
}

func (e *ServiceError) Code() string {
	return e.code
}

func (e *ServiceError) Message() string {
	return e.message
}

func (e *ServiceError) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("[%s] %s: %v", e.code, e.message, e.orig)
	}
	return fmt.Sprintf("[%s] %s", e.code, e.message)
}

func (e *ServiceError) Unwrap() error {
	return e.orig
}

func (e *ServiceError) Context() map[string]interface{} {
	return e.context
}

func MissingParamError(param string) *ServiceError {
	return New(ErrCodeMissingParam, fmt.Sprintf("missing parameter: %s", param)).WithHTTPStatus(400)
}

func InvalidParamError(param, reason string) *ServiceError {
	return New(ErrCodeInvalidParam, fmt.Sprintf("invalid %s: %s", param, reason)).WithHTTPStatus(400)
}

func NotFoundError(resource string) *ServiceError {
	return New(ErrCodeNotFound, fmt.Sprintf("%s not found", resource)).WithHTTPStatus(404)
}

func InternalError(err error) *ServiceError {
	return Wrap(err, ErrCodeInternal, "internal error").WithHTTPStatus(500)
}

func TimeoutError(operation string) *ServiceError {
	return New(ErrCodeTimeout, fmt.Sprintf("timeout: %s", operation)).WithHTTPStatus(504)
}

func Is(err error, code string) bool {
	var se *ServiceError
	return errors.As(err, &se) && se.code == code
}

func AsServiceError(err error) (*ServiceError, bool) {
	var se *ServiceError
	ok := errors.As(err, &se)
	return se, ok
}
