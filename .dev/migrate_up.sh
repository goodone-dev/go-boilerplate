#!/bin/bash

# Script to apply database migrations using golang-migrate

# Function to show usage
show_usage() {
    echo "Usage: make migrate_up DRIVER=<database_driver>"
    echo "Example: make migrate_up DRIVER=postgres"
    echo "
Available database drivers:"
    echo "  - postgres    : PostgreSQL database"
    echo "  - mysql       : MySQL database"
    echo "  - mongodb     : MongoDB database"
    exit 1
}

# Parse command line arguments
while getopts ":d:h" opt; do
    case $opt in
        d) DB_DRIVER="$OPTARG";;
        h) show_usage;;
    esac
done

# Validate required arguments
if [ -z "$DB_DRIVER" ]; then
    echo "Error: Database driver is required"
    show_usage
fi

# Load environment variables if .env file exists
if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Set database URL and migration directory based on driver
case $DB_DRIVER in
    postgres|postgresql) 
        MIGRATION_DIR="./migrations/postgres"
        # Check required environment variables
        required_vars=("POSTGRES_HOST" "POSTGRES_PORT" "POSTGRES_USERNAME" "POSTGRES_PASSWORD" "POSTGRES_SSL_MODE" "POSTGRES_DATABASE")
        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                echo "Error: Required environment variable $var is not set"
                exit 1
            fi
        done
        # Construct database URL from environment variables
        DB_URL="postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_SSL_MODE}"
        ;;
    mysql) 
        MIGRATION_DIR="./migrations/mysql"
        # Check required environment variables
        required_vars=("MYSQL_HOST" "MYSQL_PORT" "MYSQL_USERNAME" "MYSQL_PASSWORD" "MYSQL_DATABASE")
        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                echo "Error: Required environment variable $var is not set"
                exit 1
            fi
        done
        # Construct database URL from environment variables
        DB_URL="mysql://${MYSQL_USERNAME}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}"
        ;;
    mongodb) 
        MIGRATION_DIR="./migrations/mongodb"
        # Check required environment variables
        required_vars=("MONGODB_HOST" "MONGODB_PORT" "MONGODB_USERNAME" "MONGODB_PASSWORD" "MONGODB_SSL_MODE" "MONGODB_DATABASE")
        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                echo "Error: Required environment variable $var is not set"
                exit 1
            fi
        done
        # Construct database URL from environment variables
        DB_URL="mongodb://${MONGODB_USERNAME}:${MONGODB_PASSWORD}@${MONGODB_HOST}:${MONGODB_PORT}/${MONGODB_DATABASE}"
        ;;
    *)
        echo "Error: Unsupported database driver: $DB_DRIVER"
        show_usage
        ;;
esac

# Check if migration directory exists
if [ ! -d "$MIGRATION_DIR" ]; then
    echo "Error: Migration directory not found: $MIGRATION_DIR"
    exit 1
fi

# Check if golang-migrate is installed
if ! command -v migrate &> /dev/null; then
    echo "Error: 'migrate' is not installed."
    echo ""
    echo "Would you like to install 'golang-migrate'? (y/n)"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo ""
        echo "Installing 'golang-migrate'..."
        go install -tags "$DB_DRIVER" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
        
        if [ $? -eq 0 ]; then
            echo ""
            echo "✓ 'migrate' installed successfully!"
            echo ""
        else
            echo ""
            echo "✗ Failed to install 'migrate'. Please try installing manually:"
            echo "  go install -tags '$DB_DRIVER' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
            exit 1
        fi
    else
        echo ""
        echo "Installation cancelled. To install 'migrate' later, run:"
        echo "  go install -tags '$DB_DRIVER' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
fi

# Apply migrations
echo "Applying migrations for $DB_DRIVER..."
migrate -database "$DB_URL" -path "$MIGRATION_DIR" up

if [ $? -eq 0 ]; then
    echo "✅ Migrations applied successfully"
else
    echo "❌ Failed to apply migrations"
    exit 1
fi
