#!/bin/bash

set -e

# Check if 'go-test-coverage' is installed, and install it if not
$(dirname "$0")/install_go-test-coverage.sh

# Run tests with coverage
echo "🧪 Running tests with coverage..."
go test -v -race -coverprofile=coverage.out -covermode=atomic -coverpkg=./... ./internal/application/...

# Check coverage thresholds using go-test-coverage
echo ""
echo "📊 Checking coverage thresholds..."
go-test-coverage --config=./.testcoverage.yml
