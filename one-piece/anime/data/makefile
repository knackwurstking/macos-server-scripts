# TODO: Install scripts for macos
BINARY_NAME := op-anime-dl

EXTENSION_TAGS := stats,ip,journal,opmanga,pgpress

# Directories and paths
BIN_DIR := ./bin
INSTALL_DIR := /usr/local/bin
SERVICE_FILE := $(HOME)/Library/LaunchAgents/com.$(BINARY_NAME).plist
LOG_FILE := $(HOME)/Library/Application Support/$(BINARY_NAME).log

.PHONY: macos-install macos-start-service macos-stop-service \
	macos-restart-service macos-print-service macos-watch-service

# Create launchd service file for macOS
define LAUNCHD_SERVICE_FILE_CONTENT
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>$(BINARY_NAME)</string>

    <key>ProgramArguments</key>
    <array>
        <string>$(BIN_DIR)/$(BINARY_NAME)</string>
        <string>-debug</string>
        <string>-limit</string>
        <string>2</string>
        <string>-delay</string>
        <string>60</string>
        <string>-long-delay</string>
        <string>180</string>
    </array>

    <key>WorkingDirectory</key>
    <string>$(BIN_DIR)</string>

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
        <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
</dict>
</plist>
endef

export LAUNCHD_SERVICE_FILE_CONTENT

macos-install:
	@echo "Installing $(BINARY_NAME) for macOS..."
	@mkdir -p $(INSTALL_DIR)
	@sudo cp $(BIN_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$$LAUNCHD_SERVICE_FILE_CONTENT" > $(SERVICE_FILE)
	@echo "$(BINARY_NAME) installed successfully"

macos-start-service:
	@echo "Starting $(BINARY_NAME) service..."
	@launchctl load -w $(SERVICE_FILE)
	@launchctl start com.$(BINARY_NAME)

macos-stop-service:
	@echo "Stopping $(BINARY_NAME) service..."
	@launchctl stop com.$(BINARY_NAME)
	@launchctl unload -w $(SERVICE_FILE)

macos-restart-service:
	@echo "Restarting $(BINARY_NAME) service..."
	make macos-stop-service
	make macos-start-service

macos-print-service:
	@echo "$(BINARY_NAME) service information:"
	@launchctl print gui/$$(id -u)/com.$(BINARY_NAME) || echo "Service not loaded or running"

macos-watch-service:
	@echo "$(BINARY_NAME) watch server logs @ \"$(LOG_FILE)\":"
	@if [ -f "$(LOG_FILE)" ]; then \
		echo "Watching logs... Press Ctrl+C to stop"; \
		tail -f "$(LOG_FILE)"; \
	else \
		echo "Log file not found. Make sure the service is running or has been started."; \
		echo "Log file path: \"$(LOG_FILE)\""; \
	fi
