#!/bin/bash

# Script to delete docker containers

# Check if 'docker-compose' is installed, and install it if not
$(dirname "$0")/ensure_docker-compose.sh

echo "🗑️ Deleting docker containers..."

# Delete docker containers
docker-compose down

# Check if docker containers deleted successfully
if [ $? -eq 0 ]; then
    echo "✅ Docker containers deleted successfully!"
else
    echo "❌ Error: Failed to delete docker containers"
    exit 1
fi