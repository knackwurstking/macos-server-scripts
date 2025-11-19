package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"op-anime-dl/internal/anime"

	"github.com/gocolly/colly/v2"
	"github.com/lmittmann/tint"
)

var (
	a *anime.Anime
	c *Config
)

func main() {
	c = NewConfig()
	parseFlags()

	a = anime.NewAnime("https://onepiece-tube.com")

	for true {
		slog.Info("Starting anime list fetch")
		data, err := a.GetEpisodenStreams()
		if err != nil {
			slog.Error("Get anime list failed!", "err", err.Error())
		} else {
			slog.Info("Got anime list with entries", "count", len(data.Entries))
			iterAnimeList()
		}

		sleep()
	}
}

func iterAnimeList() {
	var (
		currentDownloadRequests = 0
	)

	for _, entry := range a.Data.Entries {
		if entry.Href == "" {
			slog.Debug("Skip entry (missing href attribute)", "entry.Number", entry.Number)
			continue
		}

		arc := a.Data.Arcs.Get(entry.ArcID)

		fileName := fmt.Sprintf("%04d %s (%s_SUB).mp4",
			entry.Number, entry.Name, strings.ToUpper(entry.LangSub))
		dirName := fmt.Sprintf("%03d %s", a.Data.Arcs.GetIndex(arc.ID)+1, arc.Name)

		slog.Debug("Generate file name", "dirName", dirName, "fileName", fileName, "entryNumber", entry.Number)

		mkdirAll(dirName)
		path := filepath.Join(c.Download.Dst, dirName, fileName)
		if _, err := os.Stat(path); err == nil {
			continue
		}

		currentDownloadRequests += 1
		downloadEntry(path, entry)

		duration := time.Minute * time.Duration(c.Download.Delay)
		if currentDownloadRequests >= c.Download.LimitPerDay {
			duration = time.Minute * time.Duration(c.Download.LongDelay)
			currentDownloadRequests = 0
		}

		slog.Info("Sleeping before next download", "duration", duration,
			"current_download_requests", currentDownloadRequests)
		time.Sleep(duration)
	}
}

func mkdirAll(dirName string) {
	path := filepath.Join(c.Download.Dst, dirName)
	_, err := os.Stat(path)
	if err != nil {
		slog.Debug("Create directories", "path", path, "dirName", dirName)
		err = os.MkdirAll(path, os.ModeDir|os.ModePerm)
		if err != nil {
			slog.Error("Failed to create directory", "path", path, "err", err.Error())
			// Don't panic - just log and continue execution
			return
		}
	}
}

func downloadEntry(path string, entry anime.AnimeDataEntry) {
	// Create a new collector for each download to avoid resource leaks
	collector := colly.NewCollector()

	// Set up collector for download process
	collector.OnHTML("iframe", func(h *colly.HTMLElement) {
		src := h.Attr("src")
		if src == "" {
			return
		}

		// Create nested collector for iframe content
		iframeCollector := colly.NewCollector()

		iframeCollector.OnHTML("video > source", func(h *colly.HTMLElement) {
			src := h.Attr("src")
			if src == "" {
				return
			}

			if h.Attr("type") != "video/mp4" {
				slog.Warn("HTML tag <source has not type \"video/mp4\" attribute")
				return
			}

			slog.Debug("Got url from video source", "src", src, "entryName", entry.Name)

			// Use the direct download method from anime package
			err := a.Download(entry, path)
			if err != nil {
				slog.Error("download src to dst failed", "err", err, "src", src, "dst", path)
				_ = os.Remove(path)
			}
		})

		iframeCollector.OnRequest(func(r *colly.Request) {
			slog.Debug(fmt.Sprintf("Request to \"%s\"", r.URL))
		})

		iframeCollector.OnError(func(r *colly.Response, err error) {
			slog.Error("Colly error", "err", err.Error())
		})

		if err := iframeCollector.Visit(src); err != nil {
			slog.Error(fmt.Sprintf("Visit \"%s\" failed!", src), "err", err.Error())
		}

		iframeCollector.Wait()
	})

	collector.OnRequest(func(r *colly.Request) {
		slog.Debug(fmt.Sprintf("Request to \"%s\"", r.URL))
	})

	collector.OnError(func(r *colly.Response, e error) {
		slog.Error("Main collector error", "err", e.Error())
	})

	if err := collector.Visit(entry.Href); err != nil {
		slog.Error("Download visit failed", "err", err.Error())
	}

	collector.Wait()
}

func sleep() {
	now := time.Now()

	// Calculate next update time properly
	nextUpdate := time.Date(now.Year(), now.Month(), now.Day(), c.Update.Hour, 0, 0, 0, time.Local)

	// If it's already past the update hour today, schedule for tomorrow
	if now.Hour() >= c.Update.Hour {
		nextUpdate = nextUpdate.AddDate(0, 0, 1)
	}

	duration := nextUpdate.Sub(now)
	slog.Info("Sleeping until next update day", "duration", duration, "next_update", nextUpdate, "current_time", now)

	time.Sleep(duration)

	if time.Now().Weekday() == c.Update.Weekday {
		slog.Info("Running new update now...")
	}
}

func parseFlags() {
	flag.BoolVar(&c.Debug, "debug", c.Debug, "Enable debugging")

	flag.IntVar(
		&c.Download.Delay,
		"delay",
		c.Download.Delay,
		"Set delay in minutes between downloads",
	)

	flag.IntVar(
		&c.Download.LongDelay,
		"long-delay",
		c.Download.LongDelay,
		"Set long delay in minutes if download limit was reached",
	)

	flag.StringVar(
		&c.Download.Dst,
		"dst",
		c.Download.Dst,
		"Set destination path for downloads",
	)

	flag.IntVar(
		&c.Download.LimitPerDay,
		"limit",
		c.Download.LimitPerDay,
		"Download limit (per day)",
	)

	weekday := flag.Int("update-on-day", int(c.Update.Weekday),
		"Weekday (0-6) for update the anime list")

	hour := flag.Int("update-hour", c.Update.Hour,
		"Hour (0-23) for anime list update")

	flag.Parse()

	if *weekday >= 0 && *weekday <= 6 {
		c.Update.Weekday = time.Weekday(*weekday)
	}

	if *hour >= 0 && *hour <= 23 {
		c.Update.Hour = *hour
	}

	options := &tint.Options{
		TimeFormat: time.DateTime,
		Level:      slog.LevelInfo,
	}

	if c.Debug {
		options.Level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, options)))
}
