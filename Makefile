DIR := test_install_dir/replays/
SEMVER := 0.2.0
APP := wows-fast-stats
EXE := $(APP)-$(SEMVER).exe
ZIP := $(APP).zip

.PHONY: dev
dev:
	$(eval REV := $(shell git rev-parse --short HEAD))
	$(eval LD_FLAGS := "-X main.semver=$(SEMVER) -X main.revision=$(REV)")
	wails dev -ldflags $(LD_FLAGS)

.PHONY: setup
setup:
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

.PHONY: build
build:
	$(eval REV := $(shell git rev-parse --short HEAD))
	$(eval LD_FLAGS := "-X main.semver=$(SEMVER) -X main.revision=$(REV)")
	wails build -ldflags $(LD_FLAGS) -platform windows/amd64 -o $(EXE)

.PHONY: package
package: build
	rm -rf $(APP) $(ZIP)
	mkdir $(APP)
	cp build/bin/$(EXE) $(APP)
	zip -r $(ZIP) $(APP)
	rm -rf $(APP)

.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: fmt
fmt:
	go fmt
	cd frontend/ && npx prettier --write **/*.{ts,svelte,html,css}

.PHONY: put-temp-arema-info
put-temp-arema-info:
	$(eval TEMP_ARENA_INFO := $(shell ls test_install_dir/replays | fzf))
	cp $(DIR)$(TEMP_ARENA_INFO) $(DIR)tempArenaInfo.json
