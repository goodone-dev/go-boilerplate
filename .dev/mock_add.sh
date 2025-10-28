#!/bin/bash

show_usage() {
    echo "Usage: make mock_add NAME=<interface_name>"
    echo "Example: make mock_add NAME=ICustomerRepository"
    echo ""
    exit 1
}

# Parse command line arguments
while getopts ":n:h" opt; do
    case $opt in
        n) INTERFACE_NAME="$OPTARG";;
        h) show_usage;;
    esac
done

# Validate required arguments
if [ -z "$INTERFACE_NAME" ]; then
    echo "Error: Interface name are required"
    show_usage
fi

# Check if 'mockery' is installed, and install it if not
$(dirname "$0")/ensure_mockery.sh

MODULE_PATH=$(head -n 1 go.mod | sed 's/module //')
FILE_PATH=$(grep -rl "type ${INTERFACE_NAME} interface" internal | head -n 1)

if [ -z "$FILE_PATH" ];
then
    echo "Interface ${INTERFACE_NAME} not found in internal."
    exit 1
fi

PACKAGE_DIR=$(dirname ${FILE_PATH})
PACKAGE_PATH="${MODULE_PATH}/${PACKAGE_DIR}"
BASE_FILENAME=$(basename ${FILE_PATH} .go)
MOCK_FILENAME="${BASE_FILENAME}_mock.go"

# Check if the package is already configured
if grep -q "  ${PACKAGE_PATH}:" .mockery.yml;
then
    # If package exists, check if the interface is already configured
    if grep -A 1 "  ${PACKAGE_PATH}:" .mockery.yml | grep -q "    interfaces:" && \
       grep -A 5 "    interfaces:" .mockery.yml | grep -q "      ${INTERFACE_NAME}:";
    then
        echo "Interface ${INTERFACE_NAME} is already configured in .mockery.yml."
        exit 0
    fi
fi

YAML_CONFIG="  ${PACKAGE_PATH}:
    interfaces:
      ${INTERFACE_NAME}:
        config:
          dir: \"{{.InterfaceDir}}/mocks\"
          filename: \"${MOCK_FILENAME}\""

if ! grep -q "packages:" .mockery.yml;
then
    echo -e "\npackages:" >> .mockery.yml
fi

echo -e "\n$YAML_CONFIG" >> .mockery.yml
echo "Added mock configuration for ${INTERFACE_NAME} to .mockery.yml."
