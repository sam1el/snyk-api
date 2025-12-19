package client

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Sentinel errors for common API error conditions.
var (
	ErrAuthentication = errors.New("authentication failed")
	ErrNotFound       = errors.New("resource not found")
	ErrRateLimited    = errors.New("rate limited by server")
	ErrServerError    = errors.New("server error")
	ErrTimeout        = errors.New("request timeout")
	ErrMaxRetries     = errors.New("maximum retries exceeded")
)

// APIError provides detailed error information from API responses.
type APIError struct {
	StatusCode    int
	Message       string
	RequestID     string
	SnykRequestID string
	Retryable     bool
	RetryAfter    time.Duration
	Cause         error
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.SnykRequestID != "" {
		return fmt.Sprintf("API error %d: %s (snyk-request-id: %s)", e.StatusCode, e.Message, e.SnykRequestID)
	}
	if e.RequestID != "" {
		return fmt.Sprintf("API error %d: %s (request: %s)", e.StatusCode, e.Message, e.RequestID)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// Unwrap returns the underlying cause for errors.Is/As support.
func (e *APIError) Unwrap() error {
	return e.Cause
}

// NewAPIError creates a new APIError from an HTTP response.
func NewAPIError(statusCode int, message, requestID string, headers http.Header) *APIError {
	err := &APIError{
		StatusCode: statusCode,
		Message:    message,
		RequestID:  requestID,
		Retryable:  isRetryableStatus(statusCode),
	}

	// Capture snyk-request-id from response headers
	if headers != nil {
		err.SnykRequestID = headers.Get("snyk-request-id")
	}

	// Set appropriate cause based on status code
	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		err.Cause = ErrAuthentication
	case http.StatusNotFound:
		err.Cause = ErrNotFound
	case http.StatusTooManyRequests:
		err.Cause = ErrRateLimited
		err.Retryable = true
		if headers != nil {
			err.RetryAfter = parseRetryAfter(headers.Get("Retry-After"))
		}
	case http.StatusInternalServerError, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		err.Cause = ErrServerError
		err.Retryable = true
	}

	return err
}

// parseRetryAfter parses the Retry-After header value.
// Supports both integer seconds and HTTP date formats.
func parseRetryAfter(header string) time.Duration {
	if header == "" {
		return 0
	}

	// Try parsing as integer seconds
	if seconds, err := strconv.Atoi(header); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Try parsing as HTTP date (RFC 1123)
	if t, err := http.ParseTime(header); err == nil {
		duration := time.Until(t)
		if duration > 0 {
			return duration
		}
	}

	return 0
}

// isRetryableStatus returns true if the HTTP status code indicates
// the request can be retried.
func isRetryableStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// IsRetryable returns true if the error can be retried.
func IsRetryable(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Retryable
	}
	return false
}

// IsNotFound returns true if the error is a 404 Not Found.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAuthError returns true if the error is an authentication error.
func IsAuthError(err error) bool {
	return errors.Is(err, ErrAuthentication)
}
