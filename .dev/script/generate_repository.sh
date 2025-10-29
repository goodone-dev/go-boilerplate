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
    echo "Usage: make generate_repository <entity_name>"
    echo "Example: make generate_repository CustomerAddress (PascalCase)"
    exit 1
fi

# Get entity name and convert to various cases
ENTITY_NAME=$1
ENTITY_SNAKE=$(snake_case $ENTITY_NAME)
ENTITY_KEBAB=$(kebab_case $ENTITY_NAME)
ENTITY_CAMEL=$(camel_case $ENTITY_KEBAB)

DOMAIN_NAME=$(domain_name $ENTITY_SNAKE)

# Define directories
DOMAIN_DIR="internal/domain/${DOMAIN_NAME}"
REPO_IMPL_DIR="internal/application/${DOMAIN_NAME}/repository"

# Create directories if they don't exist
mkdir -p $DOMAIN_DIR
mkdir -p $REPO_IMPL_DIR

# Create entity file
cat > "${DOMAIN_DIR}/${ENTITY_SNAKE}.entity.go" << EOF
package ${DOMAIN_NAME}

import (
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type ${ENTITY_NAME} struct {
	database.BaseEntity[uuid.UUID] \`bson:",inline"\`
	// Add your fields here
}

func (${ENTITY_NAME}) TableName() string {
	return "${ENTITY_SNAKE}s"
}

func (${ENTITY_NAME}) RepositoryName() string {
	return "${ENTITY_NAME}Repository"
}
EOF

# Create repository interface file
cat > "${DOMAIN_DIR}/${ENTITY_SNAKE}.repository.go" << EOF
package ${DOMAIN_NAME}

import (
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type I${ENTITY_NAME}Repository interface {
	database.IBaseRepository[gorm.DB, uuid.UUID, ${ENTITY_NAME}]
}
EOF

# Create repository implementation file
cat > "${REPO_IMPL_DIR}/${ENTITY_SNAKE}.repository.go" << EOF
package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/${DOMAIN_NAME}"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ${ENTITY_NAME}Repository struct {
	database.IBaseRepository[gorm.DB, uuid.UUID, ${DOMAIN_NAME}.${ENTITY_NAME}]
}

func New${ENTITY_NAME}Repository(baseRepo database.IBaseRepository[gorm.DB, uuid.UUID, ${DOMAIN_NAME}.${ENTITY_NAME}]) ${DOMAIN_NAME}.I${ENTITY_NAME}Repository {
	return &${ENTITY_NAME}Repository{
		baseRepo,
	}
}
EOF

echo "✅ Generated files for ${ENTITY_NAME} entity and repository"
echo "- ${DOMAIN_DIR}/${ENTITY_SNAKE}.entity.go"
echo "- ${DOMAIN_DIR}/${ENTITY_SNAKE}.repository.go"
echo "- ${REPO_IMPL_DIR}/${ENTITY_SNAKE}.repository.go"
echo ""
echo "Don't forget to:"
echo "1. Add your fields to the entity struct"
echo "2. Implement any additional repository methods"
echo "3. Setup repository in main.go"
echo "4. Add the repository to your dependency injection/use cases"
