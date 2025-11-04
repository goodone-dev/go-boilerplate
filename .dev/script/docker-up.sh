#!/bin/bash

# Start docker containers
echo "🚀 Starting docker containers..."
docker-compose up --build -d

# Check if docker containers started successfully
if [ $? -eq 0 ]; then
    echo "✅ Docker containers started successfully!"
else
    echo "❌ Error: Failed to start docker containers"
    exit 1
fi
