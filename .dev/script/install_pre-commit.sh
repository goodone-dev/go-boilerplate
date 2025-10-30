#!/bin/bash

# Script to install pre-commit if it is not installed

if ! command -v pre-commit &> /dev/null; then
    echo "🔧 Installing 'pre-commit'..."
    brew install pre-commit

    if [ $? -eq 0 ]; then
        echo "✅ 'pre-commit' installed successfully!"
    else
        echo "❌ Failed to install 'pre-commit'. Please try installing manually:"
        echo "  brew install pre-commit"
        exit 1
    fi
fi

echo "🔧 Installing pre-commit hooks..."
pre-commit install

if [ $? -eq 0 ]; then
    echo "✅ pre-commit hooks installed successfully!"
else
    echo "❌ Failed to install pre-commit hooks. Please try installing manually:"
    echo "  pre-commit install"
    exit 1
fi
