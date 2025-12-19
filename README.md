# snyk-api

> **Note**: This project is under active development at `github.com/sam1el/snyk-api`. Upon maturity and adoption, it will be migrated to the official Snyk organization.

[![CI](https://github.com/sam1el/snyk-api/actions/workflows/ci.yaml/badge.svg)](https://github.com/sam1el/snyk-api/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sam1el/snyk-api.svg)](https://pkg.go.dev/github.com/sam1el/snyk-api)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A comprehensive Go library and CLI tool for managing Snyk API endpoints. Provides complete coverage of both v1 and REST APIs with built-in rate limiting, retry logic, and type-safe clients.

> **Status:** 🚧 Under active development - Phase 0 (Project Setup) complete

## Features

- ✅ **Complete API Coverage**: Support for all Snyk v1 and REST API endpoints
- ✅ **Type-Safe Clients**: Generated from OpenAPI specifications
- ✅ **Rate Limiting**: Built-in request queue with configurable limits
- ✅ **Retry Logic**: Automatic retries with exponential backoff
- ✅ **Multi-Region**: Support for all Snyk regional endpoints
- ✅ **Authentication**: Token, OAuth, and PAT support
- ✅ **Library & CLI**: Use as Go library or standalone CLI tool

## Installation

### As a Library

```bash
go get github.com/sam1el/snyk-api
```

### As a CLI

```bash
go install github.com/sam1el/snyk-api/cmd/snyk-api@latest
```

## Quick Start

### Library Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sam1el/snyk-api/pkg/client"
)

func main() {
    // Create client (uses SNYK_TOKEN from environment)
    c, err := client.New()
    if err != nil {
        log.Fatal(err)
    }

    // List organizations
    ctx := context.Background()
    orgs, err := c.Orgs().List(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, org := range orgs.Data {
        fmt.Printf("Org: %s (%s)\n", org.Attributes.Name, org.ID)
    }
}
```

### CLI Usage

```bash
# Set your token
export SNYK_TOKEN=your-token-here

# List organizations
snyk-api orgs list

# Get a specific project
snyk-api projects get --org-id=xxx --project-id=yyy

# Create a target
snyk-api targets create --org-id=xxx --data=@payload.json
```

## Configuration

The tool respects standard Snyk environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SNYK_TOKEN` | API authentication token | Required |
| `SNYK_API` | Custom API endpoint | `https://api.snyk.io` |
| `SNYK_OAUTH_TOKEN` | OAuth token (alternative to SNYK_TOKEN) | - |

## Architecture

Built on top of Snyk's [go-application-framework](https://github.com/snyk/go-application-framework), ensuring consistency with the official Snyk CLI.

```
snyk-api/
├── cmd/snyk-api/       # CLI entry point
├── pkg/
│   ├── client/         # High-level client API
│   ├── apiclients/     # Generated API clients (by domain)
│   └── extension/      # CLI extension for Snyk CLI integration
├── internal/
│   └── ratelimit/      # Rate limiting implementation
└── .github/
    └── api-ref/        # OpenAPI specifications
```

## Development

### Prerequisites

- Go 1.24 or later
- Make

### Setup

```bash
# Clone the repository
git clone https://github.com/sam1el/snyk-api.git
cd snyk-api

# Install development tools
make install-tools

# Run tests
make test

# Run linters
make lint

# Build
make build
```

### Available Make Targets

```bash
make help           # Show all available targets
make build          # Build the CLI binary
make test           # Run tests
make lint           # Run linters
make fmt            # Format code
make generate       # Generate API clients from OpenAPI specs
make pull-specs     # Pull latest OpenAPI specifications
make clean          # Clean build artifacts
```

## Project Status

### Phase 0: Project Setup ✅ Complete
- [x] Go module initialization
- [x] Directory structure
- [x] Makefile and tooling
- [x] Linting configuration
- [x] CI/CD pipelines
- [x] Documentation

### Phase 1: Core Infrastructure 🚧 In Progress
- [ ] Framework integration
- [ ] Configuration layer
- [ ] Rate limiting middleware
- [ ] Base API client

### Phase 2: First API Domains + CLI
- [ ] Organizations API
- [ ] Projects API
- [ ] CLI framework
- [ ] Basic commands

### Phase 3: Expand API Coverage
- [ ] Targets API
- [ ] Issues API
- [ ] Ignores API
- [ ] Additional domains

## Documentation

- [Planning Documents](./docs/planning/) - Comprehensive analysis and design decisions
- [API Reference](https://pkg.go.dev/github.com/sam1el/snyk-api) - Go package documentation
- [Contributing Guide](./CONTRIBUTING.md) - How to contribute

## Related Projects

- [go-application-framework](https://github.com/snyk/go-application-framework) - Core Snyk framework
- [snyk-request-manager-go](https://github.com/sam1el/snyk-request-manager-go) - Rate-limited request manager
- [cli](https://github.com/snyk/cli) - Official Snyk CLI

## License

Apache-2.0 - see [LICENSE](LICENSE) for details.

## Support

For issues and questions:
- 🐛 [Report a bug](https://github.com/sam1el/snyk-api/issues/new?template=bug_report.md)
- 💡 [Request a feature](https://github.com/sam1el/snyk-api/issues/new?template=feature_request.md)
- 💬 [Discussions](https://github.com/sam1el/snyk-api/discussions)
