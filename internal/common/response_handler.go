package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mahmut-Abi/k8s-mcp-server/internal/errors"
	"github.com/sirupsen/logrus"
)

// ResponseHandler provides common HTTP response handling functionality
type ResponseHandler struct {
	serviceName string
	logger      *logrus.Entry
}

// NewResponseHandler creates a new response handler for a service
func NewResponseHandler(serviceName string) *ResponseHandler {
	return &ResponseHandler{
		serviceName: serviceName,
		logger:      logrus.WithField("component", serviceName + "-client"),
	}
}

// HandleResponse processes an HTTP response and returns the body
func (h *ResponseHandler) HandleResponse(resp *http.Response) ([]byte, error) {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, h.createError("INVALID_RESPONSE", "failed to read response body", err, 502)
	}

	if resp.StatusCode >= 400 {
		return nil, h.handleErrorResponse(resp, body)
	}

	return body, nil
}

// HandleJSONResponse processes an HTTP response and unmarshals JSON body into target
func (h *ResponseHandler) HandleJSONResponse(resp *http.Response, target interface{}) error {
	body, err := h.HandleResponse(resp)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		return h.createError("INVALID_JSON", "failed to unmarshal response body", err, 502)
	}

	return nil
}

// handleErrorResponse handles error responses based on status code
func (h *ResponseHandler) handleErrorResponse(resp *http.Response, body []byte) error {
	statusCode := resp.StatusCode
	bodyStr := string(body)

	switch statusCode {
	case 400:
		return h.createError("BAD_REQUEST", "bad request", fmt.Errorf(bodyStr), 400)
	case 401:
		return errors.New(h.prefixCode("UNAUTHORIZED"), "unauthorized access").
			WithHTTPStatus(401)
	case 403:
		return errors.New(h.prefixCode("FORBIDDEN"), "forbidden access").
			WithHTTPStatus(403)
	case 404:
		return errors.NotFoundError("resource").
			WithHTTPStatus(404)
	case 409:
		return h.createError("CONFLICT", "resource conflict", fmt.Errorf(bodyStr), 409)
	case 422:
		return h.createError("UNPROCESSABLE_ENTITY", "unprocessable entity", fmt.Errorf(bodyStr), 422)
	case 429:
		return errors.New(h.prefixCode("RATE_LIMITED"), "rate limit exceeded").
			WithHTTPStatus(429)
	case 500:
		return h.createError("SERVER_ERROR", "server error", fmt.Errorf(bodyStr), 500)
	case 502:
		return h.createError("BAD_GATEWAY", "bad gateway", fmt.Errorf(bodyStr), 502)
	case 503:
		return h.createError("SERVICE_UNAVAILABLE", "service unavailable", fmt.Errorf(bodyStr), 503)
	case 504:
		return h.createError("GATEWAY_TIMEOUT", "gateway timeout", fmt.Errorf(bodyStr), 504)
	default:
		return errors.New(h.prefixCode("API_ERROR"), fmt.Sprintf("API error (status %d): %s", statusCode, bodyStr)).
			WithHTTPStatus(statusCode).
			WithContext("status_code", statusCode).
			WithContext("service", h.serviceName)
	}
}

// createError creates a standardized error with proper error code and context
func (h *ResponseHandler) createError(code, message string, err error, httpStatus int) error {
	if err != nil {
		return errors.Wrap(err, h.prefixCode(code), message).
			WithHTTPStatus(httpStatus).
			WithContext("service", h.serviceName)
	}
	return errors.New(h.prefixCode(code), message).
		WithHTTPStatus(httpStatus).
		WithContext("service", h.serviceName)
}

// prefixCode prefixes the error code with the service name
func (h *ResponseHandler) prefixCode(code string) string {
	return fmt.Sprintf("%s_%s", h.serviceName, code)
}

// HandleSuccessResponse handles successful responses and optionally unmarshals JSON
func (h *ResponseHandler) HandleSuccessResponse(resp *http.Response, target interface{}) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return h.handleErrorResponse(resp, nil)
	}

	if target != nil {
		return h.HandleJSONResponse(resp, target)
	}

	_, err := h.HandleResponse(resp)
	return err
}