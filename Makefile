DIR := test_install_dir/replays/

.PHONY: dev
dev:
	wails dev

.PHONY: build
build:
	wails build -platform windows/amd64

.PHONY: fmt
fmt:
	go fmt
	cd frontend/ && npx prettier --write **/*.{svelte,html,css}

.PHONY: put-temp-arema-info
put-temp-arema-info:
	$(eval TEMP_ARENA_INFO := $(shell ls test_install_dir/replays | fzf))
	cp $(DIR)$(TEMP_ARENA_INFO) $(DIR)tempArenaInfo.json
