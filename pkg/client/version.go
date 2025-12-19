// Package client provides API versioning support for Snyk REST API.
package client

import "time"

// APIVersion represents a Snyk REST API version.
type APIVersion string

// DefaultAPIVersion is the latest stable GA version used by default.
const DefaultAPIVersion APIVersion = "2025-11-05"

// Common stable API versions.
const (
	Version20251105 APIVersion = "2025-11-05" // Latest GA (current spec)
	Version20240422 APIVersion = "2024-04-22" // go-application-framework default
	Version20231127 APIVersion = "2023-11-27" // Common stable version
)

// String returns the string representation of the API version.
func (v APIVersion) String() string {
	return string(v)
}

// WithStability adds a stability suffix to the version.
func (v APIVersion) WithStability(stability string) APIVersion {
	return APIVersion(string(v) + "~" + stability)
}

// Beta returns the beta version of this API version.
func (v APIVersion) Beta() APIVersion {
	return v.WithStability("beta")
}

// Experimental returns the experimental version of this API version.
func (v APIVersion) Experimental() APIVersion {
	return v.WithStability("experimental")
}

// WithVersion sets a specific API version for requests.
// If not specified, DefaultAPIVersion is used.
//
// Example:
//
//	client.WithVersion(client.Version20240422)
//	client.WithVersion(client.DefaultAPIVersion.Beta())
//	client.WithVersion("2024-01-01")
func WithVersion(version APIVersion) Option {
	return func(c *config) {
		c.apiVersion = version
	}
}

// ParseVersion parses a string into an APIVersion, validating the format.
// Returns an error if the format is invalid.
func ParseVersion(s string) (APIVersion, error) {
	// Basic validation: YYYY-MM-DD format
	if len(s) < 10 {
		return "", ErrInvalidVersion
	}

	// Try to parse the date part
	dateStr := s[:10]
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		return "", ErrInvalidVersion
	}

	// Check for optional stability suffix
	if len(s) > 10 {
		if s[10] != '~' {
			return "", ErrInvalidVersion
		}
		stability := s[11:]
		if stability != "experimental" && stability != "beta" && stability != "wip" {
			return "", ErrInvalidVersion
		}
	}

	return APIVersion(s), nil
}

// ErrInvalidVersion is returned when an API version string is invalid.
var ErrInvalidVersion = &APIError{
	Message: "invalid API version format (expected YYYY-MM-DD or YYYY-MM-DD~stability)",
}
