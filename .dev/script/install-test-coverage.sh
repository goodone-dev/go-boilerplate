#!/bin/bash

FORCED=false
DEBUG=false

while getopts ":fd" opt; do
    case $opt in
        f) FORCED=true;;
        d) DEBUG=true;;
    esac
done

if ! command -v go-test-coverage &> /dev/null; then
    if [ "$FORCED" = true ]; then
        install_go_test_coverage
    else
        echo "❌ Error: 'go-test-coverage' is not installed."
        echo ""
        echo "🤔 Would you like to install 'go-test-coverage'? (y/n)"
        read -r response

        if [[ "$response" =~ ^[Yy]$ ]]; then
            install_go_test_coverage
        else
            echo "⏸️ Installation cancelled. To install 'go-test-coverage' later, run:"
            echo "  go install github.com/vladopajic/go-test-coverage/v2@latest"
            exit 1
        fi
    fi
else
    if [ "$DEBUG" = true ]; then
        echo "✅ 'go-test-coverage' is already installed."
    fi
fi

install_go_test_coverage() {
    echo "🔧 Installing 'go-test-coverage'..."
    go install github.com/vladopajic/go-test-coverage/v2@latest

    if [ $? -eq 0 ]; then
        echo "✅ 'go-test-coverage' installed successfully!"
    else
        echo "❌ Failed to install 'go-test-coverage'. Please try installing manually:"
        echo "  go install github.com/vladopajic/go-test-coverage/v2@latest"
        exit 1
    fi
}
