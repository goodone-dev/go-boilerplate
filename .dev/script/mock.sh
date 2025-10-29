#!/bin/bash

# Script to generate mock files using mockery

# Check if 'mockery' is installed, and install it if not
$(dirname "$0")/ensure_mockery.sh

echo "🤖 Generating mock files..."

# Generate mock files
mockery --log-level=ERROR

# Check if mock files generated successfully
if [ $? -eq 0 ]; then
    echo "✅ Mock files generated successfully!"
else
    echo "❌ Error: Failed to generate mock files"
    exit 1
fi
