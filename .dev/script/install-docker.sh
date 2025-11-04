#!/bin/bash

FORCED=false
DEBUG=false

while getopts ":fd" opt; do
    case $opt in
        f) FORCED=true;;
        d) DEBUG=true;;
    esac
done

if ! command -v docker-compose &> /dev/null; then
    if [ "$FORCED" = true ]; then
        install_docker_compose
    else
        echo "❌ Error: 'docker-compose' is not installed."
        echo ""
        echo "🤔 Would you like to install 'docker-compose'? (y/n)"
        read -r response

        if [[ "$response" =~ ^[Yy]$ ]]; then
            install_docker_compose
        else
            echo "⏸️ Installation cancelled. To install 'docker-compose' later, run:"
            echo "  brew install docker-compose"
            exit 1
        fi
    fi
else
    if [ "$DEBUG" = true ]; then
        echo "✅ 'docker-compose' is already installed."
    fi
fi

install_docker_compose() {
    echo "🔧 Installing 'docker-compose'..."
    brew install docker-compose

    if [ $? -eq 0 ]; then
        echo "✅ 'docker-compose' installed successfully!"
    else
        echo "❌ Failed to install 'docker-compose'. Please try installing manually:"
        echo "  brew install docker-compose"
        exit 1
    fi
}
