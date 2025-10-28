#!/bin/bash

if ! command -v mockery &> /dev/null; then
    echo "Error: 'mockery' is not installed."
    echo ""
    echo "Would you like to install 'mockery'? (y/n)"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo ""
        echo "Installing 'mockery'..."
        go install github.com/vektra/mockery/v2@latest
        
        if [ $? -eq 0 ]; then
            echo ""
            echo "✓ 'mockery' installed successfully!"
            echo ""
        else
            echo ""
            echo "✗ Failed to install 'mockery'. Please try installing manually:"
            echo "  go install github.com/vektra/mockery/v2@latest"
            echo "  Or visit: https://vektra.github.io/mockery/latest/installation/"
            exit 1
        fi
    else
        echo ""
        echo "Installation cancelled. To install 'mockery' later, run:"
        echo "  go install github.com/vektra/mockery/v2@latest"
        echo "  Or visit: https://vektra.github.io/mockery/latest/installation/"
        exit 1
    fi
fi