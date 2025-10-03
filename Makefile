.PHONY: migration
migration:
	@.dev/migration_create.sh -n $(NAME) -d $(DRIVER)

.PHONY: migration_up
migration_up:
	@.dev/migration_up.sh -d $(DRIVER)

.PHONY: migration_down
migration_down:
	@.dev/migration_down.sh -d $(DRIVER)