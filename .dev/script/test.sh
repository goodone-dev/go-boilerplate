#!/bin/bash

echo "🧪 Running tests"

# Run tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic -coverpkg=./... ./internal/application/...
go tool cover -func=coverage.out
