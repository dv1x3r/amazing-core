package asset

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
)

type ImportResult struct {
	ImportedAssets   int `json:"imported_assets"`
	ImportedMetadata int `json:"imported_metadata"`
}

type CacheItem struct {
	File struct {
		Name string `json:"name"`
		Size int    `json:"size"`
		Type string `json:"type"`
		Hash string `json:"hash"`
		OID  int    `json:"oid"`
	} `json:"file"`
	Bundle json.RawMessage `json:"bundle"`
}

func ImportCacheItems(logger *slog.Logger, db *sql.DB, items []CacheItem) (*ImportResult, error) {
	const op = "asset.ImportCacheItems"

	const assetSQL = `
		insert into asset (cdnid, gsfoid, hash, size, file_type_id)
		values (?, ?, ?, ?, ?)
		on conflict(cdnid) do update set
			gsfoid       = excluded.gsfoid,
			hash         = excluded.hash,
			size         = excluded.size,
			file_type_id = excluded.file_type_id
		returning id`

	const metadataSQL = `
		insert into asset_metadata (asset_id, metadata)
		values (?, jsonb(?))
		on conflict(asset_id) do update set
			metadata = excluded.metadata`

	tx, err := db.Begin()
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	fileTypes := map[string]int{}
	rows, err := tx.Query("select id, name from file_type")
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, wrap.IfErr(op, err)
		}
		fileTypes[strings.ToLower(name)] = id
	}

	if err := rows.Close(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	stmtAsset, err := tx.Prepare(assetSQL)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer stmtAsset.Close()

	stmtMetadata, err := tx.Prepare(metadataSQL)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer stmtMetadata.Close()

	result := &ImportResult{}
	for _, item := range items {
		fileTypeID, ok := fileTypes[strings.ToLower(item.File.Type)]
		if !ok {
			return nil, wrap.IfErr(op, fmt.Errorf("unknown file type %q", item.File.Type))
		}

		var assetID int64
		err := stmtAsset.QueryRow(
			item.File.Name,
			item.File.OID,
			item.File.Hash,
			item.File.Size,
			fileTypeID,
		).Scan(&assetID)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		} else {
			result.ImportedAssets++
		}

		if item.Bundle == nil {
			continue
		}

		res, err := stmtMetadata.Exec(assetID, item.Bundle)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		if affected > 0 {
			result.ImportedMetadata++
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	logger.Info("import cache.json assets and metadata finished",
		"imported_assets", result.ImportedAssets,
		"imported_metadata", result.ImportedMetadata,
	)

	return result, nil
}
