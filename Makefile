DIR := test_install_dir/replays/
SEMVER := 0.6.1
APP := wows-fast-stats
EXE := $(APP)-$(SEMVER).exe
ZIP := $(APP).zip
DISCORD_WEBHOOK_URL := $(shell cat discord_webhook_url)

.PHONY: dev
dev:
	$(eval REV := $(shell git rev-parse --short HEAD))
	$(eval LD_FLAGS := "-X main.semver=$(SEMVER) -X main.revision=$(REV) -X main.env=debug -X main.discordWebhookURL=$(DISCORD_WEBHOOK_URL)")
	wails dev -ldflags $(LD_FLAGS)

.PHONY: check-prerequisite
check-prerequisite:
	command -v go > /dev/null 2>&1
	command -v npm > /dev/null 2>&1

.PHONY: setup
setup: check-prerequisite
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: build
build: lint test
	$(eval REV := $(shell git rev-parse --short HEAD))
	$(eval LD_FLAGS := "-X main.semver=$(SEMVER) -X main.revision=$(REV) -X main.env=production -X main.discordWebhookURL=$(DISCORD_WEBHOOK_URL)")
	wails build -ldflags $(LD_FLAGS) -platform windows/amd64 -o $(EXE) -trimpath

.PHONY: build-nolint
build-nolint:
	$(eval REV := $(shell git rev-parse --short HEAD))
	$(eval LD_FLAGS := "-X main.semver=$(SEMVER) -X main.revision=$(REV) -X main.env=production -X main.discordWebhookURL=$(DISCORD_WEBHOOK_URL)")
	wails build -ldflags $(LD_FLAGS) -platform windows/amd64 -o $(EXE) -trimpath

.PHONY: package
package: build
	rm -rf $(APP) $(ZIP)
	mkdir $(APP)
	cp build/bin/$(EXE) $(APP)
	zip -r $(ZIP) $(APP)
	rm -rf $(APP)

.PHONY: lint
lint:
	golangci-lint run
	cd frontend/ && npm ci && npm run check

.PHONY: fmt
fmt:
	golangci-lint run --fix
	go fmt
	cd frontend/ && npx prettier --write src/**/*.{ts,svelte,css} index.html

.PHONY: test
test:
	go test -cover ./...

.PHONY: put-temp-arema-info
put-temp-arema-info:
	$(eval TEMP_ARENA_INFO := $(shell ls test_install_dir/replays | fzf))
	cp $(DIR)$(TEMP_ARENA_INFO) $(DIR)tempArenaInfo.json
