# SimplyAuto Makefile

APP_NAME := SimplyAuto
BINARY_NAME := simplyauto.exe
BIN_DIR := bin
ICON := assets/logo.png

VERSION ?=

.PHONY: help setup build clean

# Default target
help:
	@echo "SimplyAuto Build System"
	@echo ""
	@echo "Commands:"
	@echo "  make setup              - Install required tools (fyne, etc.)"
	@echo "  make build VERSION=x.x.x - Build Windows executable with embedded icon"
	@echo "  make clean              - Remove build artifacts"
	@echo ""
	@echo "Example:"
	@echo "  make setup"
	@echo "  make build VERSION=1.0.0"

# Install required tools
setup:
	@echo "Installing required tools..."
	go install fyne.io/fyne/v2/cmd/fyne@latest
	@echo ""
	@echo "Setup complete! Make sure ~/go/bin is in your PATH"

# Build the Windows executable with embedded icon
build:
ifndef VERSION
	$(error VERSION is required. Usage: make build VERSION=x.x.x)
endif
	@echo "Building $(APP_NAME) v$(VERSION) for Windows..."
	mkdir -p $(BIN_DIR)
	cd cmd/simplyauto && \
		fyne package -os windows -icon ../../$(ICON) -name $(APP_NAME) -appVersion $(VERSION) -release
	mv cmd/simplyauto/$(APP_NAME).exe $(BIN_DIR)/$(BINARY_NAME)
	@echo ""
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)
	rm -f cmd/simplyauto/*.exe
	rm -f cmd/simplyauto/*.syso
	@echo "Done"
