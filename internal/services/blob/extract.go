package blob

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
)

type ExtractResult struct {
	ExtractedFiles int `json:"extracted_files"`
}

type ExtractOptions struct {
	ExtractPath     string `json:"extract_path"`
	ExtractMetadata bool   `json:"extract_metadata"`
}

func (s *Service) ExtractFiles(ctx context.Context, options ExtractOptions) (*ExtractResult, error) {
	const op = "blob.Service.ExtractFiles"

	if err := os.MkdirAll(options.ExtractPath, 0o755); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	rows, err := s.store.DB().QueryContext(ctx, "select cdnid, blob, json_pretty(metadata) from blob.asset_file;")
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	result := &ExtractResult{}

	for rows.Next() {
		var cdnid string
		var blob []byte
		var metadata []byte

		if err := rows.Scan(&cdnid, &blob, &metadata); err != nil {
			return nil, wrap.IfErr(op, err)
		}

		cacheFilePath := filepath.Join(options.ExtractPath, cdnid)
		if err := writeToDisk(cacheFilePath, blob); err != nil {
			return nil, wrap.IfErr(op, fmt.Errorf("cdnid: %s, err: %w", cdnid, err))
		}

		if options.ExtractMetadata && len(metadata) > 0 {
			metadataFilePath := filepath.Join(options.ExtractPath, cdnid+".meta.json")
			if err := writeToDisk(metadataFilePath, metadata); err != nil {
				return nil, wrap.IfErr(op, fmt.Errorf("cdnid: %s, err: %w", cdnid, err))
			}
		}

		s.logger.Debug(op, "extracted_file", cdnid)
		result.ExtractedFiles++
	}

	if err := rows.Err(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	s.logger.Info("cache files extraction: finished",
		"extracted_files", result.ExtractedFiles,
	)

	return result, nil
}

func writeToDisk(filePath string, content []byte) error {
	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		return fmt.Errorf("%s is a directory", filePath)
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return os.WriteFile(filePath, content, 0o644)
}
