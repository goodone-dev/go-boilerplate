.PHONY: migration
migration:
	@.dev/migration_create.sh -n $(NAME) -d $(DRIVER)

.PHONY: migration_up
migration_up:
	@.dev/migration_up.sh -d $(DRIVER)

.PHONY: migration_down
migration_down:
	@.dev/migration_down.sh -d $(DRIVER)

.PHONY: mock
mock:
	@mockery --log-level=ERROR

.PHONY: mock_config
mock_config:
	@bash .dev/mock_config.sh $(NAME)

.PHONY: seed
seed:
	@.dev/seed.sh $(DRIVER)

.PHONY: run
run:
	@go run ./cmd/api/main.go || true