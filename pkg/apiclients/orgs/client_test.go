package orgs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationType(t *testing.T) {
	// Test that the generated constants are correct
	assert.Equal(t, OrganizationType("organization"), OrganizationTypeOrganization)
}

func TestOrganizationStructure(t *testing.T) {
	// Test that Organization struct can be instantiated
	org := Organization{
		Type: OrganizationTypeOrganization,
	}

	require.NotNil(t, org)
	assert.Equal(t, OrganizationTypeOrganization, org.Type)
}

func TestListOrganizationsParams(t *testing.T) {
	limit := 10
	cursor := "test-cursor"

	params := &ListOrganizationsParams{
		Limit:         &limit,
		StartingAfter: &cursor,
	}

	require.NotNil(t, params)
	assert.Equal(t, 10, *params.Limit)
	assert.Equal(t, "test-cursor", *params.StartingAfter)
}
