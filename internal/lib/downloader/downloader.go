package downloader

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/schollz/progressbar/v3"
)

func DownloadIfNotExists(logger *slog.Logger, filePath string, url string) error {
	fileName := filepath.Base(filePath)
	if _, err := os.Stat(filePath); err == nil {
		logger.Info(fmt.Sprintf("skipping the %s download: already exists", fileName))
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

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		f.Close()
		os.Remove(filePath) // clean up partial file if user closes the app
		os.Exit(1)
	}()

	bar := progressbar.DefaultBytes(resp.ContentLength, "downloading "+fileName)
	if _, err = io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
		f.Close()
		os.Remove(filePath) // clean up partial file when something goes wrong
		return fmt.Errorf("write file: %w", err)
	}

	f.Close()
	signal.Stop(sig)

	logger.Info(fmt.Sprintf("%s has been successfully downloaded", fileName))
	return nil
}
