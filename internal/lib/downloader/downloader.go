package downloader

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

// DownloadIfNotExists downloads the file into a .part file and renames it when done,
// so an interrupted download never leaves a partial file under the real name.
func DownloadIfNotExists(logger *slog.Logger, filePath string, url string) error {
	fileName := filepath.Base(filePath)
	if _, err := os.Stat(filePath); err == nil {
		logger.Info(fmt.Sprintf("skipping the %s download: file already exists", fileName))
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	partPath := filePath + ".part"
	f, err := os.Create(partPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	bar := progressbar.DefaultBytes(resp.ContentLength, "downloading "+fileName)
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		os.Remove(partPath)
		return fmt.Errorf("write file: %w", err)
	}

	if err := os.Rename(partPath, filePath); err != nil {
		return fmt.Errorf("rename file: %w", err)
	}

	logger.Info(fmt.Sprintf("%s has been successfully downloaded", fileName))
	return nil
}
