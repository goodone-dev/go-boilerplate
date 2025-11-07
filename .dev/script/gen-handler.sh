#!/bin/bash

# Function to convert to camel case (preserve first character case)
camel_case() {
    echo $1 | awk -F"[-_]" '{$1=tolower($1); for(i=2; i<=NF; i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))} 1' OFS=""
}

# Function to convert to snake case
snake_case() {
    echo $1 | sed -r 's/([a-z0-9])([A-Z])/\1_\2/g' | tr '[:upper:]' '[:lower:]'
}

# Function to convert to kebab case
kebab_case() {
    echo $1 | sed -r 's/([a-z0-9])([A-Z])/\1-\2/g' | tr '[:upper:]' '[:lower:]'
}

# Function to get first word (for domain)
domain_name() {
    echo $1 | sed -E 's/[_-].*//' | tr '[:upper:]' '[:lower:]'
}

# Check if name is provided
if [ -z "$1" ]; then
    echo "Usage: make gen-handler NAME=<handler_name>"
    echo "Example: make gen-handler NAME=CustomerAddress"
    exit 1
fi

# Get handler name and convert to various cases
HANDLER_NAME=$1
HANDLER_SNAKE=$(snake_case $HANDLER_NAME)
HANDLER_KEBAB=$(kebab_case $HANDLER_NAME)
HANDLER_CAMEL=$(camel_case $HANDLER_KEBAB)
DOMAIN_NAME=$(domain_name $HANDLER_SNAKE)

# Define directories
DOMAIN_DIR="internal/domain/${DOMAIN_NAME}"
HANDLER_IMPL_DIR="internal/application/${DOMAIN_NAME}/delivery/http"

# Create directories if they don't exist
mkdir -p $DOMAIN_DIR
mkdir -p $HANDLER_IMPL_DIR

# Create handler interface file
cat > "${DOMAIN_DIR}/${HANDLER_SNAKE}.handler.go" << EOF
package ${DOMAIN_NAME}

type ${HANDLER_NAME}Handler interface {
	// Define your handler methods here
}
EOF

# Create handler implementation file
cat > "${HANDLER_IMPL_DIR}/${HANDLER_SNAKE}.handler.go" << EOF
package http

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/${DOMAIN_NAME}"
)

type ${HANDLER_CAMEL}Handler struct {
	${HANDLER_CAMEL}Usecase ${DOMAIN_NAME}.${HANDLER_NAME}Usecase
}

func New${HANDLER_NAME}Handler(${HANDLER_CAMEL}Usecase ${DOMAIN_NAME}.${HANDLER_NAME}Usecase) ${DOMAIN_NAME}.${HANDLER_NAME}Handler {
	return &${HANDLER_CAMEL}Handler{
		${HANDLER_CAMEL}Usecase: ${HANDLER_CAMEL}Usecase,
	}
}
EOF

echo "✅ Generated files for ${HANDLER_NAME} handler"
echo "- ${DOMAIN_DIR}/${HANDLER_SNAKE}.handler.go"
echo "- ${HANDLER_IMPL_DIR}/${HANDLER_SNAKE}.handler.go"
echo ""
echo "Don't forget to:"
echo "1. Define your handler methods in the interface"
echo "2. Implement your handler methods"
echo "3. Setup handler in main.go"
echo "4. Register the handler in your dependency injection"
