#!/bin/bash

FORCED=false
VERBOSE=false

while getopts ":fv" opt; do
    case $opt in
        f) FORCED=true;;
        v) VERBOSE=true;;
    esac
done

if ! command -v mockery &> /dev/null; then
    if [ "$FORCED" = true ]; then
        install_mockery
    else
        echo "❌ Error: 'mockery' is not installed."
        echo ""
        echo "🤔 Would you like to install 'mockery'? (y/n)"
        read -r response

        if [[ "$response" =~ ^[Yy]$ ]]; then
            install_mockery
        else
            echo "⏸️ Installation cancelled. To install 'mockery' later, run:"
            echo "  go install github.com/vektra/mockery/v3@latest"
            exit 1
        fi
    fi
else
    if [ "$VERBOSE" = true ]; then
        echo "✅ 'mockery' is already installed."
    fi
fi

install_mockery() {
    echo "🔧 Installing 'mockery'..."
    go install github.com/vektra/mockery/v3@latest

    if [ $? -eq 0 ]; then
        echo "✅ 'mockery' installed successfully!"
    else
        echo "❌ Failed to install 'mockery'. Please try installing manually:"
        echo "  go install github.com/vektra/mockery/v3@latest"
        exit 1
    fi
}
