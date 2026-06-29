package blob

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/tools"
)

type ImportResult struct {
	ImportedFiles     int `json:"imported_files"`
	SkippedFiles      int `json:"skipped_files"`
	GeneratedMetadata int `json:"generated_metadata"`
}

type ImportOptions struct {
	ImportPath       string `json:"import_path"`
	GenerateMetadata bool   `json:"generate_metadata"`
}

func (s *Service) ImportFromFolder(ctx context.Context, options ImportOptions) (*ImportResult, error) {
	const op = "blob.Service.ImportFromFolder"

	var cacheScript []byte
	if options.GenerateMetadata {
		if err := s.python.Ensure(ctx, metadataPythonPackages...); err != nil {
			return nil, wrap.IfErr(op, err)
		}
		script, err := tools.FS.ReadFile("cache.py")
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		cacheScript = script
	}

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	// some old game cache files do not have CDNID (OID) in their names
	// because of that I decided to use custom OID range for this kind of files
	// it starts with class-0 type-0 server-1 number-1
	nextGeneratedOID, err := getNextGeneratedOID(ctx, tx)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	result := &ImportResult{}
	metadataManifest := map[string]MetadataManifestEntry{}

	err = filepath.WalkDir(options.ImportPath, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			if strings.HasSuffix(entry.Name(), "_assets") {
				return filepath.SkipDir
			}
			return nil
		}

		fileName := entry.Name()
		if shouldSkipCacheFile(fileName) {
			result.SkippedFiles++
			s.logger.Debug(op, "path", path, "status", "skipped")
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

		inserted, err := insertIntoDB(ctx, tx, cdnid, data, hashString)
		if err != nil {
			return fmt.Errorf("cdnid: %s, path: %s, err: %w", cdnid, path, err)
		}

		if isGeneratedOID && inserted {
			nextGeneratedOID++
		}

		if inserted {
			result.ImportedFiles++
			s.logger.Debug(op, "cdnid", cdnid, "path", path, "status", "imported")
		} else {
			result.SkippedFiles++
			s.logger.Debug(op, "cdnid", cdnid, "path", path, "status", "skipped")
		}

		if options.GenerateMetadata {
			if _, exists := metadataManifest[hashString]; !exists {
				metadataManifest[hashString] = MetadataManifestEntry{FilePath: path}
			}
		}

		return nil
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	if options.GenerateMetadata {
		if err := s.generateMetadata(ctx, cacheScript, metadataManifest); err != nil {
			return nil, wrap.IfErr(op, err)
		}
		result.GeneratedMetadata = len(metadataManifest)
	}

	s.logger.Info("cache files import: finished",
		"imported_files", result.ImportedFiles,
		"skipped_files", result.SkippedFiles,
		"generated_metadata", result.GeneratedMetadata,
	)

	return result, nil
}

func shouldSkipCacheFile(fileName string) bool {
	fileName = strings.ToLower(fileName)

	if strings.HasPrefix(fileName, ".") {
		return true
	}

	if strings.HasSuffix(fileName, ".meta.json") || strings.HasSuffix(fileName, ".zip") {
		return true
	}

	if strings.HasPrefix(fileName, "(1)") || strings.HasSuffix(fileName, "(1)") {
		return true
	}

	if strings.HasPrefix(fileName, "(2)") || strings.HasSuffix(fileName, "(2)") {
		return true
	}

	switch fileName {
	case "index.txt", "__info":
		return true
	default:
		return false
	}
}

func cdnidFromFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	return fileName[:len(fileName)-len(ext)]
}

func getNextGeneratedOID(ctx context.Context, tx *sql.Tx) (int64, error) {
	var lastGeneratedOID int64
	row := tx.QueryRowContext(ctx, "select coalesce(max(gsfoid), 1099511627777) as last from asset where (gsfoid >> 40) & 0xFF = 1;")
	if err := row.Scan(&lastGeneratedOID); err != nil {
		return 0, err
	}
	return lastGeneratedOID + 1, nil
}

func insertIntoDB(ctx context.Context, tx *sql.Tx, cdnid string, blob []byte, hash string) (inserted bool, err error) {
	res, err := tx.ExecContext(ctx, "insert or ignore into blob.asset_file (cdnid, blob, hash) values (?, ?, ?);", cdnid, blob, hash)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}
