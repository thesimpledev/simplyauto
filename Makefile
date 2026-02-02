# SimplyAuto Makefile

APP_NAME := SimplyAuto
APP_ID := com.simplyauto.app
BINARY_NAME := simplyauto.exe
BIN_DIR := bin
ICON := assets/logo.png

VERSION ?=
GITHUB_REPO := $(shell git remote get-url origin | sed 's/.*github.com[:/]\(.*\)\.git/\1/')

.PHONY: help setup debug build-debug build clean

help:
	@echo "SimplyAuto Build System"
	@echo ""
	@echo "Commands:"
	@echo "  make setup                - Install required tools (fyne, gh cli)"
	@echo "  make debug                - Quick build with console output for debugging"
	@echo "  make build-debug          - Production-like build for testing (no console)"
	@echo "  make build VERSION=x.x.x  - Build and release to GitHub"
	@echo "  make clean                - Remove build artifacts"
	@echo ""
	@echo "Example:"
	@echo "  make setup"
	@echo "  make debug"
	@echo "  make build-debug"
	@echo "  make build VERSION=1.0.0"
	@echo ""
	@echo "Stable download URL (always points to latest release):"
	@echo "  https://github.com/$(GITHUB_REPO)/releases/latest/download/$(BINARY_NAME)"

setup:
	@echo "Installing required tools..."
	go install fyne.io/tools/cmd/fyne@latest
	@echo ""
	@echo "Make sure you also have GitHub CLI (gh) installed and authenticated:"
	@echo "  https://cli.github.com/"
	@echo "  gh auth login"
	@echo ""
	@echo "Setup complete! Make sure ~/go/bin is in your PATH"

debug:
	@echo "=== Building $(APP_NAME) (debug) for Windows ==="
	@echo "Note: Console window will show for debugging output"
	mkdir -p $(BIN_DIR)
	cd cmd/simplyauto && \
		go build -o $(APP_NAME).exe .
	mv cmd/simplyauto/$(APP_NAME).exe $(BIN_DIR)/$(BINARY_NAME)
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)"

build-debug:
	@echo "=== Building $(APP_NAME) (production-like) for Windows ==="
	mkdir -p $(BIN_DIR)
	cd cmd/simplyauto && \
		fyne package --os windows --icon ../../$(ICON) --app-id $(APP_ID) --name $(APP_NAME)
	mv cmd/simplyauto/$(APP_NAME).exe $(BIN_DIR)/$(BINARY_NAME)
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)"

build:
ifndef VERSION
	$(error VERSION is required. Usage: make build VERSION=x.x.x)
endif
	@echo "=== Building $(APP_NAME) v$(VERSION) for Windows ==="
	mkdir -p $(BIN_DIR)
	cd cmd/simplyauto && \
		fyne package --os windows --icon ../../$(ICON) --app-id $(APP_ID) --name $(APP_NAME) --app-version $(VERSION) --release && \
		go build -ldflags "-H windowsgui -X main.Version=$(VERSION)" -o $(APP_NAME).exe .
	mv cmd/simplyauto/$(APP_NAME).exe $(BIN_DIR)/$(BINARY_NAME)
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)"
	@echo ""
	@echo "=== Creating GitHub Release v$(VERSION) ==="
	@if gh release view v$(VERSION) --repo $(GITHUB_REPO) >/dev/null 2>&1; then \
		echo "Release v$(VERSION) exists, updating..."; \
		gh release upload v$(VERSION) $(BIN_DIR)/$(BINARY_NAME) --repo $(GITHUB_REPO) --clobber; \
	else \
		echo "Creating release v$(VERSION)..."; \
		gh release create v$(VERSION) $(BIN_DIR)/$(BINARY_NAME) \
			--repo $(GITHUB_REPO) \
			--title "$(APP_NAME) v$(VERSION)" \
			--notes "## $(APP_NAME) v$(VERSION)"; \
	fi
	@echo ""
	@echo "=== Done! ==="
	@echo "Release: https://github.com/$(GITHUB_REPO)/releases/tag/v$(VERSION)"
	@echo "Latest:  https://github.com/$(GITHUB_REPO)/releases/latest/download/$(BINARY_NAME)"

clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)
	rm -f cmd/simplyauto/*.exe
	rm -f cmd/simplyauto/*.syso
	@echo "Done"
