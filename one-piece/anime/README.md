# One Piece Anime Downloader

This project is a Go-based application designed to download One Piece anime from "onepiece-tube.com".

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
- `-limit`: Number of anime to download per run (default: 5)
- `-delay`: Delay between downloads in minutes (default: 30)
- `-long-delay`: Delay between anime series in minutes (default: 720)
- `-dst`: Destination directory for downloaded anime (defaults: $HOME/OnePieceAnime)
- `-update-on-day`: Weekday (0-6) for update the anime list (default: 0)
- `-update-hour`: Hour (0-23) for anime list update (default: 18)

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
