#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | awk '/=/ {print $1}')
fi

# Get the driver from the first argument
DRIVER=$1

show_usage() {
    echo "Usage: make seed DRIVER=<database_driver>"
    echo "Example: make seed DRIVER=postgres"
    echo "
Available database drivers:"
    echo "  - postgres    : PostgreSQL database"
    echo "  - mysql       : MySQL database"
    echo "  - mongodb     : MongoDB database"
    exit 1
}

# Check if driver is provided
if [ -z "$DRIVER" ]; then
    echo "Error: Database driver are required"
    show_usage
    exit 1
fi

# Validate driver
if [ "$DRIVER" = "postgres" ]; then
    if ! command -v psql &> /dev/null
    then
        echo "psql could not be found, please install it first"
        exit
    fi
elif [ "$DRIVER" = "mysql" ]; then
    if ! command -v mysql &> /dev/null
    then
        echo "mysql could not be found, please install it first"
        exit
    fi
elif [ "$DRIVER" = "mongo" ]; then
    if ! command -v mongosh &> /dev/null
    then
        echo "mongosh could not be found, please install it first"
        exit
    fi
fi

# Check if driver is postgres
if [ "$DRIVER" = "postgres" ]; then
    # Get all seeder files for postgres
    SEEDER_FILES=$(ls seeders/postgres/*.sql)

    # Loop through each seeder file and apply it
    for f in $SEEDER_FILES
    do
        echo "Applying seeder $f"
        PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USERNAME -d $POSTGRES_DATABASE < "$f"
    done
elif [ "$DRIVER" = "mysql" ]; then
    # Get all seeder files for mysql
    SEEDER_FILES=$(ls seeders/mysql/*.sql)

    # Loop through each seeder file and apply it
    for f in $SEEDER_FILES
    do
        echo "Applying seeder $f"
        mysql -h $MYSQL_HOST -P $MYSQL_PORT -u $MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE < "$f"
    done
elif [ "$DRIVER" = "mongo" ]; then
    # Get all seeder files for mongo
    SEEDER_FILES=$(ls seeders/mongo/*.js)

    # Loop through each seeder file and apply it
    for f in $SEEDER_FILES
    do
        echo "Applying seeder $f"
        mongosh --host $MONGO_HOST --port $MONGO_PORT --username $MONGO_USER --password $MONGO_PASSWORD --authenticationDatabase admin $MONGO_DATABASE < "$f"
    done
else
    echo "Driver $DRIVER is not supported"
    exit 1
fi
