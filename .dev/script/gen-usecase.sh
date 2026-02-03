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
    echo "Usage: make gen-usecase NAME=<usecase_name>"
    echo "Example: make gen-usecase NAME=CustomerAddress"
    exit 1
fi

# Get usecase name and convert to various cases
USECASE_NAME=$1
USECASE_SNAKE=$(snake_case $USECASE_NAME)
USECASE_KEBAB=$(kebab_case $USECASE_NAME)
USECASE_CAMEL=$(camel_case $USECASE_KEBAB)
DOMAIN_NAME=$(domain_name $USECASE_SNAKE)

# Define directories
DOMAIN_DIR="internal/domain/${DOMAIN_NAME}"
USECASE_IMPL_DIR="internal/application/${DOMAIN_NAME}/usecase"

# Create directories if they don't exist
mkdir -p $DOMAIN_DIR
mkdir -p $USECASE_IMPL_DIR

# Create usecase interface file
cat > "${DOMAIN_DIR}/${USECASE_SNAKE}.usecase.go" << EOF
package ${DOMAIN_NAME}

type ${USECASE_NAME}Usecase interface {
	// Define your usecase methods here
}
EOF

# Create usecase implementation file
cat > "${USECASE_IMPL_DIR}/${USECASE_SNAKE}.usecase.go" << EOF
package usecase

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/${DOMAIN_NAME}"
)

type ${USECASE_CAMEL}Usecase struct {
	${USECASE_CAMEL}Repo ${DOMAIN_NAME}.${USECASE_NAME}Repository
}

func New${USECASE_NAME}Usecase(${USECASE_CAMEL}Repo ${DOMAIN_NAME}.${USECASE_NAME}Repository) ${DOMAIN_NAME}.${USECASE_NAME}Usecase {
	return &${USECASE_CAMEL}Usecase{
		${USECASE_CAMEL}Repo: ${USECASE_CAMEL}Repo,
	}
}
EOF

echo "✅ Generated files for ${USECASE_NAME} usecase"
echo "- ${DOMAIN_DIR}/${USECASE_SNAKE}.usecase.go"
echo "- ${USECASE_IMPL_DIR}/${USECASE_SNAKE}.usecase.go"
echo ""
echo "Don't forget to:"
echo "1. Define your usecase methods in the interface"
echo "2. Implement your usecase methods"
echo "3. Setup usecase in main.go"
echo "4. Register the usecase in your dependency injection"
