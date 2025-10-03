#!/bin/bash

# Script to create database migration files using golang-migrate

# Function to show usage
show_usage() {
    echo "Usage: $0 -n <migration_name> -d <database_driver>"
    echo "Example: $0 -n create_users_table -d postgres"
    echo "
Available database drivers:"
    echo "  - postgres    : PostgreSQL database"
    echo "  - mysql       : MySQL database"
    echo "  - mongodb     : MongoDB database"
    exit 1
}

# Parse command line arguments
while getopts ":n:d:h" opt; do
    case $opt in
        n) MIGRATION_NAME="$OPTARG";;
        d) DB_DRIVER="$OPTARG";;
        h) show_usage;;
        \?) echo "Invalid option -$OPTARG"; show_usage;;
        :) echo "Option -$OPTARG requires an argument."; show_usage;;
    esac
done

# Validate required arguments
if [ -z "$MIGRATION_NAME" ] || [ -z "$DB_DRIVER" ]; then
    echo "Error: Both migration name (-n) and database driver (-d) are required"
    show_usage
fi

# Validate migration name (allow only lowercase letters, numbers, and underscores)
if ! [[ $MIGRATION_NAME =~ ^[a-z0-9_]+$ ]]; then
    echo "Error: Migration name can only contain lowercase letters, numbers, and underscores"
    exit 1
fi

# Validate database driver
case $DB_DRIVER in
    postgres|postgresql) 
        MIGRATION_DIR="./migrations/postgres"
        ;;
    mysql) 
        MIGRATION_DIR="./migrations/mysql"
        ;;
    mongodb) 
        MIGRATION_DIR="./migrations/mongodb"
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

# Create migrations directory if it doesn't exist
mkdir -p $MIGRATION_DIR

# Create migration files
echo "Creating migration files..."
migrate create -ext sql -dir $MIGRATION_DIR -format "20060102150405" -tz "Asia/Jakarta" $MIGRATION_NAME

if [ $? -eq 0 ]; then
    echo "✅ Migration files created successfully in $MIGRATION_DIR directory"
else
    echo "❌ Failed to create migration files"
    exit 1
fi
