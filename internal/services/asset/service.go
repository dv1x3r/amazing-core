package asset

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
	"github.com/dv1x3r/w2go/w2sql"

	"github.com/dustin/go-humanize"
	"github.com/huandu/go-sqlbuilder"
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

type Asset struct {
	ID          int              `json:"id"`
	CDNID       string           `json:"cdnid"`
	URL         string           `json:"url"`
	OID         int              `json:"oid"`
	Class       int              `json:"class"`
	Type        int              `json:"type"`
	Server      int              `json:"server"`
	Number      int              `json:"number"`
	FileType    w2.Dropdown      `json:"file_type"`
	AssetType   w2.Dropdown      `json:"asset_type"`
	AssetGroup  w2.Dropdown      `json:"asset_group"`
	ResName     w2.Field[string] `json:"res_name"`
	Description w2.Field[string] `json:"description"`
	Hash        string           `json:"hash"`
	Size        int              `json:"size"`
	SizeStr     string           `json:"size_str"`
	Metadata    w2.Field[string] `json:"metadata"`
	Version     w2.Field[string] `json:"version"`
}

func (s *Service) GetGridRecords(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Asset], error) {
	const op = "asset.Service.GetGridRecords"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Asset]{
		From: "asset as a",
		Select: []string{
			"a.id",
			"a.cdnid",
			"a.gsfoid",
			"(a.gsfoid >> 56) & 0xFF",
			"(a.gsfoid >> 48) & 0xFF",
			"(a.gsfoid >> 40) & 0xFF",
			"a.gsfoid & 0xFFFFFFFFFF",
			"a.file_type_id",
			"ft.name",
			"a.asset_type_id",
			"at.name",
			"a.asset_group_id",
			"ag.name",
			"a.res_name",
			"a.description",
			"a.hash",
			"a.size",
			"json(am.metadata)",
			"(am.metadata ->> '$.info.version_engine') || ' ' || (am.metadata ->> '$.assets[0].target_platform')",
		},
		WhereMapping: map[string]string{
			"id":          "a.id",
			"cdnid":       "a.cdnid",
			"oid":         "a.gsfoid",
			"class":       "(a.gsfoid >> 56) & 0xFF",
			"type":        "(a.gsfoid >> 48) & 0xFF",
			"server":      "(a.gsfoid >> 40) & 0xFF",
			"number":      "a.gsfoid & 0xFFFFFFFFFF",
			"file_type":   "a.file_type_id",
			"asset_type":  "a.asset_type_id",
			"asset_group": "a.asset_group_id",
			"res_name":    "a.res_name",
			"description": "a.description",
			"hash":        "a.hash",
			"size":        "a.size",
			"metadata":    "json(am.metadata)",
		},
		OrderByMapping: map[string]string{
			"id":          "a.id",
			"cdnid":       "a.cdnid COLLATE BINARY",
			"url":         "a.cdnid COLLATE BINARY",
			"oid":         "a.gsfoid",
			"class":       "(a.gsfoid >> 56) & 0xFF",
			"type":        "(a.gsfoid >> 48) & 0xFF",
			"server":      "(a.gsfoid >> 40) & 0xFF",
			"number":      "a.gsfoid & 0xFFFFFFFFFF",
			"file_type":   "ft.name",
			"asset_type":  "at.name",
			"asset_group": "ag.name",
			"res_name":    "a.res_name",
			"description": "a.description",
			"hash":        "a.hash",
			"size":        "a.size",
			"size_str":    "a.size",
			"version":     "(am.metadata ->> '$.info.version_engine') || ' ' || (am.metadata ->> '$.assets[0].target_platform')",
		},
		Flavor: sqlbuilder.SQLite,
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("file_type as ft", "ft.id = a.file_type_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_type as at", "at.id = a.asset_type_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_group as ag", "ag.id = a.asset_group_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_metadata as am", "am.asset_id = a.id")
		},
		Scan: func(rows *sql.Rows) (Asset, error) {
			var record Asset
			if err := rows.Scan(
				&record.ID,
				&record.CDNID,
				&record.OID,
				&record.Class,
				&record.Type,
				&record.Server,
				&record.Number,
				&record.FileType.ID,
				&record.FileType.Text,
				&record.AssetType.ID,
				&record.AssetType.Text,
				&record.AssetGroup.ID,
				&record.AssetGroup.Text,
				&record.ResName,
				&record.Description,
				&record.Hash,
				&record.Size,
				&record.Metadata,
				&record.Version,
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

func (s *Service) SaveGrid(ctx context.Context, req w2.SaveGridRequest[Asset]) error {
	const op = "asset.Service.SaveGrid"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.SaveGridContext(ctx, tx, req, w2db.SaveGridOptions[Asset]{
			Flavor: sqlbuilder.SQLite,
			BuildUpdate: func(change Asset) *sqlbuilder.UpdateBuilder {
				ub := sqlbuilder.Update("asset")
				w2sql.Set(ub, change.AssetType.ID, "asset_type_id")
				w2sql.Set(ub, change.AssetGroup.ID, "asset_group_id")
				w2sql.Set(ub, change.ResName, "res_name")
				w2sql.Set(ub, change.Description, "description")
				ub.Where(ub.EQ("id", change.ID))
				return ub
			},
		})
		return err
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteByGrid(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "asset.Service.DeleteByGrid"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "asset",
		IDField: "id",
		Flavor:  sqlbuilder.SQLite,
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DumpTable(ctx context.Context) ([]byte, error) {
	const op = "asset.Service.DumpTable"
	data, err := s.store.DumpTable(ctx, "asset")
	return data, wrap.IfErr(op, err)
}

func (s *Service) GetFileTypeDropdownRecords(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "asset.Service.GetFileTypeDropdownRecords"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "file_type",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
		Flavor:       sqlbuilder.SQLite,
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetAssetTypeDropdownRecords(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "asset.Service.GetAssetTypeDropdownRecords"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "asset_type",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
		Flavor:       sqlbuilder.SQLite,
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetAssetGroupDropdownRecords(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "asset.Service.GetAssetGroupDropdownRecords"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "asset_group",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
		Flavor:       sqlbuilder.SQLite,
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetGSFAssetByCDNID(ctx context.Context, cdnid string) (types.Asset, error) {
	const op = "asset.Service.GetGSFAssetByCDNID"
	a := types.Asset{}
	row := s.store.DB().QueryRow(`
			select
				a.gsfoid,
				coalesce(at.name, 'Undefined') as asset_type,
				a.cdnid,
				a.res_name,
				coalesce(ag.name, 'Undefined') as asset_group,
				a.size
			from asset as a
			left join asset_type as at on at.id = a.asset_type_id
			left join asset_group as ag on ag.id = a.asset_group_id
			where a.cdnid = ?;
		`, strings.TrimSpace(cdnid))
	err := row.Scan(&a.OID, &a.AssetTypeName, &a.CDNID, &a.ResName, &a.GroupName, &a.FileSize)
	if err != nil {
		return a, wrap.IfErr(op, fmt.Errorf("cdnid %s: %w", cdnid, err))
	}
	return a, nil
}

func (s *Service) ImportCacheItems(ctx context.Context, items []CacheItem) (*ImportResult, error) {
	return ImportCacheItems(ctx, s.logger, s.store.DB(), items)
}
