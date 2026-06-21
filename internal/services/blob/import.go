package blob

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

// some old game cache files do not have CDNID (OID) in their names
// because of that I decided to use custom OID range for this kind of files
// it starts with class-0 type-0 server-1 number-1
const firstGeneratedOID int64 = 1099511627777

type ImportResult struct {
	ImportedFiles int `json:"imported_files"`
	SkippedFiles  int `json:"skipped_files"`
}

func ImportFromFolder(ctx context.Context, logger *slog.Logger, db *sql.DB, dir string) (*ImportResult, error) {
	const op = "blob.ImportFromFolder"

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	result := &ImportResult{}
	nextGeneratedOID := firstGeneratedOID

	err = filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			return nil
		}

		fileName := entry.Name()
		if shouldSkipCacheFile(fileName) {
			result.SkippedFiles++
			logger.Debug(op, "path", path, "status", "skipped")
			return nil
		}

		isGeneratedOID := false
		cdnid := cdnidFromFileName(fileName)
		if len(cdnid) != 18 {
			cdnid = types.OIDFromInt64(nextGeneratedOID).CDNID()
			isGeneratedOID = true
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cdnid: %s, path: %s, err: %w", cdnid, path, err)
		}

		hashBytes := sha1.Sum(data)
		hashString := hex.EncodeToString(hashBytes[:])

		inserted, err := insertBlobToDB(tx, cdnid, data, hashString)
		if err != nil {
			return fmt.Errorf("cdnid: %s, path: %s, err: %w", cdnid, path, err)
		}

		if inserted {
			if isGeneratedOID {
				nextGeneratedOID++
			}
			result.ImportedFiles++
			logger.Debug(op, "cdnid", cdnid, "path", path, "status", "imported")
		} else {
			result.SkippedFiles++
			logger.Debug(op, "cdnid", cdnid, "path", path, "status", "skipped")
		}

		return nil
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	logger.Info("import cache files from folder: finished",
		"imported_files", result.ImportedFiles,
		"skipped_files", result.SkippedFiles,
	)

	return result, nil
}

func shouldSkipCacheFile(fileName string) bool {
	if strings.HasPrefix(fileName, ".") {
		return true
	}

	if strings.HasPrefix(fileName, "(1)") || strings.HasSuffix(fileName, "(1)") {
		return true
	}

	if strings.HasPrefix(fileName, "(2)") || strings.HasSuffix(fileName, "(2)") {
		return true
	}

	switch fileName {
	case "index.txt", "Index.txt", "__info", ".DS_Store":
		return true
	default:
		return false
	}
}

func cdnidFromFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	return fileName[:len(fileName)-len(ext)]
}

func insertBlobToDB(tx *sql.Tx, cdnid string, blob []byte, hash string) (bool, error) {
	res, err := tx.Exec(`insert or ignore into asset_file (cdnid, blob, hash) values (?, ?, ?);`, cdnid, blob, hash)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}
