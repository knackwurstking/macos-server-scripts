package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"op-manga-dl/internal/scraper"
	"op-manga-dl/internal/utils"

	"github.com/lmittmann/tint"
)

const (
	ErrCodeMagickNotAvailable int = 1
)

var (
	c *Config
)

func main() {
	c = NewConfig()
	parseFlags()

	// Check if ImageMagick is available
	if err := utils.CheckImageMagick(); err != nil {
		slog.Error("ImageMagick not available", "err", err.Error())
		os.Exit(ErrCodeMagickNotAvailable)
	}

	// Continuous loop with proper error handling
	for {
		err := downloadAllChapters()
		if err != nil {
			slog.Error("Error in download cycle", "err", err.Error())
		}

		sleep()
	}
}

func downloadAllChapters() error {
	slog.Debug("Download all chapters possible")

	ml, err := scraper.ParseMangaList()
	if err != nil {
		return fmt.Errorf("fetch & parse manga llist: %v", err)
	}

	currentDownloads := 0
	for _, chapter := range ml.Chapters {
		if chapter.Pages == 0 {
			continue
		}

		arc, i := ml.GetArc(chapter.ArcId)
		if arc == nil {
			slog.Error(fmt.Sprintf(
				"Arc for %s with the id %d not found! (This should never happen)",
				chapter.Name, chapter.ArcId,
			))
			continue
		}

		path := filepath.Join(
			c.Download.Dst,
			fmt.Sprintf("%03d %s", len(ml.Arcs)-i, arc.Name),
			fmt.Sprintf("%04d %s", chapter.Number, chapter.Name),
		)

		_, err := os.Stat(path + ".pdf")
		if err == nil {
			// File (pdf) already exists, continue to next chapter
			continue
		}

		// Make directory where the pages will be stored
		// (ignore errors if already exists)
		err = os.MkdirAll(path, 0755)
		if err != nil {
			slog.Error("Failed to create directory", "path", path, "err", err.Error())
			continue
		}

		currentDownloads += 1
		err = downloadChapter(chapter, path)
		if err != nil {
			slog.Error("Failed to download chapter", "name", chapter.Name, "error", err)
			continue
		}

		// Handle the short and long download delay based on limit
		duration := time.Minute * time.Duration(c.Download.Delay)
		if currentDownloads >= c.Download.LimitPerDay {
			duration = time.Minute * time.Duration(c.Download.LongDelay)
			currentDownloads = 0
		}

		slog.Debug("Handle the download delay",
			"duration", duration,
			"currentDownloads", currentDownloads)
		time.Sleep(duration)
	}

	return nil
}

func downloadChapter(chapter scraper.MangaList_Chapter, path string) error {
	slog.Info(fmt.Sprintf("Download chapter \"%s\" with %d pages.",
		chapter.Name, chapter.Pages))

	// download jpg/png from dURL - scrape the same script section like before
	chapterData, err := scraper.ParseChapter(chapter.Href)
	if err != nil {
		slog.Error("Parsing chapter failed!",
			"chapter.Href", chapter.Href,
			"err", err.Error())
		return err
	}

	pages := make([]string, len(chapterData.Chapter.Pages))
	for i, page := range chapterData.Chapter.Pages {
		slog.Debug(fmt.Sprintf("Download page nr %d", i+1),
			"page.Url", page.Url)

		r, err := http.Get(page.Url)
		if err != nil {
			slog.Error("Downloading page failed!",
				"page", i+1, "err", err.Error())
			return err
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Read all body data failed!",
				"page", i+1, "err", err.Error())
			return err
		}
		if len(data) == 0 {
			slog.Error("No data!", "page", i+1)
			return fmt.Errorf("no data received for page %d", i+1)
		}
		e, _ := utils.GetExtension(page.Type)
		p := filepath.Join(path, fmt.Sprintf("%02d.%s", i+1, e))
		err = os.WriteFile(p, data, 0644)
		if err != nil {
			slog.Error(fmt.Sprintf("Write file \"%s\" failed!", p),
				"err", err.Error())
			return err
		}
		pages[i] = p
	}

	if err := utils.ConvertImagesToPDF(path, pages...); err != nil {
		slog.Error("Convert pages to pdf failed!", "err", err.Error())
		return err
	}

	return nil
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
	slog.Debug("Sleep until next update day.", "duration", duration, "next_update", nextUpdate, "current_time", now)

	time.Sleep(duration)

	if time.Now().Weekday() == c.Update.Weekday {
		slog.Debug("Running new update now...")
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
