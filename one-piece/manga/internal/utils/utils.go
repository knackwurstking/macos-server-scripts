package utils

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
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

	// Check if ImageMagick is available
	if err := checkImageMagick(); err != nil {
		slog.Error("ImageMagick not available", "err", err.Error())
		return err
	}

	var stderr bytes.Buffer
	cmd := exec.Command("magick", images...)
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		slog.Error("Convert images to PDF failed!", "err", err.Error(), "stderr", stderr.String())
		_ = os.Remove(path + ".pdf")
		return err
	}

	return nil
}

// checkImageMagick verifies that ImageMagick is installed and available
func checkImageMagick() error {
	// Try to execute magick command to check if it's available
	cmd := exec.Command("magick", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ImageMagick not found or not properly installed: %w. Output: %s", err, strings.TrimSpace(string(output)))
	}

	// Check if the output contains version information
	if !strings.Contains(string(output), "ImageMagick") {
		return fmt.Errorf("ImageMagick not found in command output")
	}

	return nil
}
