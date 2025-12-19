package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIVersion_String(t *testing.T) {
	v := APIVersion("2025-11-05")
	assert.Equal(t, "2025-11-05", v.String())
}

func TestAPIVersion_WithStability(t *testing.T) {
	tests := []struct {
		name      string
		version   APIVersion
		stability string
		want      string
	}{
		{
			name:      "beta",
			version:   "2025-11-05",
			stability: "beta",
			want:      "2025-11-05~beta",
		},
		{
			name:      "experimental",
			version:   "2025-11-05",
			stability: "experimental",
			want:      "2025-11-05~experimental",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.version.WithStability(tt.stability)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestAPIVersion_Beta(t *testing.T) {
	v := APIVersion("2025-11-05")
	beta := v.Beta()
	assert.Equal(t, "2025-11-05~beta", beta.String())
}

func TestAPIVersion_Experimental(t *testing.T) {
	v := APIVersion("2025-11-05")
	exp := v.Experimental()
	assert.Equal(t, "2025-11-05~experimental", exp.String())
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    APIVersion
		wantErr bool
	}{
		{
			name:  "valid GA version",
			input: "2025-11-05",
			want:  "2025-11-05",
		},
		{
			name:  "valid beta version",
			input: "2025-11-05~beta",
			want:  "2025-11-05~beta",
		},
		{
			name:  "valid experimental version",
			input: "2025-11-05~experimental",
			want:  "2025-11-05~experimental",
		},
		{
			name:    "invalid format - too short",
			input:   "2025-11",
			wantErr: true,
		},
		{
			name:    "invalid format - bad date",
			input:   "2025-13-05",
			wantErr: true,
		},
		{
			name:    "invalid stability suffix",
			input:   "2025-11-05~invalid",
			wantErr: true,
		},
		{
			name:    "invalid separator",
			input:   "2025-11-05-beta",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, ErrInvalidVersion, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDefaultAPIVersion(t *testing.T) {
	// Ensure default version is in valid format
	_, err := ParseVersion(DefaultAPIVersion.String())
	require.NoError(t, err)

	// Ensure it's a GA version (no stability suffix)
	assert.Equal(t, 10, len(DefaultAPIVersion.String()))
}
