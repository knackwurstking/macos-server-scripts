package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type ConfigUpdate struct {
	Weekday time.Weekday `json:"weekday"`
	Hour    int          `json:"hour"`
}

type ConfigDownload struct {
	Delay       int    `json:"delay"`
	LongDelay   int    `json:"long-delay"`
	Dst         string `json:"dst"`
	LimitPerDay int    `json:"limit-per-day"`
}

type Config struct {
	Update   ConfigUpdate   `json:"update"`
	Download ConfigDownload `json:"download"`
	Debug    bool           `json:"debug"`
}

func NewConfig() *Config {
	return &Config{
		Update: ConfigUpdate{
			Weekday: time.Sunday,
			Hour:    18,
		},
		Download: ConfigDownload{
			Delay:       30,      // 20 min
			LongDelay:   12 * 60, // 12 hours
			Dst:         filepath.Join("data", "download"),
			LimitPerDay: 5,
		},
	}
}

func (c *Config) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err == nil {
		err = json.NewDecoder(file).Decode(c)
	}
	return err
}
