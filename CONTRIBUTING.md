# Contributing to snyk-api

Thank you for your interest in contributing to snyk-api!

## Development Setup

1. **Prerequisites**
   - Go 1.24 or later
   - Make
   - golangci-lint

2. **Clone and Setup**

   ```bash
   git clone https://github.com/sam1el/snyk-api.git
   cd snyk-api
   make install-tools
   ```

3. **Run Tests**

   ```bash
   make test
   make lint
   ```

## Code Style

- Follow Go best practices and idiomatic code
- Use `gofmt` and `golangci-lint`
- Write tests for new functionality
- Keep functions focused and small

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make dev` to verify
6. Submit a pull request

## Commit Messages

Use conventional commits format:

- `feat: add new feature`
- `fix: resolve bug`
- `docs: update documentation`
- `test: add tests`
- `refactor: code improvements`

## Questions?

Open an issue or discussion on GitHub.
