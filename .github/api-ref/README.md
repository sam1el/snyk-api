# API Reference Specifications

This directory contains OpenAPI specifications for Snyk APIs.

## Structure

- `rest/` - REST API (modern) specifications
- `v1/` - v1 API (legacy) specifications

## Updating Specs

Run `./scripts/pull-api-specs.sh` to fetch the latest specifications.

## Sources

- REST API: https://apidocs.snyk.io/
- v1 API: To be documented

## Generated Clients

API clients are generated from these specifications using `oapi-codegen`.
Run `make generate` after updating specs.
