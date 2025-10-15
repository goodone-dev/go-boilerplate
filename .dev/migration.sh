#!/bin/bash

# Script to create database migration files using golang-migrate

# Function to show usage
show_usage() {
    echo "Usage: make migration NAME=<migration_name> DRIVER=<database_driver>"
    echo "Example: make migration NAME=create_users_table DRIVER=postgres"
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
    esac
done

# Validate required arguments
if [ -z "$MIGRATION_NAME" ] || [ -z "$DB_DRIVER" ]; then
    echo "Error: Both migration name and database driver are required"
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

# Create migrations directory if it doesn't exist
mkdir -p $MIGRATION_DIR

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

# Create migration files
echo "Creating migration files..."
migrate create -ext sql -dir $MIGRATION_DIR -format "20060102150405" -tz "Asia/Jakarta" $MIGRATION_NAME

if [ $? -eq 0 ]; then
    echo "Migration files created successfully in $MIGRATION_DIR directory"
else
    echo "Failed to create migration files"
    exit 1
fi
