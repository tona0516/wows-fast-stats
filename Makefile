APP_NAME := wows-fast-stats
SEMVER := 0.16.2-alpha1
WG_APP_ID := e25e1a2af190880c9e33d3be7cc5313d
BINARY_NAME := $(APP_NAME)-$(SEMVER).exe
NPM := npm --prefix ./frontend
BUILD_FLAGS := -trimpath -webview2 embed
DISCORD_WEBHOOK_URL_PROD_ALERT := $(shell cat discord_webhook_url_prod_alert)
DISCORD_WEBHOOK_URL_PROD_INFO := $(shell cat discord_webhook_url_prod_info)
DISCORD_WEBHOOK_URL_DEV_ALERT := $(shell cat discord_webhook_url_dev_alert)
DISCORD_WEBHOOK_URL_DEV_INFO := $(shell cat discord_webhook_url_dev_info)
COMMON_LDFLAGS := -X main.AppName=$(APP_NAME) -X main.Semver=$(SEMVER) -X main.WGAppID=$(WG_APP_ID)
TEST_REPLAY_PATH := ./test_install_dir/replays

.PHONY: setup gen-mock lint fmt test dev build pkg clean uddep chbtl

setup:
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3
	go install go.uber.org/mock/mockgen@latest
	go install github.com/arch-go/arch-go@latest
	$(NPM) ci

gen-mock:
	go generate ./...

lint:
	golangci-lint run
	$(NPM) run check
	$(NPM) run lint

fmt:
	golangci-lint run --fix
	go fmt
	$(NPM) run fmt

test:
	arch-go
	go test -cover ./... -count=1
	$(NPM) run test

dev:
	wails dev -ldflags "$(COMMON_LDFLAGS) -X main.AlertDiscordWebhookURL=$(DISCORD_WEBHOOK_URL_DEV_ALERT) -X main.InfoDiscordWebhookURL=$(DISCORD_WEBHOOK_URL_DEV_INFO) -X main.IsDev=true"

build: test
	wails build -platform windows/amd64 -o $(BINARY_NAME) $(BUILD_FLAGS) -ldflags "$(COMMON_LDFLAGS) -X main.AlertDiscordWebhookURL=$(DISCORD_WEBHOOK_URL_PROD_ALERT) -X main.InfoDiscordWebhookURL=$(DISCORD_WEBHOOK_URL_PROD_INFO)"

pkg: build
	$(eval TEMP_DIR := $(shell mktemp -d))
	mkdir -p $(TEMP_DIR)/$(APP_NAME)
	mv ./build/bin/$(BINARY_NAME) $(TEMP_DIR)/$(APP_NAME)
	zip -r $(APP_NAME).zip $(TEMP_DIR)/$(APP_NAME)

chbtl:
	$(eval TEMP_ARENA_INFO := $(shell ls $(TEST_REPLAY_PATH) | fzf))
	cp $(TEST_REPLAY_PATH)/$(TEMP_ARENA_INFO) $(TEST_REPLAY_PATH)/tempArenaInfo.json

clean:
	rm -rf config/ cache/ temp_arena_info/ screenshot/ persistent_data/

uddep:
	go get -u
	$(NPM) update
