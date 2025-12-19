# snyk-api

> **Production-Ready**: Comprehensive Snyk API management tool with CLI interface.

[![CI](https://github.com/sam1el/snyk-api/actions/workflows/ci.yaml/badge.svg)](https://github.com/sam1el/snyk-api/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sam1el/snyk-api.svg)](https://pkg.go.dev/github.com/sam1el/snyk-api)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sam1el/snyk-api)](go.mod)

A comprehensive Go library and CLI tool for managing Snyk APIs. Built on [go-application-framework](https://github.com/snyk/go-application-framework) with type-safe clients generated from OpenAPI specifications.

> **Note**: Currently developed at `github.com/sam1el/snyk-api`. Upon maturity and official adoption, it will be migrated to `github.com/snyk/snyk-api`.

## ✨ Features

- 🎯 **Type-Safe API Clients** - Generated from official OpenAPI specs
- 🔄 **Rate Limiting** - Token bucket algorithm with worker pools
- 🔁 **Smart Retries** - Exponential backoff with jitter
- 🎨 **Multiple Output Formats** - JSON, YAML, and table views
- 🔌 **Framework Integration** - Built on go-application-framework
- 🌍 **Multi-Region Support** - All Snyk regional endpoints
- 📦 **Library + CLI** - Use as Go package or standalone tool

## 📦 Installation

### CLI Installation

```bash
go install github.com/sam1el/snyk-api/cmd/snyk-api@latest
```

### Library Installation

```bash
go get github.com/sam1el/snyk-api
```

## 🚀 Quick Start

### CLI Usage

```bash
# Set your Snyk API token
export SNYK_TOKEN=your-snyk-token-here

# List organizations
snyk-api orgs list

# List projects in an organization
snyk-api projects list --org-id=<org-id>

# Filter projects by type
snyk-api projects list --org-id=<org-id> --type npm --origin github

# Get organization details (YAML output)
snyk-api orgs get <org-id> --output yaml

# Get project details (table output)
snyk-api projects get <project-id> --org-id=<org-id> --output table

# Delete a project
snyk-api projects delete <project-id> --org-id=<org-id>
```

### Library Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sam1el/snyk-api/pkg/apiclients/orgs"
    "github.com/sam1el/snyk-api/pkg/apiclients/projects"
    "github.com/sam1el/snyk-api/pkg/client"
)

func main() {
    ctx := context.Background()

    // Create base client (uses SNYK_TOKEN from environment)
    baseClient, err := client.New(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer baseClient.Close()

    // List organizations
    orgsClient := orgs.NewOrgsClient(baseClient)
    limit := 10
    orgsList, err := orgsClient.ListOrganizations(ctx, &orgs.ListOrganizationsParams{
        Limit: &limit,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, org := range orgsList.Data {
        fmt.Printf("Org: %s (%s)\n", org.Attributes.Name, org.Id)
    }

    // List projects in first organization
    if len(orgsList.Data) > 0 {
        orgID := orgsList.Data[0].Id.String()
        projectsClient := projects.NewProjectsClient(baseClient)

        projectsList, err := projectsClient.ListProjects(ctx, orgID, &projects.ListProjectsParams{
            Limit: &limit,
        })
        if err != nil {
            log.Fatal(err)
        }

        for _, project := range projectsList.Data {
            fmt.Printf("  Project: %s (%s)\n", project.Attributes.Name, project.Attributes.Type)
        }
    }
}
```

### Advanced Client Configuration

```go
// Custom configuration
baseClient, err := client.New(ctx,
    client.WithVersion(client.APIVersion("2024-04-22")),           // Specific API version
    client.WithBaseURL("https://api.snyk.io", "https://api.snyk.io/rest"), // Custom endpoint
    client.WithRateLimitConfig(ratelimit.Config{
        BurstSize:      20,
        Period:         time.Second,
        MaxRetries:     10,
        RetryBaseDelay: 200 * time.Millisecond,
        RetryMaxDelay:  10 * time.Second,
    }),
)
```

## 🎯 Available Commands

```bash
snyk-api --help                 # Show all commands
snyk-api version               # Version information

# Organizations
snyk-api orgs list             # List organizations
snyk-api orgs get <id>         # Get organization by ID

# Projects
snyk-api projects list --org-id=<id>              # List projects
snyk-api projects list --org-id=<id> --origin github --type npm  # Filter projects
snyk-api projects get <id> --org-id=<id>         # Get project by ID
snyk-api projects delete <id> --org-id=<id>      # Delete project
```

### Global Flags

```bash
--output, -o       Output format: json, yaml, table (default: json)
--api-url          Override Snyk API URL
--api-version      Snyk API version (default: 2025-11-05)
--debug            Enable debug logging
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SNYK_TOKEN` | API authentication token | **Required** |
| `SNYK_API` | Custom API endpoint | `https://api.snyk.io` |

### API Versioning

Snyk's REST API uses date-based versioning. Specify versions per request:

```bash
# Use specific API version
snyk-api orgs list --api-version 2024-04-22

# Use beta version
snyk-api orgs list --api-version 2025-11-05~beta

# Use experimental version
snyk-api orgs list --api-version 2025-11-05~experimental
```

In Go code:

```go
// Use default latest GA version
client.New(ctx)

// Use specific version
client.New(ctx, client.WithVersion(client.Version20240422))

// Use beta
client.New(ctx, client.WithVersion(client.DefaultAPIVersion.Beta()))
```

See [API_VERSIONING.md](docs/planning/API_VERSIONING.md) for details.

## 📊 Project Status

**Current**: Production-ready with 2 API domains ✅

| Component | Status | Coverage |
|-----------|--------|----------|
| **Organizations API** | ✅ Complete | list, get |
| **Projects API** | ✅ Complete | list, get, delete |
| **Targets API** | ⏭️ Future | - |
| **Issues API** | ⏭️ Future | - |
| **Ignores API** | ⏭️ Future | - |

**Quality Metrics:**
- ✅ 21 Go source files (~3,500 lines)
- ✅ 13 passing tests with race detection
- ✅ 0 linter issues (golangci-lint)
- ✅ 0 security vulnerabilities (Snyk Code + SCA)

## 🏗️ Architecture

Built on proven patterns from Snyk's go-application-framework:

```
OpenAPI Spec → oapi-codegen → Generated Client → Wrapper → CLI
                                      ↓
                                Base Client (rate limit, retry, auth)
                                      ↓
                                Snyk REST API
```

### Project Structure

```
snyk-api/
├── cmd/
│   └── snyk-api/              # CLI entry point
├── pkg/
│   ├── client/                # Base HTTP client with rate limiting
│   ├── apiclients/
│   │   ├── orgs/             # Organizations API client
│   │   └── projects/         # Projects API client
│   └── config/               # Configuration management
├── internal/
│   ├── commands/             # Cobra CLI commands
│   ├── output/               # Output formatters
│   └── ratelimit/            # Rate limiting implementation
├── .github/
│   ├── api-ref/              # OpenAPI specifications
│   └── workflows/            # CI/CD pipelines
└── docs/
    └── planning/             # Architecture docs
```

## 🛠️ Development

### Prerequisites

- Go 1.24 or later
- Make
- golangci-lint
- oapi-codegen

### Setup

```bash
# Clone repository
git clone https://github.com/sam1el/snyk-api.git
cd snyk-api

# Install development tools
make install-tools

# Run tests
make test

# Run linters
make lint

# Build CLI
make build

# Run all checks
make all
```

### Make Targets

```bash
make help              # Show all available targets
make build             # Build the CLI binary
make test              # Run tests with coverage
make lint              # Run golangci-lint
make fmt               # Format code (gofmt + gofmt -s)
make vet               # Run go vet
make generate          # Generate API clients from OpenAPI specs
make pull-specs        # Pull latest OpenAPI specifications
make clean             # Clean build artifacts
make dev               # Run fmt + vet + lint + test
```

### Adding New API Domains

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed instructions on adding new API domains.

Quick overview:
1. Create minimal OpenAPI spec in `.github/api-ref/rest/`
2. Add `oapi-codegen.yaml` config
3. Generate client: `make generate`
4. Create wrapper client in `pkg/apiclients/<domain>/`
5. Add CLI commands in `internal/commands/`
6. Add tests

## 📚 Documentation

- [Planning Documents](docs/planning/) - Architecture decisions and analysis
- [API Versioning Strategy](docs/planning/API_VERSIONING.md)
- [Framework Analysis](docs/planning/02-framework-analysis.md)
- [Phase Completion Reports](docs/planning/)
- [Contributing Guide](CONTRIBUTING.md)

## 🤝 Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linters (`make dev`)
5. Commit your changes (`git commit -m 'feat: add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 🔗 Related Projects

- [go-application-framework](https://github.com/snyk/go-application-framework) - Core Snyk framework
- [cli](https://github.com/snyk/cli) - Official Snyk CLI
- [snyk-request-manager-go](https://github.com/snyk/snyk-request-manager-go) - Rate-limited request manager

## 📄 License

Apache-2.0 - see [LICENSE](LICENSE) for details.

## 💬 Support

- 🐛 [Report a bug](https://github.com/sam1el/snyk-api/issues/new?labels=bug)
- 💡 [Request a feature](https://github.com/sam1el/snyk-api/issues/new?labels=enhancement)
- 📖 [Documentation](https://pkg.go.dev/github.com/sam1el/snyk-api)
- 💬 [Discussions](https://github.com/sam1el/snyk-api/discussions)

---

**Built with ❤️ for the Snyk community**
