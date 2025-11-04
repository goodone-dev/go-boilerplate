#!/bin/bash

FORCED=false
DEBUG=false

while getopts ":fd" opt; do
    case $opt in
        f) FORCED=true;;
        d) DEBUG=true;;
    esac
done

if ! command -v pre-commit &> /dev/null; then
    if [ "$FORCED" = true ]; then
        install_pre_commit
    else
        echo "❌ Error: 'pre-commit' is not installed."
        echo ""
        echo "🤔 Would you like to install 'pre-commit'? (y/n)"
        read -r response

        if [[ "$response" =~ ^[Yy]$ ]]; then
            install_pre_commit
        else
            echo "⏸️ Installation cancelled. To install 'pre-commit' later, run:"
            echo "  brew install pre-commit"
            exit 1
        fi
    fi
else
    if [ "$DEBUG" = true ]; then
        echo "✅ 'pre-commit' is already installed."
    fi
fi

install_pre_commit() {
    echo "🔧 Installing 'pre-commit'..."
    brew install pre-commit

    if [ $? -eq 0 ]; then
        echo "✅ 'pre-commit' installed successfully!"
    else
        echo "❌ Failed to install 'pre-commit'. Please try installing manually:"
        echo "  brew install pre-commit"
        exit 1
    fi
}

if [ -f ".git/hooks/pre-commit" ]; then
    if [ "$DEBUG" = true ]; then
        echo "✅ pre-commit hooks are already installed."
    fi
else
    echo "🔧 Installing pre-commit hooks..."
    pre-commit install

    if [ $? -eq 0 ]; then
        echo "✅ pre-commit hooks installed successfully!"
    else
        echo "❌ Failed to install pre-commit hooks. Please try installing manually:"
        echo "  pre-commit install"
        exit 1
    fi
fi
