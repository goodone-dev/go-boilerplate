#!/bin/bash

# Script to install air if it is not installed

if ! command -v air &> /dev/null; then
    echo "❌ Error: 'air' is not installed."
    echo ""
    echo "🤔 Would you like to install 'air'? (y/n)"
    read -r response

    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo "🔧 Installing 'air'..."
        go install github.com/air-verse/air@latest

        if [ $? -eq 0 ]; then
            echo "✅ 'air' installed successfully!"
        else
            echo "❌ Failed to install 'air'. Please try installing manually:"
            echo "  go install github.com/air-verse/air@latest"
            exit 1
        fi
    else
        echo "⏸️ Installation cancelled. To install 'air' later, run:"
        echo "  go install github.com/air-verse/air@latest"
        exit 1
    fi
fi
