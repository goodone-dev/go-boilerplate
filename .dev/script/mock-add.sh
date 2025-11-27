#!/bin/bash

show_usage() {
    echo "Usage: make mock-add NAME=<interface_name>"
    echo "Example: make mock-add NAME=CustomerRepository"
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
    echo "❌ Error: Interface name is required"
    show_usage
fi

MODULE_PATH=$(head -n 1 go.mod | sed 's/module //')
FILE_PATH=$(grep -rl "type ${INTERFACE_NAME} interface" internal | head -n 1)

if [ -z "$FILE_PATH" ];
then
    echo "🔍 Interface ${INTERFACE_NAME} not found in internal directory"
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
    if grep -A 100 "  ${PACKAGE_PATH}:" .mockery.yml | grep -q "      ${INTERFACE_NAME}:";
    then
        echo "🚫 Interface ${INTERFACE_NAME} is already configured in .mockery.yml."
        exit 0
    fi

    # Package exists, add interface to existing package
    # Find the line number after "interfaces:" within the package and insert there
    awk -v pkg="  ${PACKAGE_PATH}:" -v iname="${INTERFACE_NAME}" -v fname="${MOCK_FILENAME}" '
    BEGIN { found_pkg=0; inserted=0 }
    {
        print $0
        if ($0 == pkg) {
            found_pkg=1
        }
        # Insert after "interfaces:" line within the package
        if (found_pkg && !inserted && $0 ~ /^    interfaces:/) {
            print "      " iname ":"
            print "        config:"
            print "          dir: \"{{.InterfaceDir}}/mocks\""
            print "          filename: \"" fname "\""
            inserted=1
            found_pkg=0
        }
    }
    ' .mockery.yml > .mockery.yml.tmp && mv .mockery.yml.tmp .mockery.yml

    echo "✅ Added interface ${INTERFACE_NAME} to existing package ${PACKAGE_PATH} in .mockery.yml."
else
    # Package doesn't exist, create new package entry
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
    echo "✅ Added mock configuration for ${INTERFACE_NAME} to .mockery.yml."
fi
