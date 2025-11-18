package utils

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

const (
	JPEG = "image/jpeg"
	PNG  = "image/png"
)

func GetExtension(t string) (ext string, err error) {
	switch t {
	case PNG:
		ext = "png"
	case JPEG:
		ext = "jpg"
	default:
		ext = "unknown"
		err = fmt.Errorf("unknown extension from type \"%s\"", t)
	}

	return ext, err
}

func ConvertImagesToPDF(path string, images ...string) error {
    slog.Debug("Convert images to pdf", "dst", path)
	images = append(images, "-quality", "100", "-density", "150", path+".pdf")

	var stderr bytes.Buffer
	cmd := exec.Command("magick", images...)
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		slog.Error("Convert images to PDF failed!", "err", err.Error(), "stderr", stderr.String())
		_ = os.Remove(path + ".pdf")
	}

	return err
}
