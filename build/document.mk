#
# Makefile fragment for auto-documenting Golang projects
#

GODOC        ?= godocdown

document:
	@echo "=== $(PROJECT_NAME) === [ documentation    ]: Generating Godoc in Markdown..."
	@for p in $(PACKAGES); do \
		echo "=== $(PROJECT_NAME) === [ documentation    ]:     $$p"; \
		$(GODOC) $$p > $$p/README.md ; \
	done
	@for c in $(COMMANDS); do \
		echo "=== $(PROJECT_NAME) === [ documentation    ]:     $$c"; \
		$(GODOC) $$c > $$c/README.md ; \
	done

