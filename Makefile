migration:
	@.dev/migration.sh -n $(NAME) -d $(DRIVER)

migrate_up:
	@.dev/migrate_up.sh -d $(DRIVER)

migrate_down:
	@.dev/migrate_down.sh -d $(DRIVER)

mock_add:
	@.dev/mock_add.sh -n $(NAME)

mock:
	@.dev/mock.sh

seed:
	@.dev/seed.sh -d $(DRIVER)

seeder:
	@.dev/seeder.sh -n $(NAME) -d $(DRIVER)

run:
	@.dev/run.sh

watch:
	@.dev/run.sh -w

up:
	@docker-compose up --build -d

down:
	@docker-compose down

stop:
	@docker-compose stop

setup:
	@chmod +x .dev/*.sh

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Migration targets:"
	@echo "  migration NAME=<migration_name> DRIVER=<driver>    	Create a new migration file"
	@echo "  migrate_up DRIVER=<driver>                       	Run all migrations"
	@echo "  migrate_down DRIVER=<driver>                     	Rollback the last migration"
	@echo ""
	@echo "Seeder targets:"
	@echo "  seeder NAME=<seeder_name> DRIVER=<driver>        	Create a new seeder file"
	@echo "  seed DRIVER=<driver>                             	Apply all seeders"
	@echo ""
	@echo "Development targets:"
	@echo "  setup                                              	Make all .sh files in .dev directory executable"
	@echo "  run                                               	Run the application"
	@echo "  watch                                             	Run the application with live reloading"
	@echo "  mock                                              	Generate mocks"
	@echo "  mock_add NAME=<interface_name>                     	Add mock config"
	@echo ""
	@echo "Docker targets:"
	@echo "  up                                               	Start the application with docker-compose"
	@echo "  down                                             	Remove the application with docker-compose"
	@echo "  stop                                             	Stop the application with docker-compose"

.PHONY: help setup run watch \
		migration migrate_up migrate_down \
		seeder seed \
		mock mock_add \
		up down stop