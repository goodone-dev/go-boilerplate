#!/bin/bash

# Script to apply database migrations using golang-migrate

# Function to show usage
show_usage() {
    echo "Usage: $0 -d <database_driver>"
    echo "Example: $0 -d postgres"
    echo "
Available database drivers:"
    echo "  - postgres    : PostgreSQL database"
    echo "  - mysql       : MySQL database"
    echo "  - mongodb     : MongoDB database"
    echo "  - sqlite      : SQLite database"
    exit 1
}

# Parse command line arguments
while getopts ":d:h" opt; do
    case $opt in
        d) DB_DRIVER="$OPTARG";;
        h) show_usage;;
        \?) echo "Invalid option -$OPTARG"; show_usage;;
        :) echo "Option -$OPTARG requires an argument."; show_usage;;
    esac
done

# Validate required arguments
if [ -z "$DB_DRIVER" ]; then
    echo "Error: Database driver (-d) is required"
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
        required_vars=("POSTGRES_MASTER_HOST" "POSTGRES_MASTER_PORT" "POSTGRES_MASTER_USERNAME" "POSTGRES_MASTER_PASSWORD" "POSTGRES_MASTER_SSL_MODE" "POSTGRES_DATABASE")
        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                echo "Error: Required environment variable $var is not set"
                exit 1
            fi
        done
        # Construct database URL from environment variables
        DB_URL="postgresql://${POSTGRES_MASTER_USERNAME}:${POSTGRES_MASTER_PASSWORD}@${POSTGRES_MASTER_HOST}:${POSTGRES_MASTER_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_MASTER_SSL_MODE}"
        ;;
    mysql) 
        MIGRATION_DIR="./migrations/mysql"
        # Check required environment variables
        required_vars=("MYSQL_MASTER_HOST" "MYSQL_MASTER_PORT" "MYSQL_MASTER_USERNAME" "MYSQL_MASTER_PASSWORD" "MYSQL_MASTER_SSL_MODE" "MYSQL_DATABASE")
        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                echo "Error: Required environment variable $var is not set"
                exit 1
            fi
        done
        # Construct database URL from environment variables
        DB_URL="mysql://${MYSQL_MASTER_USERNAME}:${MYSQL_MASTER_PASSWORD}@tcp(${MYSQL_MASTER_HOST}:${MYSQL_MASTER_PORT})/${MYSQL_DATABASE}"
        ;;
    mongodb) 
        MIGRATION_DIR="./migrations/mongodb"
        # Check required environment variables
        required_vars=("MONGODB_MASTER_HOST" "MONGODB_MASTER_PORT" "MONGODB_MASTER_USERNAME" "MONGODB_MASTER_PASSWORD" "MONGODB_MASTER_SSL_MODE" "MONGODB_DATABASE")
        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                echo "Error: Required environment variable $var is not set"
                exit 1
            fi
        done
        # Construct database URL from environment variables
        DB_URL="mongodb://${MONGODB_MASTER_USERNAME}:${MONGODB_MASTER_PASSWORD}@${MONGODB_MASTER_HOST}:${MONGODB_MASTER_PORT}/${MONGODB_DATABASE}"
        ;;
    sqlite) 
        MIGRATION_DIR="./migrations/sqlite"
        # Check required environment variables
        required_vars=("SQLITE_DATABASE")
        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                echo "Error: Required environment variable $var is not set"
                exit 1
            fi
        done
        # Construct database URL from environment variables
        DB_URL="sqlite3://${SQLITE_DATABASE}"
        ;;
    *)
        echo "Error: Unsupported database driver: $DB_DRIVER"
        show_usage
        ;;
esac

# Check if golang-migrate is installed
if ! command -v migrate &> /dev/null; then
    echo "Error: golang-migrate is not installed"
    echo "Install it using: go install -tags '$DB_DRIVER' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Check if migration directory exists
if [ ! -d "$MIGRATION_DIR" ]; then
    echo "Error: Migration directory not found: $MIGRATION_DIR"
    exit 1
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
