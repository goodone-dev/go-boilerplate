#!/bin/bash

# Script to run the application

# Check if --watch argument is provided
WATCH_MODE=false
for arg in "$@"; do
    if [ "$arg" = "-w" ]; then
        WATCH_MODE=true
        break
    fi
done

if [ "$WATCH_MODE" = true ]; then
    # Check if 'air' is installed, and install it if not
    $(dirname "$0")/ensure_air.sh

    echo "🚀 Starting application with live reloading..."
    air -c .air.toml
else
    echo "🚀 Starting application..."
    go run ./cmd/api/main.go
fi
