#!/usr/bin/env bash
set -euo pipefail

# Script to pull OpenAPI specifications for Snyk APIs
# This will be expanded as we identify where Snyk publishes their API specs

API_REF_DIR=".github/api-ref"
TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

echo "========================================="
echo "Pulling Snyk API Specifications"
echo "Timestamp: $TIMESTAMP"
echo "========================================="

# Create directories
mkdir -p "$API_REF_DIR/rest"
mkdir -p "$API_REF_DIR/v1"

# NOTE: Snyk API specs must be downloaded manually from https://apidocs.snyk.io/
# The API endpoint https://api.snyk.io/rest/openapi returns a list of versions,
# not the actual OpenAPI specification.
#
# Current specs in repo:
#   - rest/rest-spec-2025-11-05.json (70,697 lines) - Full REST API
#   - v1/spec.yaml (23,929 lines) - Legacy v1 API
#
# To update:
#   1. Visit https://apidocs.snyk.io/
#   2. Download the latest OpenAPI spec
#   3. Save to the appropriate directory with date suffix
#   4. Run 'make generate' to regenerate clients

echo ""
echo "✅ Current Specifications:"
echo ""
echo "REST API:"
ls -lh "$API_REF_DIR/rest/"*.json 2>/dev/null | awk '{print "  " $9 " (" $5 ")"}'
echo ""
echo "v1 API:"
ls -lh "$API_REF_DIR/v1/"*.yaml 2>/dev/null | awk '{print "  " $9 " (" $5 ")"}'
echo ""

# Create a placeholder file
cat > "$API_REF_DIR/README.md" << 'EOF'
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
EOF

echo "✅ API spec directories created"
echo "📝 Next steps:"
echo "   1. Identify Snyk API spec sources"
echo "   2. Add download logic to this script"
echo "   3. Run 'make generate' to create clients"
echo ""
echo "Done!"
