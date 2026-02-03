setup:
	@chmod -v +x .dev/script/*.sh
	@echo "✅ All .sh files in .dev/script directory are executable"
	@.dev/script/install-air.sh -f -v
	@.dev/script/install-docker.sh -f -v
	@.dev/script/install-test-coverage.sh -f -v
	@.dev/script/install-migrate.sh -f -v
	@.dev/script/install-mockery.sh -f -v
	@.dev/script/install-pre-commit.sh -f -v

install-air:
	@.dev/script/install-air.sh

install-docker:
	@.dev/script/install-docker.sh

install-test-coverage:
	@.dev/script/install-test-coverage.sh

install-migrate:
	@.dev/script/install-migrate.sh

install-mockery:
	@.dev/script/install-mockery.sh

install-pre-commit:
	@.dev/script/install-pre-commit.sh

db-migrate-new: install-migrate
	@.dev/script/db-migrate-new.sh -n $(NAME) -d $(DRIVER)

db-migrate-up: install-migrate
	@.dev/script/db-migrate-up.sh -d $(DRIVER)

db-migrate-down: install-migrate
	@.dev/script/db-migrate-down.sh -d $(DRIVER)

mock: install-mockery
	@.dev/script/mock.sh

mock-add: install-mockery
	@.dev/script/mock-add.sh -n $(NAME)

test: install-test-coverage
	@.dev/script/test.sh

db-seed-new: install-migrate
	@.dev/script/db-seed-new.sh -n $(NAME) -d $(DRIVER)

db-seed-up: install-migrate
	@.dev/script/db-seed-up.sh -d $(DRIVER)

run:
	@.dev/script/run.sh

watch: install-air
	@.dev/script/run.sh -w

docker-up: install-docker
	@.dev/script/docker-up.sh

docker-down: install-docker
	@.dev/script/docker-down.sh

docker-stop: install-docker
	@.dev/script/docker-stop.sh

gen-repo:
	@.dev/script/gen-repo.sh $(NAME)

gen-usecase:
	@.dev/script/gen-usecase.sh $(NAME)

gen-handler:
	@.dev/script/gen-handler.sh $(NAME)

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Setup targets:"
	@echo "  setup                                             Setup development environment"
	@echo ""
	@echo "Development targets:"
	@echo "  run                                               Run application"
	@echo "  watch                                             Run application with live reload"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-up                                         Start containers"
	@echo "  docker-down                                       Stop and remove containers"
	@echo "  docker-stop                                       Stop containers"
	@echo ""
	@echo "Layer generation targets:"
	@echo "  gen-repo NAME=<name>                              Generate repository layer"
	@echo "  gen-usecase NAME=<name>                           Generate usecase layer"
	@echo "  gen-handler NAME=<name>                           Generate handler layer"
	@echo ""
	@echo "Migration targets:"
	@echo "  db-migrate-new NAME=<name> DRIVER=<driver>        Create new migration file"
	@echo "  db-migrate-up DRIVER=<driver>                     Apply all pending migrations"
	@echo "  db-migrate-down DRIVER=<driver>                   Rollback last migration"
	@echo ""
	@echo "Seeder targets:"
	@echo "  db-seed-new NAME=<name> DRIVER=<driver>           Create new seeder file"
	@echo "  db-seed-up DRIVER=<driver>                        Apply all seeders"
	@echo ""
	@echo "Test targets:"
	@echo "  test                                              Run tests with coverage check"
	@echo ""
	@echo "Mock targets:"
	@echo "  mock                                              Generate all mocks"
	@echo "  mock-add NAME=<interface>                         Add new mock configuration"

.PHONY: help setup \
		install-air install-docker install-test-coverage install-migrate install-mockery install-pre-commit \
		run watch \
		docker-up docker-down docker-stop \
		gen-repo gen-usecase gen-handler \
		db-migrate-new db-migrate-up db-migrate-down \
		db-seed-new db-seed-up \
		test \
		mock mock-add
