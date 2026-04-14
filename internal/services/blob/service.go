package blob

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/url"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/dustin/go-humanize"
	"github.com/huandu/go-sqlbuilder"
)

var (
	ErrFileNotFound = errors.New("file is not found")
	ErrFileExists   = errors.New("file with the same name already exists")
)

type Service struct {
	logger      *slog.Logger
	store       db.Store
	deliveryURL string
}

func NewService(logger *slog.Logger, store db.Store, deliveryURL string) *Service {
	return &Service{
		logger:      logger,
		store:       store,
		deliveryURL: deliveryURL,
	}
}

type File struct {
	ID      int    `json:"id"`
	CDNID   string `json:"cdnid"`
	Hash    string `json:"hash"`
	Size    int    `json:"size"`
	SizeStr string `json:"size_str"`
	URL     string `json:"url"`
}

func (s *Service) FetchFileBlob(ctx context.Context, cdnid string) ([]byte, error) {
	const op = "blob.Service.FetchFileBlob"

	var data []byte

	const query = "select blob from asset_file where cdnid = ?;"
	row := s.store.DB().QueryRowContext(ctx, query, cdnid)
	if err := row.Scan(&data); err == sql.ErrNoRows {
		return nil, ErrFileNotFound
	} else if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	return data, nil
}

func (s *Service) FetchFilesList(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[File], error) {
	const op = "blob.Service.FetchFilesList"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[File]{
		From: "asset_file",
		Select: []string{
			"id",
			"cdnid",
			"hash",
			"length(blob) as size",
		},
		WhereMapping: map[string]string{
			"id":    "id",
			"cdnid": "cdnid",
			"hash":  "hash",
			"size":  "length(blob)",
		},
		OrderByMapping: map[string]string{
			"id":       "id",
			"cdnid":    "cdnid COLLATE BINARY",
			"url":      "cdnid COLLATE BINARY",
			"hash":     "hash",
			"size":     "length(blob)",
			"size_str": "length(blob)",
		},
		Flavor: sqlbuilder.SQLite,
		Scan: func(rows *sql.Rows) (File, error) {
			var record File
			if err := rows.Scan(
				&record.ID,
				&record.CDNID,
				&record.Hash,
				&record.Size,
			); err != nil {
				return record, err
			}
			record.SizeStr = humanize.Bytes(uint64(record.Size))
			record.URL, _ = url.JoinPath(s.deliveryURL, record.CDNID)
			return record, nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) DeleteFiles(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "blob.Service.DeleteFiles"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "asset_file",
		IDField: "id",
		Flavor:  sqlbuilder.SQLite,
	})
	return wrap.IfErr(op, err)
}

func (s *Service) SaveFiles(ctx context.Context, headers []*multipart.FileHeader) error {
	const op = "blob.Service.SaveFiles"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	// s.logger.Debug(op, "headers", headers)

	for _, header := range headers {
		file, err := header.Open()
		if err != nil {
			return wrap.IfErr(op, err)
		}
		defer file.Close()

		blob, err := io.ReadAll(file)
		if err != nil {
			return wrap.IfErr(op, err)
		}

		hashSum := sha1.Sum(blob)
		blobHash := hex.EncodeToString(hashSum[:])

		const query = "insert into asset_file (cdnid, blob, hash) values (?, ?, ?);"
		_, err = tx.ExecContext(ctx, query, header.Filename, blob, blobHash)
		if s.store.IsErrConstraintUnique(err) {
			return wrap.IfErr(op, fmt.Errorf("%w: %s", ErrFileExists, header.Filename))
		} else if err != nil {
			return wrap.IfErr(op, err)
		}
	}

	return wrap.IfErr(op, tx.Commit())
}

func (s *Service) ImportFromFolder(ctx context.Context) (*ImportResult, error) {
	return ImportFromFolder(ctx, s.logger, s.store.DB(), "cache")
}

func (s *Service) ExportToFolder(ctx context.Context) (*ExportResult, error) {
	return ExportToFolder(ctx, s.logger, s.store.DB(), "cache", true)
}

func (s *Service) SyncToS3(ctx context.Context, cfg S3Config) (*S3SyncResult, error) {
	return SyncToS3(ctx, s.logger, s.store.DB(), cfg)
}
