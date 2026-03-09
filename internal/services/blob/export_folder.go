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

type ExportResult struct {
	ExportedFiles int `json:"exported_files"`
	SkippedFiles  int `json:"skipped_files"`
}

func (s *Service) ExportToFolder(ctx context.Context, dir string, overwrite bool) (*ExportResult, error) {
	const op = "blob.Service.ExportToFolder"

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	rows, err := s.store.DB().Query("select cdnid, blob from asset_file;")
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	result := &ExportResult{}

	for rows.Next() {
		var cdnid string
		var blob []byte

		if err := rows.Scan(&cdnid, &blob); err != nil {
			return nil, wrap.IfErr(op, err)
		}

		assetPath := filepath.Join(dir, cdnid)
		wrote, err := writeBlobToDisk(assetPath, blob, overwrite)
		if err != nil {
			return nil, wrap.IfErr(op, fmt.Errorf("cdnid: %s, err: %w", cdnid, err))
		}

		if wrote {
			result.ExportedFiles++
			s.logger.Debug(op, "cdnid", cdnid, "status", "exported")
		} else {
			result.SkippedFiles++
			s.logger.Debug(op, "cdnid", cdnid, "status", "skipped", "reason", "exists")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	s.logger.Info("export cache files to folder: finished",
		"exported_files", result.ExportedFiles,
		"skipped_files", result.SkippedFiles,
	)

	return result, nil
}

func writeBlobToDisk(filePath string, blob []byte, overwrite bool) (bool, error) {
	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		return false, fmt.Errorf("%s is a directory", filePath)
	} else if err == nil && !overwrite {
		return false, nil
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return false, err
	}

	if err := os.WriteFile(filePath, blob, 0o644); err != nil {
		return false, err
	}

	return true, nil
}
