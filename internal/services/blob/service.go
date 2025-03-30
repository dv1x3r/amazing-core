package blob

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"

	"github.com/dustin/go-humanize"
	"github.com/dv1x3r/amazing-core/internal/config"
	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/logger"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"

	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2sql/w2sqlbuilder"
	"github.com/huandu/go-sqlbuilder"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFileExists   = errors.New("file with the same name already exists")
)

type FileInfo struct {
	ID      int    `json:"id"`
	CDNID   string `json:"cdnid"`
	Hash    string `json:"hash"`
	Size    int    `json:"size"`
	SizeStr string `json:"size_str"`
	URL     string `json:"url"`
}

func (i *FileInfo) ScanRow(scan func(dest ...any) error) error {
	return scan(&i.ID, &i.CDNID, &i.Hash, &i.Size)
}

type Service struct {
	store db.Store
}

func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) FetchFileBlob(ctx context.Context, cdnid string) ([]byte, error) {
	const op = "blob.Service.FetchFileBlob"
	const query = "select blob from blob.asset_file where cdnid = ?;"
	var blob []byte
	err := s.store.DB().QueryRowContext(ctx, query, cdnid).Scan(&blob)
	if err == sql.ErrNoRows {
		return nil, ErrFileNotFound
	}
	return blob, wrap.IfErr(op, err)
}

func (s *Service) FetchFilesList(ctx context.Context, r w2.GridDataRequest) ([]FileInfo, int, error) {
	const op = "blob.Service.FetchFilesList"

	var total int
	var records []FileInfo

	sb := sqlbuilder.Select("count(*)")
	sb.From("blob.asset_file")

	w2sqlbuilder.Where(sb, r, map[string]string{
		"cdnid": "cdnid",
		"hash":  "hash",
	})

	query, args := sb.BuildWithFlavor(sqlbuilder.SQLite)
	row := s.store.DB().QueryRowContext(ctx, query, args...)
	if err := row.Scan(&total); err != nil && err != sql.ErrNoRows {
		return nil, 0, wrap.IfErr(op, err)
	}

	sb.Select(
		"id",
		"cdnid",
		"hash",
		"length(blob) as size",
	)

	w2sqlbuilder.OrderBy(sb, r, map[string]string{
		"id":       "id",
		"cdnid":    "cdnid",
		"url":      "cdnid",
		"hash":     "hash",
		"size":     "length(blob)",
		"size_str": "length(blob)",
	})

	w2sqlbuilder.Limit(sb, r)
	w2sqlbuilder.Offset(sb, r)

	query, args = sb.BuildWithFlavor(sqlbuilder.SQLite)
	rows, err := s.store.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, wrap.IfErr(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var info FileInfo
		if err := info.ScanRow(rows.Scan); err != nil {
			return nil, 0, wrap.IfErr(op, err)
		}
		records = append(records, info)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, wrap.IfErr(op, err)
	}

	assetDeliveryURL := config.Get().Settings.AssetDeliveryURL
	for i := range records {
		records[i].SizeStr = humanize.Bytes(uint64(records[i].Size))
		records[i].URL, err = url.JoinPath(assetDeliveryURL, records[i].CDNID)
		if err != nil {
			return nil, 0, wrap.IfErr(op, err)
		}
	}

	return records, total, nil
}

func (s *Service) SaveFiles(ctx context.Context, files []*multipart.FileHeader) error {
	const op = "blob.Service.SaveFiles"

	logger.Get().Debug(op, "files", files)

	tx, err := s.store.DB().Begin()
	if err != nil {
		return wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	for _, header := range files {
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

		const query = "insert into blob.asset_file (cdnid, blob, hash) values (?, ?, ?);"
		_, err = tx.ExecContext(ctx, query, header.Filename, blob, blobHash)
		if s.store.IsErrConstraintUnique(err) {
			return fmt.Errorf("%w: %s", ErrFileExists, header.Filename)
		} else if err != nil {
			return wrap.IfErr(op, err)
		}
	}

	return wrap.IfErr(op, tx.Commit())
}

func (s *Service) DeleteFiles(ctx context.Context, ids []int) error {
	const op = "blob.Service.DeleteFiles"
	dlb := sqlbuilder.DeleteFrom("blob.asset_file")
	dlb.Where(dlb.In("id", sqlbuilder.List(ids)))
	query, args := dlb.BuildWithFlavor(sqlbuilder.SQLite)
	_, err := s.store.DB().ExecContext(ctx, query, args...)
	return wrap.IfErr(op, err)
}
