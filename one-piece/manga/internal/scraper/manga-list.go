package scraper

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gocolly/colly/v2"
)

type MangaList_Special struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MangaList_Arc struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
}

type MangaList_Chapter struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Number      int    `json:"number"`
	CategoryId  int    `json:"category_id"`
	ArcId       int    `json:"arc_id"`
	SpecialsId  int    `json:"specials_id"`
	Lang        string `json:"lang"`
	Pages       int    `json:"pages"`
	IsAvailable bool   `json:"is_available"`
	Date        string `json:"date"`
	Href        string `json:"href"`
}

type MangaList struct {
	Specials []MangaList_Special `json:"specials"`
	Arcs     []MangaList_Arc     `json:"arcs"`
	Chapters []MangaList_Chapter `json:"entries"`
}

func (mangaList *MangaList) GetArc(id int) (arc *MangaList_Arc, index int) {
	for i := 0; i < len(mangaList.Arcs); i++ {
		if mangaList.Arcs[i].Id == id {
			return &mangaList.Arcs[i], i
		}
	}

	return nil, -1
}

// ParseMangaList
func ParseMangaList() (data *MangaList, err error) {
	data = &MangaList{}
	c := colly.NewCollector()

	c.OnHTML("script", func(e *colly.HTMLElement) {
		e.Text = strings.Trim(e.Text, " ")
		if len(e.Text) < 13 {
			return
		}

		if e.Text[0:13] == "window.__data" {
			e.Text, _ = strings.CutPrefix(e.Text, "window.__data")
			e.Text = strings.TrimLeft(e.Text, " =")
			e.Text = strings.TrimRight(e.Text, "; ")

			if err := json.Unmarshal([]byte(e.Text), data); err != nil {
				slog.Error(err.Error())
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		slog.Debug(fmt.Sprintf("Request to \"%s\"", r.URL))
	})

	c.OnError(func(r *colly.Response, e error) {
		err = e
	})

	if err := c.Visit("https://onepiece-tube.com/manga/kapitel-mangaliste"); err != nil {
		return data, err
	}

	c.Wait()

	return data, err
}
