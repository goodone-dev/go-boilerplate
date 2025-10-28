#!/bin/bash

if ! command -v air &> /dev/null; then
    echo "Error: 'air' is not installed."
    echo ""
    echo "Would you like to install 'air' for live reloading? (y/n)"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo ""
        echo "Installing 'air'..."
        go install github.com/air-verse/air@latest
        
        if [ $? -eq 0 ]; then
            echo ""
            echo "✓ 'air' installed successfully!"
            echo ""
            echo "Starting application with live reloading..."
            air -c .air.toml
        else
            echo ""
            echo "✗ Failed to install 'air'. Please try installing manually:"
            echo "  go install github.com/air-verse/air@latest"
            exit 1
        fi
    else
        echo ""
        echo "Installation cancelled. To install 'air' later, run:"
        echo "  go install github.com/air-verse/air@latest"
        exit 1
    fi
fi