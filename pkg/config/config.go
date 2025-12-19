// Package config provides configuration loading and management for the Snyk API client,
// integrating with go-application-framework's configuration system.
package config

import (
	"time"

	"github.com/snyk/go-application-framework/pkg/configuration"
)

// Config keys used in go-application-framework configuration.
const (
	// API endpoint configuration
	KeyAPIURL     = "api"
	KeyRESTAPIURL = "rest_api_url"

	// Rate limiting configuration
	KeyRateLimitBurst  = "api_rate_limit_burst"
	KeyRateLimitPeriod = "api_rate_limit_period"

	// Retry configuration
	KeyMaxRetries     = "api_max_retries"
	KeyRetryBaseDelay = "api_retry_base_delay"
	KeyRetryMaxDelay  = "api_retry_max_delay"

	// Request configuration
	KeyUserAgent = "api_user_agent"
)

// APIConfig holds API client configuration loaded from the framework.
type APIConfig struct {
	BaseURL         string
	RestBaseURL     string
	RateLimitBurst  int
	RateLimitPeriod time.Duration
	MaxRetries      int
	RetryBaseDelay  time.Duration
	RetryMaxDelay   time.Duration
	UserAgent       string
}

// LoadFromConfiguration extracts API configuration from go-application-framework configuration.
func LoadFromConfiguration(cfg configuration.Configuration) *APIConfig {
	apiCfg := &APIConfig{
		BaseURL:         cfg.GetString(KeyAPIURL),
		RestBaseURL:     cfg.GetString(KeyRESTAPIURL),
		RateLimitBurst:  cfg.GetInt(KeyRateLimitBurst),
		RateLimitPeriod: time.Duration(cfg.GetInt(KeyRateLimitPeriod)) * time.Millisecond,
		MaxRetries:      cfg.GetInt(KeyMaxRetries),
		RetryBaseDelay:  time.Duration(cfg.GetInt(KeyRetryBaseDelay)) * time.Millisecond,
		RetryMaxDelay:   time.Duration(cfg.GetInt(KeyRetryMaxDelay)) * time.Second,
		UserAgent:       cfg.GetString(KeyUserAgent),
	}

	// Apply defaults for missing values
	if apiCfg.BaseURL == "" {
		apiCfg.BaseURL = "https://api.snyk.io"
	}
	if apiCfg.RestBaseURL == "" {
		apiCfg.RestBaseURL = "https://api.snyk.io/rest"
	}
	if apiCfg.RateLimitBurst == 0 {
		apiCfg.RateLimitBurst = 10
	}
	if apiCfg.RateLimitPeriod == 0 {
		apiCfg.RateLimitPeriod = 500 * time.Millisecond
	}
	if apiCfg.MaxRetries == 0 {
		apiCfg.MaxRetries = 5
	}
	if apiCfg.RetryBaseDelay == 0 {
		apiCfg.RetryBaseDelay = 100 * time.Millisecond
	}
	if apiCfg.RetryMaxDelay == 0 {
		apiCfg.RetryMaxDelay = 30 * time.Second
	}
	if apiCfg.UserAgent == "" {
		apiCfg.UserAgent = "snyk-api-go/0.1.0"
	}

	return apiCfg
}
