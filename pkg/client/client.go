// Package client provides a comprehensive API client for Snyk with rate limiting,
// retry logic, and integration with go-application-framework.
package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/sam1el/snyk-api/internal/ratelimit"
	"github.com/snyk/go-application-framework/pkg/app"
	"github.com/snyk/go-application-framework/pkg/configuration"
	"github.com/snyk/go-application-framework/pkg/networking"
	"github.com/snyk/go-application-framework/pkg/workflow"
)

// Client is the main API client that handles authentication, rate limiting,
// and request execution using go-application-framework.
type Client struct {
	config        *config
	engine        workflow.Engine
	networkAccess networking.NetworkAccess
	limiter       *ratelimit.Limiter
	logger        *zerolog.Logger
}

// New creates a new API client with the given options.
func New(ctx context.Context, opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	c := &Client{
		config: cfg,
		logger: cfg.logger,
	}

	// Setup workflow engine
	if cfg.engine == nil {
		c.engine = app.CreateAppEngine()
	} else {
		c.engine = cfg.engine
	}

	// Get network access from engine
	c.networkAccess = c.engine.GetNetworkAccess()

	// Setup rate limiter
	c.limiter = ratelimit.New(cfg.rateLimit)
	c.limiter.Start(ctx, cfg.rateLimit.BurstSize)

	c.logger.Debug().
		Str("base_url", cfg.baseURL).
		Str("rest_base_url", cfg.restBaseURL).
		Int("burst_size", cfg.rateLimit.BurstSize).
		Dur("period", cfg.rateLimit.Period).
		Msg("initialized snyk api client")

	return c, nil
}

// Close gracefully shuts down the client.
func (c *Client) Close() error {
	if c.limiter != nil {
		c.limiter.Stop()
	}
	return nil
}

// GetNetworkAccess returns the underlying network access for direct requests.
func (c *Client) GetNetworkAccess() networking.NetworkAccess {
	return c.networkAccess
}

// GetConfiguration returns the framework configuration.
func (c *Client) GetConfiguration() configuration.Configuration {
	return c.engine.GetConfiguration()
}

// GetLogger returns the client logger.
func (c *Client) GetLogger() *zerolog.Logger {
	return c.logger
}

// Execute performs an HTTP request with rate limiting and retries.
func (c *Client) Execute(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Add user agent
	req.Header.Set("User-Agent", c.config.userAgent)

	// Create result channel for rate-limited execution
	resultChan := make(chan error, 1)
	var response *http.Response

	rlReq := &ratelimit.Request{
		ID:  req.URL.String(),
		Ctx: ctx,
		Execute: func(execCtx context.Context) error {
			// Execute with retries
			// nolint:bodyclose // Response body is closed by caller of Execute()
			resp, err := c.executeWithRetries(execCtx, req)
			response = resp
			return err
		},
		Result: resultChan,
	}

	// Enqueue for rate-limited execution
	c.limiter.Enqueue(rlReq)

	// Wait for result
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-resultChan:
		if err != nil {
			return response, err
		}
		return response, nil
	}
}

// executeWithRetries performs the HTTP request with exponential backoff retries.
func (c *Client) executeWithRetries(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error
	maxRetries := c.config.rateLimit.MaxRetries
	baseDelay := c.config.rateLimit.RetryBaseDelay
	maxDelay := c.config.rateLimit.RetryMaxDelay

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Clone request to ensure it's fresh for each attempt
		reqClone := req.Clone(ctx)

		// Add authentication and default headers
		if err := c.networkAccess.AddHeaders(reqClone); err != nil {
			return nil, fmt.Errorf("failed to add headers: %w", err)
		}

		// Execute request
		resp, err := c.networkAccess.GetHttpClient().Do(reqClone)
		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				delay := calculateBackoff(attempt, baseDelay, maxDelay)
				c.logger.Debug().
					Err(err).
					Int("attempt", attempt+1).
					Dur("delay", delay).
					Msg("request failed, retrying")
				time.Sleep(delay)
				continue
			}
			break
		}

		// Check response status
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}

		// Create API error
		apiErr := NewAPIError(resp.StatusCode, resp.Status, "", resp.Header)
		lastErr = apiErr

		// Close response body since we're not returning it
		// Error is intentionally ignored as we're already handling a higher-level error
		_ = resp.Body.Close() //nolint:errcheck // Intentional: cleanup on error path

		// Don't retry if not retryable
		if !apiErr.Retryable || attempt >= maxRetries {
			return nil, apiErr
		}

		// Use server-specified retry delay if provided
		delay := apiErr.RetryAfter
		if delay == 0 {
			delay = calculateBackoff(attempt, baseDelay, maxDelay)
		}

		c.logger.Debug().
			Int("status", resp.StatusCode).
			Int("attempt", attempt+1).
			Dur("delay", delay).
			Str("snyk_request_id", apiErr.SnykRequestID).
			Msg("retrying request")

		time.Sleep(delay)
	}

	if lastErr == nil {
		lastErr = ErrMaxRetries
	}

	return nil, fmt.Errorf("%w: %w", ErrMaxRetries, lastErr)
}

// calculateBackoff returns the delay duration for exponential backoff with jitter.
func calculateBackoff(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	delay := baseDelay * (1 << uint(attempt))
	if delay > maxDelay {
		delay = maxDelay
	}

	// Add jitter: 80-100% of calculated delay
	jitter := time.Duration(float64(delay) * 0.8)
	return jitter
}

// BaseURL returns the v1 API base URL.
func (c *Client) BaseURL() string {
	return c.config.baseURL
}

// RestBaseURL returns the REST API base URL.
func (c *Client) RestBaseURL() string {
	return c.config.restBaseURL
}

// APIVersion returns the configured API version.
func (c *Client) APIVersion() APIVersion {
	return c.config.apiVersion
}

// HTTPClient returns the underlying HTTP client for direct requests.
func (c *Client) HTTPClient() *http.Client {
	return c.networkAccess.GetHttpClient()
}

// Token returns the authentication token.
func (c *Client) Token() string {
	cfg := c.engine.GetConfiguration()
	token := cfg.GetString(configuration.AUTHENTICATION_TOKEN)
	return token
}
