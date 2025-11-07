#!/bin/bash

FORCED=false
VERBOSE=false

while getopts ":fv" opt; do
    case $opt in
        f) FORCED=true;;
        v) VERBOSE=true;;
    esac
done

if ! command -v air &> /dev/null; then
    if [ "$FORCED" = true ]; then
        install_air
    else
        echo "❌ Error: 'air' is not installed."
        echo ""
        echo "🤔 Would you like to install 'air'? (y/n)"
        read -r response

        if [[ "$response" =~ ^[Yy]$ ]]; then
            install_air
        else
            echo "⏸️ Installation cancelled. To install 'air' later, run:"
            echo "  go install github.com/air-verse/air@latest"
            exit 1
        fi
    fi
else
    if [ "$VERBOSE" = true ]; then
        echo "✅ 'air' is already installed."
    fi
fi

install_air() {
    echo "🔧 Installing 'air'..."
    go install github.com/air-verse/air@latest

    if [ $? -eq 0 ]; then
        echo "✅ 'air' installed successfully!"
    else
        echo "❌ Failed to install 'air'. Please try installing manually:"
        echo "  go install github.com/air-verse/air@latest"
        exit 1
    fi
}
