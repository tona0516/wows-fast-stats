APP := wows-fast-stats
SEMVER := 0.7.0
REVISION := $(shell git rev-parse --short HEAD)
EXE := $(APP)-$(SEMVER).exe
ZIP := $(APP).zip
DISCORD_WEBHOOK_URL := $(shell cat discord_webhook_url)

LDFLAGS_COMMON := -X main.semver=$(SEMVER) -X main.revision=$(REVISION) -X main.discordWebhookURL=$(DISCORD_WEBHOOK_URL)
LDFLAGS_DEV := -X main.env=debug $(LDFLAGS_COMMON)
LDFLAGS_PROD := -X main.env=production $(LDFLAGS_COMMON)

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
	rm -rf $(APP) $(ZIP)
	mkdir $(APP)
	cp build/bin/$(EXE) $(APP)
	zip -r $(ZIP) $(APP)
	rm -rf $(APP)

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
