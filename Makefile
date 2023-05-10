DIR := test_install_dir/replays/
SEMVER := 0.3.0
APP := wows-fast-stats
EXE := $(APP)-$(SEMVER).exe
ZIP := $(APP).zip

.PHONY: dev
dev:
	$(eval REV := $(shell git rev-parse --short HEAD))
	$(eval LD_FLAGS := "-X main.semver=$(SEMVER) -X main.revision=$(REV) -X main.env=debug")
	wails dev -ldflags $(LD_FLAGS)

.PHONY: setup
setup:
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

.PHONY: build
build: lint
	$(eval REV := $(shell git rev-parse --short HEAD))
	$(eval LD_FLAGS := "-X main.semver=$(SEMVER) -X main.revision=$(REV) -X main.env=production")
	wails build -ldflags $(LD_FLAGS) -platform windows/amd64 -o $(EXE) -trimpath

.PHONY: build-nolint
build-nolint:
	$(eval REV := $(shell git rev-parse --short HEAD))
	$(eval LD_FLAGS := "-X main.semver=$(SEMVER) -X main.revision=$(REV) -X main.env=production")
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
	cd frontend/ && npx prettier --write **/*.{ts,svelte,html,css}

.PHONY: put-temp-arema-info
put-temp-arema-info:
	$(eval TEMP_ARENA_INFO := $(shell ls test_install_dir/replays | fzf))
	cp $(DIR)$(TEMP_ARENA_INFO) $(DIR)tempArenaInfo.json
