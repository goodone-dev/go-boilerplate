#!/bin/bash

# Run tests with coverage
echo "🧪 Running tests with coverage..."
go test -v -race -coverprofile=coverage.out -covermode=atomic -coverpkg=./... ./internal/application/...
go tool cover -html=coverage.out -o=coverage.html

# Check coverage thresholds using go-test-coverage
echo ""
echo "📊 Checking coverage thresholds..."
go-test-coverage --config=./.testcoverage.yml
