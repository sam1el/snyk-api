# Contributing to snyk-api

Thank you for your interest in contributing to `snyk-api`! We welcome contributions from the community.

## 📋 Table of Contents

- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Adding New API Domains](#adding-new-api-domains)
- [Code Style](#code-style)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Commit Messages](#commit-messages)
- [Questions](#questions)

## 🛠️ Development Setup

### Prerequisites

- **Go 1.24 or later** - [Install Go](https://go.dev/doc/install)
- **Make** - Usually pre-installed on Unix systems
- **golangci-lint** - [Installation guide](https://golangci-lint.run/usage/install/)
- **oapi-codegen** - Installed via `go get` during setup

### Clone and Setup

```bash
# Fork the repository on GitHub first, then:
git clone https://github.com/YOUR_USERNAME/snyk-api.git
cd snyk-api

# Install development tools
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Verify setup
make help
```

### Environment Configuration

```bash
# Required for testing
export SNYK_TOKEN=your-snyk-token-here

# Optional: Custom API endpoint
export SNYK_API=https://api.snyk.io
```

## 🔄 Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

Follow the [Code Style](#code-style) guidelines.

### 3. Run Development Checks

```bash
# Format code
make fmt

# Run static analysis
make vet

# Run linters
make lint

# Run tests
make test

# Or run all checks at once
make dev
```

### 4. Test Your Changes

```bash
# Build the CLI
make build

# Test manually
./bin/snyk-api --help
./bin/snyk-api orgs list

# Run examples
go run examples/basic/main.go
```

### 5. Commit and Push

```bash
git add .
git commit -m "feat: add amazing feature"
git push origin feature/your-feature-name
```

### 6. Create Pull Request

Open a PR on GitHub with:
- Clear title and description
- Reference to any related issues
- Screenshots/examples if applicable

## ➕ Adding New API Domains

See the comprehensive guide: **[docs/ADDING_API_DOMAINS.md](docs/ADDING_API_DOMAINS.md)**

Quick overview:
1. Create OpenAPI spec in `.github/api-ref/rest/`
2. Configure `oapi-codegen`
3. Generate client with `make generate`
4. Write wrapper client
5. Add CLI commands
6. Write tests
7. Update documentation

## 🎨 Code Style

### Go Guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting (automatic with `make fmt`)
- Keep functions small and focused
- Write clear, descriptive names
- Add package-level comments

### Example: Good Function

```go
// ListOrganizations retrieves all organizations accessible to the authenticated user.
// It supports pagination via the params.StartingAfter cursor.
func (c *OrgsClient) ListOrganizations(ctx context.Context, params *ListOrganizationsParams) (*OrganizationList, error) {
    // Implementation...
}
```

### Project Patterns

**Error Handling:**
```go
// Always wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to list organizations: %w", err)
}
```

**Resource Cleanup:**
```go
// Always defer cleanup
defer func() {
    if err := resp.Body.Close(); err != nil {
        log.Printf("failed to close body: %v", err)
    }
}()
```

**UUID Conversion:**
```go
// Use openapi_types.ParseUUID for UUIDs
uuid, err := openapi_types.ParseUUID(stringID)
if err != nil {
    return nil, fmt.Errorf("invalid ID format: %w", err)
}
```

## 🧪 Testing

### Unit Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./pkg/client/
```

### Writing Tests

```go
func TestNewOrgsClient(t *testing.T) {
    // Use testify for assertions
    assert := assert.New(t)
    
    // Create test client
    client := NewOrgsClient(baseClient)
    
    assert.NotNil(client)
    assert.NotNil(client.apiClient)
}
```

### Integration Tests

```bash
# Set up test environment
export SNYK_TOKEN=your-test-token

# Run examples as integration tests
go run examples/basic/main.go
go run examples/filtering/main.go
```

### Linting

```bash
# Run all linters
make lint

# Run specific linters
golangci-lint run --enable-only=errcheck
golangci-lint run --enable-only=staticcheck
```

## 🔀 Pull Request Process

### Before Submitting

- [ ] All tests pass (`make test`)
- [ ] Linters pass (`make lint`)
- [ ] Code is formatted (`make fmt`)
- [ ] Documentation is updated
- [ ] Examples are added (if applicable)
- [ ] Commit messages follow conventions

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No new warnings
```

### Review Process

1. Automated checks run (CI/CD)
2. Maintainer reviews code
3. Address feedback
4. Approval and merge

## 📝 Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/).

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Test additions or changes
- `chore:` - Maintenance tasks

### Examples

**Good:**
```bash
feat: add targets API support

- Created targets-minimal.yaml OpenAPI spec
- Generated client with oapi-codegen
- Added CLI commands for list and get operations
- Updated documentation

Closes #42
```

**Simple:**
```bash
fix: correct UUID parsing in projects client
```

**Breaking Change:**
```bash
feat!: change API client interface

BREAKING CHANGE: Client.Do() renamed to Client.Execute()
```

## 📚 Documentation

### What to Document

- New API domains
- CLI commands
- Configuration options
- Breaking changes
- Migration guides

### Where to Document

- **README.md** - High-level overview, quick start
- **docs/ADDING_API_DOMAINS.md** - Developer guide
- **examples/** - Practical usage examples
- **Code comments** - Implementation details
- **CHANGELOG.md** - Release notes

### Documentation Style

- Use clear, concise language
- Include code examples
- Add CLI usage examples
- Explain "why" not just "what"

## 🤝 Code Review

### As a Reviewer

- Be respectful and constructive
- Explain reasoning for suggestions
- Approve when ready
- Test locally if needed

### As a Contributor

- Respond promptly to feedback
- Ask questions if unclear
- Make requested changes
- Update tests/docs as needed

## 🐛 Reporting Bugs

Use the [Bug Report template](https://github.com/sam1el/snyk-api/issues/new?labels=bug).

Include:
- Go version (`go version`)
- CLI version (`snyk-api version`)
- Steps to reproduce
- Expected vs actual behavior
- Error messages

## 💡 Requesting Features

Use the [Feature Request template](https://github.com/sam1el/snyk-api/issues/new?labels=enhancement).

Include:
- Use case description
- Expected behavior
- Example usage
- Alternative solutions considered

## ❓ Questions?

- 💬 [Open a Discussion](https://github.com/sam1el/snyk-api/discussions)
- 🐛 [Report an Issue](https://github.com/sam1el/snyk-api/issues)
- 📖 [Read the Docs](docs/)
- 📧 Contact maintainers

## 📜 License

By contributing, you agree that your contributions will be licensed under the Apache-2.0 License.

---

**Thank you for contributing to snyk-api!** 🎉
