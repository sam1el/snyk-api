package client

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/sam1el/snyk-api/internal/ratelimit"
	cfgpkg "github.com/sam1el/snyk-api/pkg/config"
	"github.com/snyk/go-application-framework/pkg/workflow"
)

// Option is a functional option for configuring the Client.
type Option func(*config)

// config holds the internal configuration for a Client.
type config struct {
	engine      workflow.Engine
	rateLimit   ratelimit.Config
	logger      *zerolog.Logger
	userAgent   string
	baseURL     string
	restBaseURL string
	apiVersion  APIVersion
}

// defaultConfig returns the default client configuration.
func defaultConfig(res cfgpkg.Resolved) *config {
	nopLogger := zerolog.Nop()
	version := DefaultAPIVersion
	if parsed, err := ParseVersion(res.APIVersion); err == nil {
		version = parsed
	}
	baseURL := res.APIURL
	if baseURL == "" {
		baseURL = "https://api.snyk.io"
	}
	restBase := res.RestAPIURL
	if restBase == "" {
		restBase = "https://api.snyk.io/rest"
	}

	return &config{
		rateLimit:   ratelimit.DefaultConfig(),
		logger:      &nopLogger,
		userAgent:   "snyk-api-go/0.1.0",
		baseURL:     baseURL,
		restBaseURL: restBase,
		apiVersion:  version,
	}
}

// WithEngine sets the workflow engine to use.
// If not provided, a standalone engine will be created.
func WithEngine(engine workflow.Engine) Option {
	return func(c *config) {
		c.engine = engine
	}
}

// WithRateLimit sets rate limiting configuration.
func WithRateLimit(burstSize int, period time.Duration) Option {
	return func(c *config) {
		c.rateLimit.BurstSize = burstSize
		c.rateLimit.Period = period
	}
}

// WithRetryPolicy sets retry configuration.
func WithRetryPolicy(maxRetries int, baseDelay, maxDelay time.Duration) Option {
	return func(c *config) {
		c.rateLimit.MaxRetries = maxRetries
		c.rateLimit.RetryBaseDelay = baseDelay
		c.rateLimit.RetryMaxDelay = maxDelay
	}
}

// WithLogger sets the logger for the client.
func WithLogger(logger *zerolog.Logger) Option {
	return func(c *config) {
		c.logger = logger
	}
}

// WithUserAgent sets a custom user agent prefix.
func WithUserAgent(ua string) Option {
	return func(c *config) {
		c.userAgent = ua
	}
}

// WithBaseURL sets custom API base URLs.
// This is useful for regional endpoints or testing.
func WithBaseURL(v1URL, restURL string) Option {
	return func(c *config) {
		if v1URL != "" {
			c.baseURL = v1URL
		}
		if restURL != "" {
			c.restBaseURL = restURL
		}
	}
}
