#!/bin/bash

# Stop docker containers
echo "🛑 Stopping docker containers..."
docker-compose stop

# Check if docker containers stopped successfully
if [ $? -eq 0 ]; then
    echo "✅ Docker containers stopped successfully!"
else
    echo "❌ Error: Failed to stop docker containers"
    exit 1
fi
