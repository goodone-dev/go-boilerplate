#!/bin/bash

# Delete docker containers
echo "🗑️ Deleting docker containers..."
docker-compose down

# Check if docker containers deleted successfully
if [ $? -eq 0 ]; then
    echo "✅ Docker containers deleted successfully!"
else
    echo "❌ Error: Failed to delete docker containers"
    exit 1
fi
