# SimplyAuto Makefile

APP_NAME := SimplyAuto
APP_ID := com.simplyauto.app
BINARY_NAME := simplyauto.exe
BIN_DIR := bin
ICON := assets/logo.png

VERSION ?=
GITHUB_REPO := $(shell git remote get-url origin | sed 's/.*github.com[:/]\(.*\)\.git/\1/')

.PHONY: help setup build clean

# Default target
help:
	@echo "SimplyAuto Build System"
	@echo ""
	@echo "Commands:"
	@echo "  make setup                - Install required tools (fyne, gh cli)"
	@echo "  make build VERSION=x.x.x  - Build and release to GitHub"
	@echo "  make clean                - Remove build artifacts"
	@echo ""
	@echo "Example:"
	@echo "  make setup"
	@echo "  make build VERSION=1.0.0"

# Install required tools
setup:
	@echo "Installing required tools..."
	go install fyne.io/tools/cmd/fyne@latest
	@echo ""
	@echo "Make sure you also have GitHub CLI (gh) installed and authenticated:"
	@echo "  https://cli.github.com/"
	@echo "  gh auth login"
	@echo ""
	@echo "Setup complete! Make sure ~/go/bin is in your PATH"

# Build and release
build:
ifndef VERSION
	$(error VERSION is required. Usage: make build VERSION=x.x.x)
endif
	@echo "=== Building $(APP_NAME) v$(VERSION) for Windows ==="
	mkdir -p $(BIN_DIR)
	cd cmd/simplyauto && \
		fyne package --os windows --icon ../../$(ICON) --app-id $(APP_ID) --name $(APP_NAME) --app-version $(VERSION) --release
	mv cmd/simplyauto/$(APP_NAME).exe $(BIN_DIR)/$(BINARY_NAME)
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)"
	@echo ""
	@echo "=== Creating GitHub Release v$(VERSION) ==="
	@if gh release view v$(VERSION) >/dev/null 2>&1; then \
		echo "Release v$(VERSION) exists, updating..."; \
		gh release upload v$(VERSION) $(BIN_DIR)/$(BINARY_NAME) --clobber; \
	else \
		echo "Creating new release v$(VERSION)..."; \
		gh release create v$(VERSION) $(BIN_DIR)/$(BINARY_NAME) \
			--title "$(APP_NAME) v$(VERSION)" \
			--notes "## $(APP_NAME) v$(VERSION)"; \
	fi
	@echo ""
	@echo "=== Updating 'latest' Release ==="
	@if gh release view latest >/dev/null 2>&1; then \
		echo "Updating latest release..."; \
		gh release delete latest -y; \
		git push origin :refs/tags/latest 2>/dev/null || true; \
	fi
	gh release create latest $(BIN_DIR)/$(BINARY_NAME) \
		--title "$(APP_NAME) Latest (v$(VERSION))" \
		--notes "Latest stable release. Current version: v$(VERSION)"
	@echo ""
	@echo "=== Done! ==="
	@echo "Versioned: https://github.com/$(GITHUB_REPO)/releases/tag/v$(VERSION)"
	@echo "Latest:    https://github.com/$(GITHUB_REPO)/releases/download/latest/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)
	rm -f cmd/simplyauto/*.exe
	rm -f cmd/simplyauto/*.syso
	@echo "Done"
