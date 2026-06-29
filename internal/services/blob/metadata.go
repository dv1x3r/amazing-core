package blob

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/dv1x3r/amazing-core/internal/lib/python"
)

var metadataPythonPackages = []python.Package{
	{
		Spec:        "git+https://github.com/dv1x3r/UnityPy.git@test",
		ImportCheck: "import UnityPy",
	},
	{
		Spec:        "mutagen",
		ImportCheck: "import mutagen",
	},
	{
		Spec:        "Pillow",
		ImportCheck: "from PIL import Image",
	},
}

type MetadataManifest struct {
	Files []MetadataManifestEntry `json:"files"`
}

type MetadataManifestEntry struct {
	FilePath string `json:"file_path"`
	JSONPath string `json:"json_path"`
}

func (s *Service) generateMetadata(ctx context.Context, cacheScript []byte, entries map[string]MetadataManifestEntry) error {
	const op = "blob.Service.generateMetadata"
	if len(entries) == 0 {
		return nil
	}

	tmpDir, err := os.MkdirTemp("", "amazing-core-blob-metadata-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	metadataDir := filepath.Join(tmpDir, "metadata")
	if err := os.MkdirAll(metadataDir, 0o755); err != nil {
		return err
	}

	for hash, entry := range entries {
		entry.JSONPath = filepath.Join(metadataDir, hash+".meta.json")
		entries[hash] = entry
	}

	manifestPath := filepath.Join(tmpDir, "manifest.json")
	if err := writeMetadataManifest(manifestPath, entries); err != nil {
		return err
	}

	if s.logger != nil {
		s.logger.Debug(op, "tmp_dir", tmpDir, "manifest", manifestPath, "files", len(entries))
	}

	streamOutput := s.logger != nil && s.logger.Handler().Enabled(ctx, slog.LevelDebug)
	_, err = s.python.RunScriptWithOptions(ctx, cacheScript, python.RunOptions{
		StreamOutput: streamOutput,
	}, "--manifest", manifestPath)
	if err != nil {
		return fmt.Errorf("run cache tool: %w", err)
	}

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for hash, entry := range entries {
		if err := s.updateMetadata(ctx, tx, hash, entry); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func writeMetadataManifest(path string, entries map[string]MetadataManifestEntry) error {
	manifest := MetadataManifest{
		Files: make([]MetadataManifestEntry, 0, len(entries)),
	}
	for _, entry := range entries {
		manifest.Files = append(manifest.Files, entry)
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (s *Service) updateMetadata(ctx context.Context, tx *sql.Tx, hash string, entry MetadataManifestEntry) error {
	metadata, err := os.ReadFile(entry.JSONPath)
	if err != nil {
		return fmt.Errorf("read metadata file: hash: %s, path: %s, err: %w", hash, entry.FilePath, err)
	}

	const query = "update blob.asset_file set metadata = jsonb(?), updated_at = unixepoch() where hash = ?;"
	if _, err = tx.ExecContext(ctx, query, string(metadata), hash); err != nil {
		return fmt.Errorf("update metadata: hash: %s, path: %s, err: %w", hash, entry.FilePath, err)
	}

	return nil
}
