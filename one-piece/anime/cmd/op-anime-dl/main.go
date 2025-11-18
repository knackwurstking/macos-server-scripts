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
		slog.Debug("Get anime list.", "url", a.GetUrl(anime.PathEpisodenStreams))
        _, err := a.GetEpisodenStreams()
		if err != nil {
			slog.Error("Get anime list failed!", "err", err.Error())
		} else {
			iterAnimeList()
		}

        sleep()
	}
}

func iterAnimeList() {
    var (
        currentDownloads = 0
    )

	for _, entry := range a.Data.Entries {
        if entry.Href == "" {
            slog.Debug("Skip entry (missing href attribute)", "entry.Number", entry.Number)
            continue
        }

        arc := a.Data.Arcs.Get(entry.ArcID)

		fileName := fmt.Sprintf("%04d %s (%s_SUB).mp4",
			entry.Number, entry.Name, strings.ToUpper(entry.LangSub))
        dirName := fmt.Sprintf("%03d %s", a.Data.Arcs.GetIndex(arc.ID) + 1, arc.Name)

		//slog.Debug("Generate file name", "dirName", dirName, "fileName", fileName)

        path := filepath.Join(c.Download.Dst, dirName, fileName)
        _, err := os.Stat(path)
        if err != nil {
            mkdirAll(dirName)
        } else {
            continue
        }

        currentDownloads += 1
        downloadEntry(path, entry)

        duration := time.Minute * time.Duration(c.Download.Delay)
        if currentDownloads >= c.Download.LimitPerDay {
            duration = time.Minute * time.Duration(c.Download.LongDelay)
            currentDownloads = 0
        }

        slog.Debug("Download delay", "duration", duration, "currentDownloads", currentDownloads)
        time.Sleep(duration)
	}
}



func mkdirAll(dirName string) {
	path := filepath.Join(c.Download.Dst, dirName)
    _, err := os.Stat(path)
	if err != nil {
		slog.Debug("Create directories", "path", path)
		err = os.MkdirAll(path, os.ModeDir|os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func downloadEntry(path string, entry anime.AnimeDataEntry) {
    if err := a.Download(entry, path); err != nil {
        slog.Error("Download error!", "err", err.Error())
    }
}

func sleep() {
	for true {
		now := time.Now()

		var day int
		if now.Hour() < c.Update.Hour {
			day = now.Day()
		} else {
			day = now.Day() + 1
		}
		next := time.Date(now.Year(), now.Month(), day, c.Update.Hour, 0, 0, 0, time.Local)

		duration := next.Sub(now)
		slog.Debug("Sleep until next update day.", "duration", duration)

		time.Sleep(duration)

		if time.Now().Weekday() == c.Update.Weekday {
			slog.Debug("Running new update now...")
			break
		}
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
		Level: slog.LevelInfo,
    }

	if c.Debug {
        options.Level = slog.LevelDebug
	}

    slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, options)))
}
