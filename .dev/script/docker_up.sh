#!/bin/bash

# Script to start docker containers

# Check if 'docker-compose' is installed, and install it if not
$(dirname "$0")/install_docker-compose.sh

echo "🚀 Starting docker containers..."

# Start docker containers
docker-compose up --build -d

# Check if docker containers started successfully
if [ $? -eq 0 ]; then
    echo "✅ Docker containers started successfully!"
else
    echo "❌ Error: Failed to start docker containers"
    exit 1
fi
