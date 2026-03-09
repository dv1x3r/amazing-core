package blob

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
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

	tx, err := s.store.DB().BeginTx(ctx, nil)
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
			s.logger.Debug(op, "cdnid", cdnid, "status", "skipped", "reason", "invalid_cdnid")
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, fileName))
		if err != nil {
			return nil, wrap.IfErr(op, fmt.Errorf("cdnid: %s, err: %w", cdnid, err))
		}

		hashBytes := sha1.Sum(data)
		hashString := hex.EncodeToString(hashBytes[:])

		inserted, err := insertBlobToDB(tx, cdnid, data, hashString)
		if err != nil {
			return nil, wrap.IfErr(op, fmt.Errorf("cdnid: %s, err: %w", cdnid, err))
		}

		if inserted {
			result.ImportedFiles++
			s.logger.Debug(op, "cdnid", cdnid, "status", "imported")
		} else {
			result.SkippedFiles++
			s.logger.Debug(op, "cdnid", cdnid, "status", "skipped", "reason", "hash_match")
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	s.logger.Info("import cache files from folder: finished",
		"imported_files", result.ImportedFiles,
		"skipped_files", result.SkippedFiles,
	)

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
