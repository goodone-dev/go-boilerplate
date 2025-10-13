migration:
	@.dev/migration.sh -n $(NAME) -d $(DRIVER)

migrate_up:
	@.dev/migrate_up.sh -d $(DRIVER)

migrate_down:
	@.dev/migrate_down.sh -d $(DRIVER)

mock:
	@mockery --log-level=ERROR

mock_config:
	@.dev/mock_config.sh $(NAME)

seed:
	@.dev/seed.sh $(DRIVER)

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

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Migration targets:"
	@echo "  migration NAME=<migration_name> DRIVER=<driver>    	Create a new migration file"
	@echo "  migrate_up DRIVER=<driver>                       	Run all migrations"
	@echo "  migrate_down DRIVER=<driver>                     	Rollback the last migration"
	@echo "  seed DRIVER=<driver>                             	Seed the database"
	@echo ""
	@echo "Development targets:"
	@echo "  run                                               	Run the application"
	@echo "  watch                                             	Run the application with live reloading"
	@echo "  mock                                              	Generate mocks"
	@echo "  mock_config NAME=<name>                          	Generate mock config"
	@echo ""
	@echo "Docker targets:"
	@echo "  up                                               	Start the application with docker-compose"
	@echo "  down                                             	Remove the application with docker-compose"
	@echo "  stop                                             	Stop the application with docker-compose"

.PHONY: help run watch \
		migration migrate_up migrate_down seed \
		mock mock_config \
		up down stop