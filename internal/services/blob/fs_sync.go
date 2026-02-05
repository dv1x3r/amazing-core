package blob

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
)

type ImportResult struct {
	ImportedFiles int `json:"imported_files"`
	SkippedFiles  int `json:"skipped_files"`
}

func (s *Service) ImportFromFolder(ctx context.Context, dir string) (*ImportResult, error) {
	const op = "blob.Service.ImportFromFolder"

	tx, err := s.store.DB().Begin()
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	result := &ImportResult{}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		fileName := dirEntry.Name()
		fileExt := filepath.Ext(fileName)

		cdnid := fileName[:len(fileName)-len(fileExt)]
		if len(cdnid) != 18 {
			result.SkippedFiles++
			s.logger.Debug("import", "cdnid", cdnid, "status", "skipped", "reason", "invalid_cdnid")
			continue
		}

		blob, err := os.ReadFile(filepath.Join(dir, fileName))
		if err != nil {
			return nil, wrap.IfErr(op, fmt.Errorf("cdnid: %s, err: %w", cdnid, err))
		}

		hashBytes := sha1.Sum(blob)
		hashString := hex.EncodeToString(hashBytes[:])

		inserted, err := insertBlobToDB(tx, cdnid, blob, hashString)
		if err != nil {
			return nil, wrap.IfErr(op, fmt.Errorf("cdnid: %s, err: %w", cdnid, err))
		}

		if inserted {
			result.ImportedFiles++
			s.logger.Debug("import", "cdnid", cdnid, "status", "imported")
		} else {
			result.SkippedFiles++
			s.logger.Debug("import", "cdnid", cdnid, "status", "skipped", "reason", "hash_match")
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	s.logger.Info("import", "status", "finished", "imported", result.ImportedFiles, "skipped", result.SkippedFiles)
	return result, nil
}

func insertBlobToDB(tx *sql.Tx, cdnid string, blob []byte, hash string) (bool, error) {
	var dbHash string
	row := tx.QueryRow("select hash from asset_file where cdnid = ?;", cdnid)
	if err := row.Scan(&dbHash); err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if dbHash == hash {
		return false, nil
	} else if dbHash != "" {
		return false, fmt.Errorf("hashes do not match")
	}

	_, err := tx.Exec("insert into asset_file (cdnid, blob, hash) values (?, ?, ?);", cdnid, blob, hash)
	if err != nil {
		return false, err
	}

	return true, nil
}

type ExportResult struct {
	ExportedFiles int `json:"exported_files"`
	SkippedFiles  int `json:"skipped_files"`
}

func (s *Service) ExportFromFolder(ctx context.Context, dir string, overwrite bool) (*ExportResult, error) {
	const op = "blob.Service.ExportFromFolder"

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
			s.logger.Debug("export", "cdnid", cdnid, "status", "exported")
		} else {
			result.SkippedFiles++
			s.logger.Debug("export", "cdnid", cdnid, "status", "skipped", "reason", "exists")
		}
	}

	if err := rows.Err(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	s.logger.Info("export", "status", "finished", "exported", result.ExportedFiles, "skipped", result.SkippedFiles)
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
