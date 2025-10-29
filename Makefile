migration:
	@.dev/script/migration.sh -n $(NAME) -d $(DRIVER)

migration_up:
	@.dev/script/migration_up.sh -d $(DRIVER)

migration_down:
	@.dev/script/migration_down.sh -d $(DRIVER)

mock:
	@.dev/script/mock.sh

mock_add:
	@.dev/script/mock_add.sh -n $(NAME)

seeder:
	@.dev/script/seeder.sh -n $(NAME) -d $(DRIVER)

seeder_up:
	@.dev/script/seeder_up.sh -d $(DRIVER)

run:
	@.dev/script/run.sh

watch:
	@.dev/script/run.sh -w

up:
	@.dev/script/docker_up.sh

down:
	@.dev/script/docker_down.sh

stop:
	@.dev/script/docker_stop.sh

setup:
	@echo "🔧 Making all .sh files in .dev directory executable..."
	@chmod +x .dev/script/*.sh
	@echo "✅ All .sh files in .dev directory are executable"
	@echo "🔧 Installing pre-commit..."
	@pre-commit install
	@echo "✅ pre-commit installed"

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Migration targets:"
	@echo "  migration NAME=<migration_name> DRIVER=<database_name>    	Create a new migration file"
	@echo "  migration_up DRIVER=<database_name>                       	Run all migrations"
	@echo "  migration_down DRIVER=<database_name>                     	Rollback the last migration"
	@echo ""
	@echo "Seeder targets:"
	@echo "  seeder NAME=<seeder_name> DRIVER=<database_name>        	Create a new seeder file"
	@echo "  seeder_up DRIVER=<database_name>                         	Apply all seeders"
	@echo ""
	@echo "Development targets:"
	@echo "  setup                                              		Make all .sh files in .dev directory executable"
	@echo "  run                                               		Run the application"
	@echo "  watch                                             		Run the application with live reloading"
	@echo ""
	@echo "Mock targets:"
	@echo "  mock                                              		Generate mocks"
	@echo "  mock_add NAME=<interface_name>                     		Add mock config"
	@echo ""
	@echo "Docker targets:"
	@echo "  up                                               		Start the application with docker-compose"
	@echo "  down                                             		Remove the application with docker-compose"
	@echo "  stop                                             		Stop the application with docker-compose"

.PHONY: help setup run watch \
		migration migration_up migration_down \
		seeder seeder_up \
		mock mock_add \
		up down stop
