#!/bin/bash

# Script to install golang-migrate if it is not installed

if ! command -v migrate &> /dev/null; then
    echo "❌ Error: 'migrate' is not installed."
    echo ""
    echo "🤔 Would you like to install 'golang-migrate'? (y/n)"
    read -r response

    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo ""
        echo "🔧 Installing 'golang-migrate'..."
        go install -tags "$DB_DRIVER" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

        if [ $? -eq 0 ]; then
            echo ""
            echo "✅ 'migrate' installed successfully!"
            echo ""
        else
            echo ""
            echo "❌ Failed to install 'migrate'. Please try installing manually:"
            echo "  go install -tags '$DB_DRIVER' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
            exit 1
        fi
    else
        echo ""
        echo "⏸️ Installation cancelled. To install 'migrate' later, run:"
        echo "  go install -tags '$DB_DRIVER' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
fi
