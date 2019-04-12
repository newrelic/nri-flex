#
# Makefile fragment for Packaging
#

package-all: package-linux package-darwin package-windows

package: compile-only
	@echo "=== $(PROJECT_NAME) === [ package          ]: Creating package for $(GOOS)..."
	@for b in $(BINS); do \
		echo "=== $(PROJECT_NAME) === [ package          ]:     $$b"; \
		./build/package/create-tar.sh $$b $(GOOS) $(PROJECT_VER) ; \
	done

package-linux: compile-linux
	@echo "=== $(PROJECT_NAME) === [ package-linux    ]: Creating package for Linux..."
	@for b in $(BINS); do \
		echo "=== $(PROJECT_NAME) === [ package-linux    ]:     $$b"; \
		./build/package/create-tar.sh $$b linux $(PROJECT_VER) ; \
	done

package-darwin: compile-darwin
	@echo "=== $(PROJECT_NAME) === [ package-darwin   ]: Creating package for Darwin..."
	@for b in $(BINS); do \
		echo "=== $(PROJECT_NAME) === [ package-darwin   ]:     $$b"; \
		./build/package/create-tar.sh $$b darwin $(PROJECT_VER) ; \
	done

package-windows: compile-windows
	@echo "=== $(PROJECT_NAME) === [ package-windows  ]: Creating package for Windows..."
	@for b in $(BINS); do \
		echo "=== $(PROJECT_NAME) === [ package-windows  ]:     $$b"; \
		./build/package/create-tar.sh $$b windows $(PROJECT_VER) ; \
	done

