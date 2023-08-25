APP_NAME := wows-fast-stats
SEMVER := 0.9.0
EXE := $(APP_NAME)-$(SEMVER).exe
ZIP := $(APP_NAME).zip
DISCORD_WEBHOOK_URL_PROD := $(shell cat discord_webhook_url_prod)
DISCORD_WEBHOOK_URL_DEV := $(shell cat discord_webhook_url_dev)

LDFLAGS_COMMON := -X main.AppName=$(APP_NAME) -X main.Semver=$(SEMVER)
LDFLAGS_DEV := $(LDFLAGS_COMMON) -X main.DiscordWebhookURL=$(DISCORD_WEBHOOK_URL_DEV) -X main.IsDev=true
LDFLAGS_PROD := $(LDFLAGS_COMMON) -X main.DiscordWebhookURL=$(DISCORD_WEBHOOK_URL_PROD)

TEST_DIR := test_install_dir/replays/

.PHONY: dev
dev: gen
	wails dev -ldflags "$(LDFLAGS_DEV)"

.PHONY: check-prerequisite
check-prerequisite:
	command -v go > /dev/null 2>&1
	command -v npm > /dev/null 2>&1

.PHONY: setup
setup: check-prerequisite
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/stringer@latest
	npm ci

.PHONY: build
build: gen lint test
	wails build -ldflags "$(LDFLAGS_PROD)" -platform windows/amd64 -o $(EXE) -trimpath

.PHONY: build-nolint
build-nolint: gen
	wails build -ldflags "$(LDFLAGS_PROD)" -platform windows/amd64 -o $(EXE) -trimpath

.PHONY: package
package: build
	rm -rf $(APP_NAME) $(ZIP)
	mkdir $(APP_NAME)
	cp build/bin/$(EXE) $(APP_NAME)
	zip -r $(ZIP) $(APP_NAME)
	rm -rf $(APP_NAME)

.PHONY: gen
gen:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run
	cd frontend/ && npm run check

.PHONY: fmt
fmt:
	golangci-lint run --fix
	go fmt
	cd frontend/ && npx prettier --write "src/**/*.{ts,svelte,css}" index.html

.PHONY: test
test:
	go test -cover ./...

.PHONY: put-temp-arema-info
put-temp-arema-info:
	$(eval TEMP_ARENA_INFO := $(shell ls test_install_dir/replays | fzf))
	cp $(TEST_DIR)$(TEMP_ARENA_INFO) $(TEST_DIR)tempArenaInfo.json

.PHONY: remove-files
remove-files:
	rm -rf config/ cache/ temp_arena_info/ screenshot/
