# Examples

This directory contains practical examples demonstrating how to use the `snyk-api` library.

## Prerequisites

All examples require:
1. Go 1.24 or later
2. A Snyk API token set in the `SNYK_TOKEN` environment variable

```bash
export SNYK_TOKEN=your-token-here
```

## Available Examples

### Basic Usage ([basic/](basic/))

Demonstrates fundamental operations:
- Creating a client
- Listing organizations
- Listing projects in an organization

**Run:**
```bash
go run examples/basic/main.go
```

### Advanced Usage ([advanced/](advanced/))

Shows advanced features:
- Custom client configuration
- Custom rate limiting
- Pagination through results
- Grouping and analyzing project data
- Using specific API versions

**Run:**
```bash
# Edit the file to set your organization ID first
go run examples/advanced/main.go
```

### Filtering ([filtering/](filtering/))

Demonstrates filtering and searching:
- Filter projects by type (npm, maven, etc.)
- Filter projects by origin (github, cli, etc.)
- Combined filters
- Client-side search patterns

**Run:**
```bash
# Edit the file to set your organization ID first
go run examples/filtering/main.go
```

## Getting Organization IDs

If you don't know your organization ID, you can list all organizations:

```bash
# Using the CLI
snyk-api orgs list

# Or run the basic example
go run examples/basic/main.go
```

## Common Patterns

### Error Handling

```go
result, err := orgsClient.ListOrganizations(ctx, params)
if err != nil {
    log.Fatalf("Failed to list organizations: %v", err)
}
```

### Resource Cleanup

Always close the client when done:

```go
defer func() {
    if err := baseClient.Close(); err != nil {
        log.Printf("Failed to close client: %v", err)
    }
}()
```

### Pagination

```go
var cursor *string
for {
    params := &projects.ListProjectsParams{
        Limit:         &limit,
        StartingAfter: cursor,
    }
    
    result, err := projectsClient.ListProjects(ctx, orgID, params)
    if err != nil {
        return err
    }
    
    // Process result.Data...
    
    if result.Links == nil || result.Links.Next == nil {
        break // No more pages
    }
    
    // Extract cursor from next link
    cursor = extractCursor(result.Links.Next)
}
```

### Filtering

```go
// Filter by type
npmType := "npm"
params := &projects.ListProjectsParams{
    Type: &npmType,
}

// Filter by origin
githubOrigin := "github"
params := &projects.ListProjectsParams{
    Origin: &githubOrigin,
}

// Combined filters
params := &projects.ListProjectsParams{
    Type:   &npmType,
    Origin: &githubOrigin,
}
```

## Contributing Examples

Have a useful example? Please contribute!

1. Create a new directory under `examples/`
2. Add a well-commented `main.go`
3. Update this README
4. Submit a pull request

## Troubleshooting

### "Failed to create client: SNYK_TOKEN not set"

Set your Snyk API token:
```bash
export SNYK_TOKEN=your-token-here
```

### "Failed to list organizations: 401 Unauthorized"

Your token is invalid or expired. Generate a new one at https://app.snyk.io/account

### "Failed to list projects: 404 Not Found"

The organization ID doesn't exist or you don't have access. Verify the ID:
```bash
snyk-api orgs list
```

## Additional Resources

- [Main README](../README.md)
- [API Documentation](https://pkg.go.dev/github.com/sam1el/snyk-api)
- [Contributing Guide](../CONTRIBUTING.md)

