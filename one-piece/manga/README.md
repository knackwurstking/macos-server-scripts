# One Piece Manga Downloader

This project is a Go-based application designed to download One Piece manga from "onepiece-tube.com".

## Building

To build the project, run:
```bash
make
```

## Installation

For macOS installation:
```bash
make macos-install
```

This will install the binary to `/usr/local/bin` and set up a launchd service.

## Configuration

The application supports various flags:
- `-debug`: Enable debug logging
- `-limit`: Number of manga to download per run (default: 2)
- `-delay`: Delay between downloads in seconds (default: 60)
- `-long-delay`: Delay between manga series in seconds (default: 180)
- `-dst`: Destination directory for downloaded manga

## Service Management

The project includes macOS service management commands:
- `make macos-start-service` - Start the service
- `make macos-stop-service` - Stop the service  
- `make macos-restart-service` - Restart the service
- `make macos-watch-service` - Watch service logs

## Dependencies

The project uses:
- [gocolly/colly/v2](https://github.com/gocolly/colly) - Web scraping library
- [lmittmann/tint](https://github.com/lmittmann/tint) - Logging library
