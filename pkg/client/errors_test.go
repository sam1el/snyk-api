package client

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		message        string
		requestID      string
		headers        http.Header
		wantRetryable  bool
		wantCause      error
		wantRetryAfter time.Duration
	}{
		{
			name:          "401 unauthorized",
			statusCode:    http.StatusUnauthorized,
			message:       "Invalid credentials",
			wantRetryable: false,
			wantCause:     ErrAuthentication,
		},
		{
			name:          "404 not found",
			statusCode:    http.StatusNotFound,
			message:       "Resource not found",
			wantRetryable: false,
			wantCause:     ErrNotFound,
		},
		{
			name:           "429 rate limited",
			statusCode:     http.StatusTooManyRequests,
			message:        "Rate limit exceeded",
			headers:        http.Header{"Retry-After": []string{"60"}},
			wantRetryable:  true,
			wantCause:      ErrRateLimited,
			wantRetryAfter: 60 * time.Second,
		},
		{
			name:          "500 server error",
			statusCode:    http.StatusInternalServerError,
			message:       "Internal server error",
			wantRetryable: true,
			wantCause:     ErrServerError,
		},
		{
			name:          "503 service unavailable",
			statusCode:    http.StatusServiceUnavailable,
			message:       "Service temporarily unavailable",
			wantRetryable: true,
			wantCause:     ErrServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.statusCode, tt.message, tt.requestID, tt.headers)

			assert.Equal(t, tt.statusCode, err.StatusCode)
			assert.Equal(t, tt.message, err.Message)
			assert.Equal(t, tt.wantRetryable, err.Retryable)
			assert.Equal(t, tt.wantRetryAfter, err.RetryAfter)

			if tt.wantCause != nil {
				assert.True(t, errors.Is(err, tt.wantCause))
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *APIError
		wantMsg string
	}{
		{
			name: "with snyk-request-id",
			err: &APIError{
				StatusCode:    500,
				Message:       "Server error",
				SnykRequestID: "abc-123",
			},
			wantMsg: "API error 500: Server error (snyk-request-id: abc-123)",
		},
		{
			name: "with request-id",
			err: &APIError{
				StatusCode: 404,
				Message:    "Not found",
				RequestID:  "req-456",
			},
			wantMsg: "API error 404: Not found (request: req-456)",
		},
		{
			name: "without request ids",
			err: &APIError{
				StatusCode: 429,
				Message:    "Rate limited",
			},
			wantMsg: "API error 429: Rate limited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantMsg, tt.err.Error())
		})
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "retryable api error",
			err:  &APIError{StatusCode: 503, Retryable: true},
			want: true,
		},
		{
			name: "non-retryable api error",
			err:  &APIError{StatusCode: 404, Retryable: false},
			want: false,
		},
		{
			name: "non-api error",
			err:  errors.New("network error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsRetryable(tt.err))
		})
	}
}

func TestParseRetryAfter(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   time.Duration
	}{
		{
			name:   "integer seconds",
			header: "120",
			want:   120 * time.Second,
		},
		{
			name:   "empty header",
			header: "",
			want:   0,
		},
		{
			name:   "invalid format",
			header: "invalid",
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseRetryAfter(tt.header)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsRetryableStatus(t *testing.T) {
	tests := []struct {
		status int
		want   bool
	}{
		{http.StatusOK, false},
		{http.StatusBadRequest, false},
		{http.StatusUnauthorized, false},
		{http.StatusNotFound, false},
		{http.StatusRequestTimeout, true},
		{http.StatusTooManyRequests, true},
		{http.StatusInternalServerError, true},
		{http.StatusBadGateway, true},
		{http.StatusServiceUnavailable, true},
		{http.StatusGatewayTimeout, true},
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.want, isRetryableStatus(tt.status))
		})
	}
}
