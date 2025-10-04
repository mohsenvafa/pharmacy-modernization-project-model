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

.PHONY: setup tailwind-watch dev dev-watch mock-iris check-tools build-ts watch-ts

setup:
	@make -f .dev/Makefile.setup setup

check-tools:
	@make -f .dev/Makefile.setup check-tools

tailwind-watch:
	@$(MAKE) check-tools
	@set -euo pipefail; \
	"$(TAILWIND_BIN)" -c web/tailwind.config.js -i web/styles/input.css -o web/public/app.css --watch

# Run templ in watch mode (Tailwind watcher can run separately).
dev-watch:
	@set -euo pipefail; \
	templ generate -watch \
		-proxyport=7332 \
		-proxy="http://localhost:8080" \
		-cmd="go run ./cmd/server" \
		-open-browser=false

# Convenience target to run all watchers together.
dev:
	@set -euo pipefail; \
	$(MAKE) tailwind-watch & \
	TAILWIND_PID=$$!; \
	$(MAKE) watch-ts & \
	TS_PID=$$!; \
	trap 'kill $$TAILWIND_PID >/dev/null 2>&1 || true; kill $$TS_PID >/dev/null 2>&1 || true' EXIT INT TERM; \
	$(MAKE) dev-watch

mock-iris:
	go run ./cmd/iris_mock

# Build TypeScript
build-ts:
	@cd web && npm run build

# Watch TypeScript for changes
watch-ts:
	@cd web && npm run watch
