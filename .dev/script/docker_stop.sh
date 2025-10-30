#!/bin/bash

# Script to stop docker containers

# Check if 'docker-compose' is installed, and install it if not
$(dirname "$0")/install_docker-compose.sh

echo "🛑 Stopping docker containers..."

# Stop docker containers
docker-compose stop

# Check if docker containers stopped successfully
if [ $? -eq 0 ]; then
    echo "✅ Docker containers stopped successfully!"
else
    echo "❌ Error: Failed to stop docker containers"
    exit 1
fi
