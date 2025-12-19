# Adding New API Domains

This guide explains how to add support for a new Snyk API domain (e.g., Targets, Issues, Ignores).

## Overview

The architecture uses a code generation pipeline:

```
OpenAPI Spec → oapi-codegen → Generated Client → Wrapper → CLI Commands
```

Adding a new domain requires 5 steps:

1. **Create OpenAPI specification**
2. **Configure code generation**
3. **Generate the client**
4. **Write wrapper client**
5. **Add CLI commands**

## Step 1: Create OpenAPI Specification

Create a minimal OpenAPI spec in `.github/api-ref/rest/<domain>-minimal.yaml`.

### Example: `targets-minimal.yaml`

```yaml
openapi: 3.0.3
info:
  title: Snyk Targets API (Minimal)
  version: 2024-12-19
  description: |
    Minimal OpenAPI spec for Snyk Targets API.
    Full spec: https://apidocs.snyk.io/

servers:
  - url: https://api.snyk.io/rest
    description: Snyk REST API

security:
  - BearerAuth: []

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: token

  schemas:
    Target:
      type: object
      required:
        - id
        - type
        - attributes
      properties:
        id:
          type: string
          format: uuid
        type:
          type: string
          enum: [target]
        attributes:
          type: object
          required:
            - display_name
            - url
          properties:
            display_name:
              type: string
            url:
              type: string
              format: uri
            origin:
              type: string

    TargetList:
      type: object
      required:
        - data
      properties:
        data:
          type: array
          items:
            $ref: '#/components/schemas/Target'
        links:
          type: object
          properties:
            next:
              type: string
              format: uri

    Error:
      type: object
      required:
        - errors
      properties:
        errors:
          type: array
          items:
            type: object
            properties:
              status:
                type: string
              detail:
                type: string

paths:
  /orgs/{org_id}/targets:
    get:
      operationId: listTargets
      summary: List targets
      parameters:
        - name: org_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
        - name: limit
          in: query
          schema:
            type: integer
            default: 10
      responses:
        '200':
          description: Target list
          content:
            application/vnd.api+json:
              schema:
                $ref: '#/components/schemas/TargetList'
        '401':
          description: Unauthorized
          content:
            application/vnd.api+json:
              schema:
                $ref: '#/components/schemas/Error'

  /orgs/{org_id}/targets/{target_id}:
    get:
      operationId: getTarget
      summary: Get target by ID
      parameters:
        - name: org_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
        - name: target_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Target details
          content:
            application/vnd.api+json:
              schema:
                type: object
                required:
                  - data
                properties:
                  data:
                    $ref: '#/components/schemas/Target'
```

### Tips for OpenAPI Specs

- ✅ Start minimal - add only the fields you need
- ✅ Use `operationId` for meaningful function names
- ✅ Include `application/vnd.api+json` content type
- ✅ Use UUID format for IDs
- ✅ Define error schemas
- ❌ Don't copy the full Snyk spec (too large)
- ❌ Don't forget security schemes

## Step 2: Configure Code Generation

Create `pkg/apiclients/<domain>/oapi-codegen.yaml`:

```yaml
package: targets
generate:
  client: true
  models: true
  embedded-spec: true
output: client_generated.go
```

Add generation directive to `pkg/apiclients/<domain>/doc.go`:

```go
// Package targets provides a client for the Snyk Targets API.
//
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config oapi-codegen.yaml ../../../.github/api-ref/rest/targets-minimal.yaml
package targets
```

## Step 3: Generate the Client

```bash
# Generate all clients
make generate

# Or generate specific package
go generate ./pkg/apiclients/targets/
```

This creates `client_generated.go` (~500-1500 lines).

## Step 4: Write Wrapper Client

Create `pkg/apiclients/<domain>/client.go` to wrap the generated client:

```go
package targets

import (
    "context"
    "fmt"
    "io"
    "net/http"

    "github.com/sam1el/snyk-api/pkg/client"
    openapi_types "github.com/oapi-codegen/runtime/types"
)

// TargetsClient wraps the generated OpenAPI client for the Targets API.
type TargetsClient struct {
    apiClient  *ClientWithResponses
    baseClient *client.Client
}

// NewTargetsClient creates a new TargetsClient.
func NewTargetsClient(baseClient *client.Client) *TargetsClient {
    // The generated client needs an http.Client that uses our baseClient's Execute method.
    httpClient := &http.Client{
        Transport: &roundTripperFunc{
            roundTrip: func(req *http.Request) (*http.Response, error) {
                return baseClient.Execute(req.Context(), req)
            },
        },
    }

    // Create the generated client
    apiClient, err := NewClientWithResponses(baseClient.RestBaseURL(), WithHTTPClient(httpClient))
    if err != nil {
        baseClient.GetLogger().Fatal().Err(err).Msg("Failed to create generated targets API client")
    }

    return &TargetsClient{
        apiClient:  apiClient,
        baseClient: baseClient,
    }
}

// ListTargets lists all targets for a given organization.
func (c *TargetsClient) ListTargets(ctx context.Context, orgID string, params *ListTargetsParams) (*TargetList, error) {
    uuidOrgID, err := openapi_types.ParseUUID(orgID)
    if err != nil {
        return nil, fmt.Errorf("invalid organization ID format: %w", err)
    }

    resp, err := c.apiClient.ListTargetsWithResponse(ctx, uuidOrgID, params)
    if err != nil {
        return nil, fmt.Errorf("failed to list targets for organization %s: %w", orgID, err)
    }

    // Handle response (check for ApplicationvndApiJSON200)
    if resp.ApplicationvndApiJSON200 == nil {
        if resp.Body != nil {
            defer func() { _ = resp.HTTPResponse.Body.Close() }()
            bodyBytes, _ := io.ReadAll(resp.HTTPResponse.Body)
            return nil, fmt.Errorf("unexpected response from API: status %s, body: %s", resp.Status(), string(bodyBytes))
        }
        return nil, fmt.Errorf("unexpected response from API: status %s", resp.Status())
    }

    return resp.ApplicationvndApiJSON200, nil
}

// GetTarget retrieves a single target by ID.
func (c *TargetsClient) GetTarget(ctx context.Context, orgID, targetID string) (*TargetResponse, error) {
    uuidOrgID, err := openapi_types.ParseUUID(orgID)
    if err != nil {
        return nil, fmt.Errorf("invalid organization ID format: %w", err)
    }
    uuidTargetID, err := openapi_types.ParseUUID(targetID)
    if err != nil {
        return nil, fmt.Errorf("invalid target ID format: %w", err)
    }

    resp, err := c.apiClient.GetTargetWithResponse(ctx, uuidOrgID, uuidTargetID)
    if err != nil {
        return nil, fmt.Errorf("failed to get target %s for organization %s: %w", targetID, orgID, err)
    }

    if resp.ApplicationvndApiJSON200 == nil {
        if resp.Body != nil {
            defer func() { _ = resp.HTTPResponse.Body.Close() }()
            bodyBytes, _ := io.ReadAll(resp.HTTPResponse.Body)
            return nil, fmt.Errorf("unexpected response from API: status %s, body: %s", resp.Status(), string(bodyBytes))
        }
        return nil, fmt.Errorf("unexpected response from API: status %s", resp.Status())
    }

    return resp.ApplicationvndApiJSON200, nil
}

// roundTripperFunc is a helper to allow a function to implement http.RoundTripper.
type roundTripperFunc struct {
    roundTrip func(*http.Request) (*http.Response, error)
}

func (rt *roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
    return rt.roundTrip(req)
}
```

### Wrapper Client Patterns

- ✅ Convert string IDs to UUIDs
- ✅ Use `baseClient.Execute()` for rate limiting/retry
- ✅ Handle response body closing
- ✅ Return descriptive errors
- ✅ Check `ApplicationvndApiJSON200` (or appropriate response field)

## Step 5: Add CLI Commands

Create `internal/commands/<domain>.go`:

```go
package commands

import (
    "context"
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "github.com/sam1el/snyk-api/internal/output"
    "github.com/sam1el/snyk-api/pkg/apiclients/targets"
)

// NewTargetsCmd creates the targets subcommand.
func NewTargetsCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:     "targets",
        Aliases: []string{"target"},
        Short:   "Manage Snyk targets.",
        Long:    `Manage Snyk targets (repositories and containers).`,
    }

    cmd.AddCommand(newTargetsListCmd())
    cmd.AddCommand(newTargetsGetCmd())

    return cmd
}

func newTargetsListCmd() *cobra.Command {
    var orgID string
    var limit int

    cmd := &cobra.Command{
        Use:   "list",
        Short: "List targets in an organization",
        Long: `List all targets within a specified Snyk organization.

Examples:
  # List targets
  snyk-api targets list --org-id=<org-id>

  # List with custom limit
  snyk-api targets list --org-id=<org-id> --limit 50
`,
        RunE: func(cmd *cobra.Command, args []string) error {
            if orgID == "" {
                return fmt.Errorf("organization ID is required. Use --org-id flag")
            }

            ctx := cmd.Context()
            baseClient, err := createClient(ctx)
            if err != nil {
                return fmt.Errorf("failed to create client: %w", err)
            }
            defer func() { _ = baseClient.Close() }()

            targetsClient := targets.NewTargetsClient(baseClient)

            params := &targets.ListTargetsParams{
                Limit: &limit,
            }

            targetList, err := targetsClient.ListTargets(ctx, orgID, params)
            if err != nil {
                return fmt.Errorf("failed to list targets: %w", err)
            }

            return output.Output(os.Stdout, targetList.Data, output.FormatType(outputFormat))
        },
    }

    cmd.Flags().StringVar(&orgID, "org-id", "", "Snyk organization ID")
    cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results")
    _ = cmd.MarkFlagRequired("org-id")

    return cmd
}

func newTargetsGetCmd() *cobra.Command {
    var orgID string

    cmd := &cobra.Command{
        Use:   "get <target-id>",
        Short: "Get target by ID",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            if orgID == "" {
                return fmt.Errorf("organization ID is required. Use --org-id flag")
            }

            ctx := cmd.Context()
            baseClient, err := createClient(ctx)
            if err != nil {
                return fmt.Errorf("failed to create client: %w", err)
            }
            defer func() { _ = baseClient.Close() }()

            targetsClient := targets.NewTargetsClient(baseClient)
            targetID := args[0]

            target, err := targetsClient.GetTarget(ctx, orgID, targetID)
            if err != nil {
                return fmt.Errorf("failed to get target: %w", err)
            }

            return output.Output(os.Stdout, target, output.FormatType(outputFormat))
        },
    }

    cmd.Flags().StringVar(&orgID, "org-id", "", "Snyk organization ID")
    _ = cmd.MarkFlagRequired("org-id")

    return cmd
}
```

Register the command in `internal/commands/root.go`:

```go
func init() {
    // ... existing commands ...
    rootCmd.AddCommand(NewTargetsCmd())
}
```

## Step 6: Test

### Unit Tests

Create `pkg/apiclients/<domain>/client_test.go`:

```go
package targets

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewTargetsClient(t *testing.T) {
    // Test client creation
    // (Similar to orgs/client_test.go)
}
```

### Manual Testing

```bash
# Build
make build

# Test CLI commands
./bin/snyk-api targets list --org-id=<org-id>
./bin/snyk-api targets get <target-id> --org-id=<org-id>

# Test with different output formats
./bin/snyk-api targets list --org-id=<org-id> --output table
./bin/snyk-api targets list --org-id=<org-id> --output yaml
```

### Run Linter

```bash
make lint
```

## Step 7: Document

Update documentation:

1. Add to `README.md` project status
2. Update `PHASE*_COMPLETE.md` or create new phase doc
3. Add example usage to `examples/`

## Checklist

Before submitting a PR:

- [ ] OpenAPI spec created in `.github/api-ref/rest/`
- [ ] Code generation configured (`oapi-codegen.yaml` + `doc.go`)
- [ ] Client generated successfully (`make generate`)
- [ ] Wrapper client implemented
- [ ] CLI commands added and registered
- [ ] Tests written (at least basic smoke tests)
- [ ] Linter passes (`make lint`)
- [ ] Manual testing completed
- [ ] Documentation updated
- [ ] Examples added (if applicable)

## Common Issues

### Issue: Generated client has compilation errors

**Cause**: OpenAPI spec has invalid schemas or missing references.

**Fix**: Validate the spec:
```bash
# Use online validator
https://editor.swagger.io/

# Or install validator
npm install -g @apidevtools/swagger-cli
swagger-cli validate .github/api-ref/rest/<domain>-minimal.yaml
```

### Issue: `undefined: resp.JSON200`

**Cause**: Generated client uses different response field names based on content type.

**Fix**: Check generated code for actual field name (usually `ApplicationvndApiJSON200`):
```go
// Not: resp.JSON200
// Use: resp.ApplicationvndApiJSON200
```

### Issue: Linter complains about unused code

**Cause**: Generated client includes code you're not using.

**Fix**: Add `//nolint` comments or exclude generated files:
```yaml
# .golangci.yaml
issues:
  exclude-rules:
    - path: _generated\.go$
      linters:
        - all
```

### Issue: UUID conversion errors

**Cause**: Need to convert strings to `openapi_types.UUID`.

**Fix**: Use `openapi_types.ParseUUID()`:
```go
uuid, err := openapi_types.ParseUUID(stringID)
if err != nil {
    return nil, fmt.Errorf("invalid ID: %w", err)
}
```

## Reference Implementations

See existing implementations for patterns:

- **Organizations API**: `pkg/apiclients/orgs/` - Simple list/get operations
- **Projects API**: `pkg/apiclients/projects/` - CRUD operations with filters

## Questions?

- 📖 [Architecture Decisions](planning/05-architecture-decisions.md)
- 💬 [Open a Discussion](https://github.com/sam1el/snyk-api/discussions)
- 🐛 [Report an Issue](https://github.com/sam1el/snyk-api/issues)

