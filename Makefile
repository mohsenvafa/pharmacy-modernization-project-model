UNAME_S := $(shell uname -s 2>/dev/null)

ifeq ($(OS),Windows_NT)
  DEFAULT_TAILWIND_BIN := ./bin/tailwindcss.exe
else ifneq (,$(findstring MINGW,$(UNAME_S)))
  DEFAULT_TAILWIND_BIN := ./bin/tailwindcss.exe
else ifneq (,$(findstring MSYS,$(UNAME_S)))
  DEFAULT_TAILWIND_BIN := ./bin/tailwindcss.exe
else ifneq (,$(findstring CYGWIN,$(UNAME_S)))
  DEFAULT_TAILWIND_BIN := ./bin/tailwindcss.exe
else
  DEFAULT_TAILWIND_BIN := ./bin/tailwindcss
endif

TAILWIND_BIN ?= $(DEFAULT_TAILWIND_BIN)

.PHONY: setup tailwind-watch dev dev-watch mock-iris check-tools build-ts watch-ts graphql-generate graphql-install

setup:
	@make -f .dev/Makefile.setup setup

check-tools:
	@make -f .dev/Makefile.setup check-tools

tailwind-watch:
	@$(MAKE) check-tools
	@set -euo pipefail; \
	cd web && npx tailwindcss -c tailwind.config.js -i styles/input.css -o public/app.css --watch

# Run templ in watch mode (Tailwind watcher can run separately).
dev-watch:
	@set -euo pipefail; \
	templ generate -watch \
		-proxyport=7332 \
		-proxy="http://localhost:8080" \
		-cmd="go run -gcflags=all=-N -gcflags=all=-l ./cmd/server" \
		-open-browser=false

# Convenience target to run all watchers together.
dev:
	@set -euo pipefail; \
	$(MAKE) tailwind-watch & \
	TAILWIND_PID=$$!; \
	$(MAKE) watch-ts & \
	TS_PID=$$!; \
	trap 'kill $$TAILWIND_PID >/dev/null 2>&1 || true; kill $$TS_PID >/dev/null 2>&1 || true' EXIT INT TERM; \
	echo "🚀 Starting development server..."; \
	$(MAKE) dev-watch & \
	DEV_PID=$$!; \
	sleep 2; \
	echo ""; \
	echo "📋 Debug Information:"; \
	echo "   VS Code Debug Configuration: Use process ID below"; \
	echo "   Current Go server PID: $$(ps aux | grep 'go run.*cmd/server' | grep -v grep | awk '{print $$2}' | head -1)"; \
	echo ""; \
	wait $$DEV_PID

mock-iris:
	go run ./cmd/iris_mock

# Build TypeScript
build-ts:
	@cd web && npm run build

# Watch TypeScript for changes
watch-ts:
	@cd web && npm run watch

# Generate GraphQL code from schemas
graphql-generate:
	@echo "🔄 Generating GraphQL code..."
	@gqlgen generate
	@echo "✅ GraphQL code generated successfully!"

# Install gqlgen CLI tool
graphql-install:
	@echo "📦 Installing gqlgen..."
	@go install github.com/99designs/gqlgen@latest
	@echo "✅ gqlgen installed successfully!"
