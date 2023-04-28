DIR := test_install_dir/replays/
SEMVER := 0.0.3
APP_NAME := wows-fast-stats-$(SEMVER).exe

.PHONY: dev
dev:
	$(eval REV := $(shell git rev-parse --short HEAD))
	$(eval LD_FLAGS := "-X main.semver=$(SEMVER) -X main.revision=$(REV)")
	wails dev -ldflags $(LD_FLAGS)

.PHONY: build
build:
	$(eval REV := $(shell git rev-parse --short HEAD))
	$(eval LD_FLAGS := "-X main.semver=$(SEMVER) -X main.revision=$(REV)")
	wails build -ldflags $(LD_FLAGS) -platform windows/amd64 -o $(APP_NAME)

.PHONY: fmt
fmt:
	go fmt
	cd frontend/ && npx prettier --write **/*.{svelte,html,css}

.PHONY: put-temp-arema-info
put-temp-arema-info:
	$(eval TEMP_ARENA_INFO := $(shell ls test_install_dir/replays | fzf))
	cp $(DIR)$(TEMP_ARENA_INFO) $(DIR)tempArenaInfo.json
