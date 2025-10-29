#!/bin/bash

# Script to apply database seeders using golang-migrate

# Function to show usage
show_usage() {
    echo "Usage: make seeder_up DRIVER=<database_driver>"
    echo "Example: make seeder_up DRIVER=postgres"
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
    echo "❌ Error: Database driver is required"
    show_usage
fi

# Load environment variables if .env file exists
if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Set database URL and seeder directory based on driver
case $DB_DRIVER in
    postgres|postgresql)
        SEEDER_DIR="./seeders/postgres"
        # Check required environment variables
        required_vars=("POSTGRES_HOST" "POSTGRES_PORT" "POSTGRES_USERNAME" "POSTGRES_PASSWORD" "POSTGRES_SSL_MODE" "POSTGRES_DATABASE")
        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                echo "❌ Error: Required environment variable $var is not set"
                exit 1
            fi
        done
        # Construct database URL from environment variables
        DB_URL="postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_SSL_MODE}"
        ;;
    mysql)
        SEEDER_DIR="./seeders/mysql"
        # Check required environment variables
        required_vars=("MYSQL_HOST" "MYSQL_PORT" "MYSQL_USERNAME" "MYSQL_PASSWORD" "MYSQL_DATABASE")
        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                echo "❌ Error: Required environment variable $var is not set"
                exit 1
            fi
        done
        # Construct database URL from environment variables
        DB_URL="mysql://${MYSQL_USERNAME}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}"
        ;;
    mongodb)
        SEEDER_DIR="./seeders/mongodb"
        # Check required environment variables
        required_vars=("MONGODB_HOST" "MONGODB_PORT" "MONGODB_USERNAME" "MONGODB_PASSWORD" "MONGODB_SSL_MODE" "MONGODB_DATABASE")
        for var in "${required_vars[@]}"; do
            if [ -z "${!var}" ]; then
                echo "❌ Error: Required environment variable $var is not set"
                exit 1
            fi
        done
        # Construct database URL from environment variables
        DB_URL="mongodb://${MONGODB_USERNAME}:${MONGODB_PASSWORD}@${MONGODB_HOST}:${MONGODB_PORT}/${MONGODB_DATABASE}"
        ;;
    *)
        echo "❌ Error: Unsupported database driver: $DB_DRIVER"
        show_usage
        ;;
esac

# Check if seeder directory exists
if [ ! -d "$SEEDER_DIR" ]; then
    echo "Error: Seeder directory not found: $SEEDER_DIR"
    exit 1
fi

# Check if 'golang-migrate' is installed, and install it if not
$(dirname "$0")/ensure_migrate.sh

# Apply seeders
echo "🌱 Applying seeders for $DB_DRIVER..."
migrate -database "$DB_URL" -path "$SEEDER_DIR" up

# Check if seeders applied successfully
if [ $? -eq 0 ]; then
    echo "✅ Seeders applied successfully!"
else
    echo "❌ Error: Failed to apply seeders"
    exit 1
fi
