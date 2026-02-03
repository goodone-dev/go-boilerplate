#!/bin/bash

# Check if --watch argument is provided
WATCH_MODE=false
for arg in "$@"; do
    if [ "$arg" = "-w" ]; then
        WATCH_MODE=true
        break
    fi
done

if [ "$WATCH_MODE" = true ]; then
    air -c .air.toml
else
    go run ./cmd/rest/main.go
fi
