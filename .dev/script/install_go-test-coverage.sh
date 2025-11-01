#!/bin/bash

# Script to install go-test-coverage if it is not installed

if ! command -v go-test-coverage &> /dev/null; then
    echo "❌ Error: 'go-test-coverage' is not installed."
    echo ""
    echo "🤔 Would you like to install 'go-test-coverage'? (y/n)"
    read -r response

    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo "🔧 Installing 'go-test-coverage'..."
        go install github.com/vladopajic/go-test-coverage/v2@latest

        if [ $? -eq 0 ]; then
            echo "✅ 'go-test-coverage' installed successfully!"
        else
            echo "❌ Failed to install 'go-test-coverage'. Please try installing manually:"
            echo "  go install github.com/vladopajic/go-test-coverage/v2@latest"
            exit 1
        fi
    else
        echo "⏸️ Installation cancelled. To install 'go-test-coverage' later, run:"
        echo "  go install github.com/vladopajic/go-test-coverage/v2@latest"
        exit 1
    fi
fi
