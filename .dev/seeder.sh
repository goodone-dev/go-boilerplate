#!/bin/bash

# Script to create database seeder files using golang-migrate

# Function to show usage
show_usage() {
    echo "Usage: make seeder NAME=<seeder_name> DRIVER=<database_driver>"
    echo "Example: make seeder NAME=seed_users_table DRIVER=postgres"
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
        n) SEEDER_NAME="$OPTARG";;
        d) DB_DRIVER="$OPTARG";;
        h) show_usage;;
    esac
done

# Validate required arguments
if [ -z "$SEEDER_NAME" ] || [ -z "$DB_DRIVER" ]; then
    echo "❌ Error: Both seeder name and database driver are required"
    show_usage
fi

# Validate seeder name (allow only lowercase letters, numbers, and underscores)
if ! [[ $SEEDER_NAME =~ ^[a-z0-9_]+$ ]]; then
    echo "❌ Error: Seeder name can only contain lowercase letters, numbers, and underscores"
    exit 1
fi

# Validate database driver
case $DB_DRIVER in
    postgres|postgresql) 
        SEEDER_DIR="./seeders/postgres"
        ;;
    mysql) 
        SEEDER_DIR="./seeders/mysql"
        ;;
    mongodb) 
        SEEDER_DIR="./seeders/mongodb"
        ;;
    *)
        echo "❌ Error: Unsupported database driver: $DB_DRIVER"
        show_usage
        ;;
esac

# Create seeders directory if it doesn't exist
mkdir -p $SEEDER_DIR

# Check if 'golang-migrate' is installed, and install it if not
$(dirname "$0")/ensure_migrate.sh

# Create seeder files
echo "🌱 Creating seeder files for $SEEDER_NAME..."
migrate create -ext sql -dir $SEEDER_DIR -format "20060102150405" -tz "Asia/Jakarta" $SEEDER_NAME

if [ $? -eq 0 ]; then
    echo "✅ Seeder files created successfully!"
else
    echo "❌ Error: Failed to create seeder files"
    exit 1
fi
