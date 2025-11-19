BINARY_NAME := op-manga-dl

# Directories and paths
BIN_DIR     := ./bin
INSTALL_DIR := /usr/local/bin
SERVICE_FILE:= $(HOME)/Library/LaunchAgents/com.$(BINARY_NAME).plist
LOG_FILE    := $(HOME)/Library/Application Support/$(BINARY_NAME)/$(BINARY_NAME).log

LIMIT := 2
DELAY := 60
LONG_DELAY := 180
DST := /Volumes/NAS/Videos/One Piece/Manga

.PHONY: all init build clean macos-install macos-start-service macos-stop-service \
        macos-restart-service macos-print-service macos-watch-service

all: init build

init:
	go mod tidy -v

build:
	go build -v -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/op-manga-dl

clean:
	rm -rf $(BIN_DIR)

# Create launchd service file for macOS
define LAUNCHD_SERVICE_FILE_CONTENT
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
"http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.$(BINARY_NAME)</string>

    <key>ProgramArguments</key>
    <array>
        <string>$(INSTALL_DIR)/$(BINARY_NAME)</string>
        <string>-debug</string>
        <string>-limit</string>
        <string>$(LIMIT)</string>
        <string>-delay</string>
        <string>$(DELAY)</string>
        <string>-long-delay</string>
        <string>$(LONG_DELAY)</string>
        <string>-dst</string>
        <string>$(DST)</string>
    </array>

    <key>RunAtLoad</key>
    <true/>

    <key>KeepAlive</key>
    <false/>

    <key>StandardOutPath</key>
    <string>$(LOG_FILE)</string>

    <key>StandardErrorPath</key>
    <string>$(LOG_FILE)</string>

    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
</dict>
</plist>
endef

export LAUNCHD_SERVICE_FILE_CONTENT

macos-install:
	@echo "Installing $(BINARY_NAME) for macOS..."
	mkdir -p $(INSTALL_DIR)
	sudo cp $(BIN_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	echo "$$LAUNCHD_SERVICE_FILE_CONTENT" > $(SERVICE_FILE)
	@echo "$(BINARY_NAME) installed successfully"

macos-start-service:
	@echo "Starting $(BINARY_NAME) service..."
	launchctl load -w $(SERVICE_FILE)
	launchctl start com.$(BINARY_NAME)

macos-stop-service:
	@echo "Stopping $(BINARY_NAME) service..."
	launchctl stop com.$(BINARY_NAME) || exit 0
	launchctl unload -w $(SERVICE_FILE)

macos-restart-service:
	@echo "Restarting $(BINARY_NAME) service..."
	make macos-stop-service || exit 0
	make macos-start-service

macos-print-service:
	@echo "$(BINARY_NAME) service information:"
	launchctl print gui/$(shell id -u)/com.$(BINARY_NAME) || echo "Service not loaded or running"

macos-watch-service:
	@echo "$(BINARY_NAME) watch server logs @ \"$(LOG_FILE)\":"
	if [ -f "$(LOG_FILE)" ]; then \
		echo "Watching logs... Press Ctrl+C to stop"; \
		tail -f "$(LOG_FILE)"; \
	else \
		echo "Log file not found. Make sure the service is running or has been started."; \
		echo "Log file path: \"$(LOG_FILE)\""; \
	fi
