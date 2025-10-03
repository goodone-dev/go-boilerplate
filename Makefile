.PHONY: migration-create
migration-create: #Create new migration file
	@.dev/create-migration.sh -n $(NAME) -d $(DRIVER)

.PHONY: migration-apply
migration-apply: #Run migration up
	@.dev/apply-migration.sh -d $(DRIVER)

.PHONY: migration-rollback
migration-rollback: #Run migration down
	@.dev/rollback-migration.sh -d $(DRIVER)