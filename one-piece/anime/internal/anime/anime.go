package anime

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

const (
	PathEpisodenStreams Path = "/anime/episoden-streams"
)

type Path string

type Anime struct {
	Origin string     `json:"origin"`
	Data   *AnimeData `json:"data"`
}

func NewAnime(origin string) *Anime {
	return &Anime{
		Origin: origin,
		Data:   NewAnimeData(),
	}
}

func (anime *Anime) GetUrl(path Path) string {
	switch path {
	case PathEpisodenStreams:
		return fmt.Sprintf("%s%s", anime.Origin, PathEpisodenStreams)
	default:
		panic(fmt.Sprintf("Name \"%s\" not found!", path))
	}
}

func (anime *Anime) GetEpisodenStreams() (*AnimeData, error) {
	var (
		c   = colly.NewCollector()
		err error
	)

	c.OnHTML("script", func(h *colly.HTMLElement) {
		var (
			dataVar = "window.__data"
			text    = strings.Trim(h.Text, " ")
		)

		if len(text) < len(dataVar) {
			return
		}

		if text[0:13] == dataVar {
			text, _ = strings.CutPrefix(text, dataVar)
			text = strings.TrimLeft(text, " =")
			text = strings.TrimRight(text, "; ")

			if err := json.Unmarshal([]byte(text), anime.Data); err != nil {
				slog.Error("Unmarshal data failed!", "err", err.Error())
			} else {
				slog.Debug("Successfully unmarshaled data", "entries_count", len(anime.Data.Entries))
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		slog.Debug(fmt.Sprintf("Request to \"%s\"", r.URL))
	})

	c.OnError(func(r *colly.Response, e error) {
		if e != nil {
			err = e
			slog.Error("HTTP request error", "error", e.Error())
		}
	})

	if err := c.Visit(anime.GetUrl(PathEpisodenStreams)); err != nil {
		return anime.Data, err
	}

	c.Wait()

	return anime.Data, err
}

func (anime *Anime) Download(entry AnimeDataEntry, path string) error {
	var (
		c   = colly.NewCollector()
		err error
	)

	c.OnHTML("iframe", func(h *colly.HTMLElement) {
		src := h.Attr("src")
		if src == "" {
			slog.Warn("Empty iframe src found")
			return
		}

		iframeCollector := colly.NewCollector()

		iframeCollector.OnHTML("video > source", func(h *colly.HTMLElement) {
			src := h.Attr("src")
			if src == "" {
				slog.Warn("Empty video source found")
				return
			}

			if h.Attr("type") != "video/mp4" {
				slog.Warn("HTML tag <source has not type \"video/mp4\" attribute")
				return
			}

			slog.Debug("Got url from video source", "src", src, "path", path)
			if err := anime.downloadSource(src, path); err != nil {
				slog.Error("download src to dst failed", "err", err, "src", src, "dst", path)
				_ = os.Remove(path)
			}
		})

		iframeCollector.OnRequest(func(r *colly.Request) {
			slog.Debug(fmt.Sprintf("Request to \"%s\"", r.URL))
		})

		iframeCollector.OnError(func(r *colly.Response, err error) {
			slog.Error("Iframe collector error", "error", err.Error())
		})

		if err := iframeCollector.Visit(src); err != nil {
			slog.Error(fmt.Sprintf("Visit \"%s\" failed!", src), "err", err.Error())
		}

		iframeCollector.Wait()
	})

	c.OnRequest(func(r *colly.Request) {
		slog.Debug(fmt.Sprintf("Request to \"%s\"", r.URL))
	})

	c.OnError(func(r *colly.Response, e error) {
		if e != nil {
			err = e
			slog.Error("Main collector error", "error", e.Error())
		}
	})

	if err := c.Visit(entry.Href); err != nil {
		return err
	}

	c.Wait()

	return err
}

func (anime *Anime) downloadSource(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		slog.Warn("file already exists", "dst", dst)
		return nil
	}

	slog.Debug("Starting download", "src", src, "dst", dst)
	response, err := http.Get(src)
	if err != nil {
		slog.Error("HTTP GET failed", "src", src, "error", err.Error())
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(dst)
	if err != nil {
		slog.Error("Failed to create file", "dst", dst, "error", err.Error())
		return err
	}
	defer file.Close()

	n, err := io.Copy(bufio.NewWriter(file), response.Body)
	slog.Debug("io.Copy completed", "dst", dst, "written", n, "err", err)

	if err != nil {
		slog.Error("Copy failed", "dst", dst, "error", err.Error())
	}

	return err
}

type AnimeData struct {
	Arcs    AnimeDataArcs    `json:"arcs"`
	Entries []AnimeDataEntry `json:"entries"`
}

func NewAnimeData() *AnimeData {
	return &AnimeData{
		Arcs:    make([]AnimeDataArc, 0),
		Entries: make([]AnimeDataEntry, 0),
	}
}

type AnimeDataArcs []AnimeDataArc

func (arcs AnimeDataArcs) Get(id int) *AnimeDataArc {
	for i := 0; i < len(arcs); i++ {
		if arcs[i].ID == id {
			return &arcs[i]
		}
	}

	return nil
}

// GetIndex in reversed order
func (arcs AnimeDataArcs) GetIndex(id int) int {
	for i := 0; i < len(arcs); i++ {
		if arcs[i].ID == id {
			return len(arcs) - 1 - i
		}
	}

	panic(fmt.Sprintf("GetOrder failed for id \"%d\"", id))
}

type AnimeDataArc struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AnimeDataEntry struct {
	Name        string `json:"name"`
	Number      int    `json:"number"`
	ArcID       int    `json:"arc_id"`
	LangSub     string `json:"lang_sub"`
	LangDub     string `json:"lang_dub"`
	IsAvailable bool   `json:"is_available"`
	Href        string `json:"href"`
}
