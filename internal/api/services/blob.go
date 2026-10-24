package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/dustin/go-humanize"
)

var ErrBlobFileNotFound = errors.New("file is not found")

type File struct {
	ID       int              `json:"id"`
	CDNID    string           `json:"cdnid"`
	Hash     string           `json:"hash"`
	Size     int              `json:"size"`
	SizeStr  string           `json:"size_str"`
	Metadata w2.Field[string] `json:"metadata"`
}

func (s *Service) GetBlobData(ctx context.Context, cdnid string) ([]byte, error) {
	const op = "api.Service.GetBlobData"

	var data []byte

	const query = "select blob from blob.asset_file where cdnid = ?;"
	row := s.store.DB().QueryRowContext(ctx, query, cdnid)
	if err := row.Scan(&data); err == sql.ErrNoRows {
		return nil, ErrBlobFileNotFound
	} else if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	return data, nil
}

func (s *Service) GetBlobMetadata(ctx context.Context, cdnid string) (string, error) {
	const op = "api.Service.GetBlobMetadata"

	var metadata sql.Null[string]

	const query = "select json_pretty(metadata) from blob.asset_file where cdnid = ?;"
	row := s.store.DB().QueryRowContext(ctx, query, cdnid)
	if err := row.Scan(&metadata); err == sql.ErrNoRows {
		return "", ErrBlobFileNotFound
	} else if err != nil {
		return "", wrap.IfErr(op, err)
	}

	return metadata.V, nil
}

func (s *Service) GetBlobGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[File], error) {
	const op = "api.Service.GetBlobGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[File]{
		From: "blob.asset_file",
		Select: []string{
			"id",
			"cdnid",
			"hash",
			"length(blob) as size",
			"json_pretty(metadata) as metadata",
		},
		Where: map[string]string{
			"id":    "id",
			"cdnid": "cdnid",
			"hash":  "hash",
		},
		OrderBy: map[string]string{
			"id":       "id",
			"cdnid":    "cdnid COLLATE BINARY",
			"size_str": "length(blob)",
		},
		Scan: func(rows *sql.Rows, record *File) error {
			if err := rows.Scan(
				&record.ID,
				&record.CDNID,
				&record.Hash,
				&record.Size,
				&record.Metadata,
			); err != nil {
				return err
			}
			record.SizeStr = humanize.Bytes(uint64(record.Size))
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) DeleteFiles(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteFiles"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "blob.asset_file",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
