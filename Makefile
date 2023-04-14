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
