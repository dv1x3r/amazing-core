package asset

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type Asset struct {
	ID         int              `json:"id"`
	OID        string           `json:"oid"`
	OIDStr     string           `json:"oid_str"`
	CDNID      string           `json:"cdnid"`
	URL        string           `json:"url"`
	FileType   w2.Dropdown      `json:"file_type"`
	AssetType  w2.Dropdown      `json:"asset_type"`
	AssetGroup w2.Dropdown      `json:"asset_group"`
	ResName    w2.Field[string] `json:"res_name"`
	Hash       string           `json:"hash"`
	Size       int              `json:"size"`
	SizeStr    string           `json:"size_str"`
	Metadata   w2.Field[string] `json:"metadata"`
	Version    w2.Field[string] `json:"version"`
	Icon       string           `json:"icon"`
}

func (s *Service) GetAssetGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Asset], error) {
	const op = "asset.Service.GetAssetGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Asset]{
		From: "asset as a",
		Select: []string{
			"a.id",
			"a.gsfoid",
			"a.cdnid",
			"a.file_type_id",
			"ft.name as file_type_name",
			"a.asset_type_id",
			"at.name as type_name",
			"a.asset_group_id",
			"ag.name as group_name",
			"a.res_name",
			"a.hash",
			"a.size",
			"json(am.metadata) as metadata",
			"(am.metadata ->> '$.assets[0].target_platform') || ' ' || (am.metadata ->> '$.info.version_engine') as platform",
		},
		WhereMapping: map[string]string{
			"id":          "a.id",
			"oid":         "a.gsfoid",
			"cdnid":       "a.cdnid",
			"file_type":   "a.file_type_id",
			"asset_type":  "a.asset_type_id",
			"asset_group": "a.asset_group_id",
			"res_name":    "a.res_name",
			"hash":        "a.hash",
			"size":        "a.size",
			"version":     "(am.metadata ->> '$.assets[0].target_platform') || ' ' || (am.metadata ->> '$.info.version_engine')",
			"metadata":    "json(am.metadata)",
		},
		OrderByMapping: map[string]string{
			"id":          "a.id",
			"oid":         "a.gsfoid",
			"oid_str":     "a.gsfoid",
			"cdnid":       "a.cdnid COLLATE BINARY",
			"file_type":   "ft.name",
			"asset_type":  "at.name",
			"asset_group": "ag.name",
			"res_name":    "a.res_name",
			"hash":        "a.hash",
			"size":        "a.size",
			"size_str":    "a.size",
			"version":     "(am.metadata ->> '$.assets[0].target_platform') || ' ' || (am.metadata ->> '$.info.version_engine')",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("file_type as ft", "ft.id = a.file_type_id")
			sb.Join("asset_type as at", "at.id = a.asset_type_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_group as ag", "ag.id = a.asset_group_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_metadata as am", "am.asset_id = a.id")
		},
		Scan: func(rows *sql.Rows, record *Asset) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.CDNID,
				&record.FileType.ID,
				&record.FileType.Text,
				&record.AssetType.ID,
				&record.AssetType.Text,
				&record.AssetGroup.ID,
				&record.AssetGroup.Text,
				&record.ResName,
				&record.Hash,
				&record.Size,
				&record.Metadata,
				&record.Version,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromString(record.OID).String()
			record.SizeStr = humanize.Bytes(uint64(record.Size))
			record.URL, _ = url.JoinPath(s.deliveryURL, record.CDNID)
			if record.FileType.Text.V == "image/png" {
				record.Icon, _ = url.JoinPath(s.deliveryURL, record.CDNID)
			}
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) UpdateAssets(ctx context.Context, req w2.SaveGridRequest[Asset]) error {
	const op = "asset.Service.UpdateAssets"
	_, err := w2db.SaveGridContext(ctx, s.store.DB(), req, w2db.SaveGridOptions[Asset]{
		BuildOptions: func(change Asset) w2db.UpdateOptions {
			return w2db.UpdateOptions{
				Update:  "asset",
				Cols:    []string{"asset_type_id", "asset_group_id", "res_name"},
				Values:  []any{change.AssetType.ID, change.AssetGroup.ID, change.ResName},
				IDField: "id",
				IDValue: change.ID,
			}
		},
	})
	return wrap.IfErr(op, err)
}

func (s *Service) GetGSFAssetByCDNID(ctx context.Context, cdnid string) (types.Asset, error) {
	const op = "asset.Service.GetGSFAssetByCDNID"
	a := types.Asset{}
	row := s.store.DB().QueryRowContext(ctx, `
			select
				a.gsfoid,
				at.name as type_name,
				a.cdnid,
				coalesce(a.res_name, 'Undefined') as res_name,
				coalesce(ag.name, 'Undefined') as group_name,
				a.size
			from asset as a
			join asset_type as at on at.id = a.asset_type_id
			left join asset_group as ag on ag.id = a.asset_group_id
			where a.cdnid = ?;
		`, strings.TrimSpace(cdnid))
	if err := row.Scan(&a.OID, &a.AssetTypeName, &a.CDNID, &a.ResName, &a.GroupName, &a.FileSize); err != nil {
		return a, wrap.IfErr(op, fmt.Errorf("cdnid %s: %w", cdnid, err))
	}
	return a, nil
}
