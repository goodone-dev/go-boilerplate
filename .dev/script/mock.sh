#!/bin/bash

# Generate mock files
echo "🤖 Generating mock files..."
mockery --log-level=ERROR

# Check if mock files generated successfully
if [ $? -eq 0 ]; then
    echo "✅ Mock files generated successfully!"
else
    echo "❌ Error: Failed to generate mock files"
    exit 1
fi
