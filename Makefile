migration:
	@.dev/migration.sh -n $(NAME) -d $(DRIVER)

migrate_up:
	@.dev/migrate_up.sh -d $(DRIVER)

migrate_down:
	@.dev/migrate_down.sh -d $(DRIVER)

mock:
	@mockery --log-level=ERROR

mock_config:
	@bash .dev/mock_config.sh $(NAME)

seed:
	@.dev/seed.sh $(DRIVER)

run:
	@go run ./cmd/api/main.go

up:
	@docker-compose up --build -d

down:
	@docker-compose down

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
	@echo "  run                                              	Run the application"
	@echo "  mock                                             	Generate mocks"
	@echo "  mock_config NAME=<name>                          	Generate mock config"
	@echo ""
	@echo "Docker targets:"
	@echo "  up                                               	Start the application with docker-compose"

.PHONY: help run \
		migration migrate_up migrate_down seed \
		mock mock_config \
		up down