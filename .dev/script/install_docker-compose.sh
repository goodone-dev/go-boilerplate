#!/bin/bash

# Script to install docker-compose if it is not installed

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Error: 'docker-compose' is not installed."
    echo ""
    echo "🤔 Would you like to install 'docker-compose'? (y/n)"
    read -r response

    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo "🔧 Installing 'docker-compose'..."
        brew install docker-compose

        if [ $? -eq 0 ]; then
            echo "✅ 'docker-compose' installed successfully!"
        else
            echo "❌ Failed to install 'docker-compose'. Please try installing manually:"
            echo "  brew install docker-compose"
            exit 1
        fi
    else
        echo "⏸️ Installation cancelled. To install 'docker-compose' later, run:"
        echo "  brew install docker-compose"
        exit 1
    fi
fi
