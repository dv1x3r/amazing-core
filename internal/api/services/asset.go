package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/dustin/go-humanize"
	"github.com/huandu/go-sqlbuilder"
)

type Asset struct {
	ID            int              `json:"id"`
	OID           string           `json:"oid"`
	OIDStr        string           `json:"oid_str"`
	CDNID         string           `json:"cdnid"`
	FileType      w2.Dropdown      `json:"file_type"`
	AssetType     w2.Dropdown      `json:"asset_type"`
	AssetGroup    w2.Dropdown      `json:"asset_group"`
	ResName       w2.Field[string] `json:"res_name"`
	Hash          string           `json:"hash"`
	Size          int              `json:"size"`
	SizeStr       string           `json:"size_str"`
	BundleVersion w2.Field[string] `json:"bundle_version"`
	Icon          string           `json:"icon"`
	Containers    []w2.Dropdown    `json:"containers"`
}

func (s *Service) GetAssetGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Asset], error) {
	const op = "api.Service.GetAssetGrid"
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
			"a.bundle_version",
			"cnt.containers",
		},
		Where: map[string]string{
			"id":             "a.id",
			"oid":            "a.gsfoid",
			"cdnid":          "a.cdnid",
			"file_type":      "a.file_type_id",
			"asset_type":     "a.asset_type_id",
			"asset_group":    "a.asset_group_id",
			"res_name":       "a.res_name",
			"hash":           "a.hash",
			"bundle_version": "a.bundle_version",
		},
		OrderBy: map[string]string{
			"id":             "a.id",
			"oid":            "a.gsfoid",
			"oid_str":        "a.gsfoid",
			"cdnid":          "a.cdnid COLLATE BINARY",
			"file_type":      "ft.name",
			"asset_type":     "at.name",
			"asset_group":    "ag.name",
			"res_name":       "a.res_name",
			"size_str":       "a.size",
			"bundle_version": "a.bundle_version",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("file_type as ft", "ft.id = a.file_type_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_type as at", "at.id = a.asset_type_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_group as ag", "ag.id = a.asset_group_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
					select
						am.asset_id,
						json_group_array(
							json_object('id', c.id, 'text', c.name)
							order by c.name
						) as containers
					from (
						select container_id, win_asset_id as asset_id from asset_container_assetmap
						union
						select container_id, osx_asset_id as asset_id from asset_container_assetmap where osx_asset_id is not null
					) as am
					join asset_container as c on c.id = am.container_id
					group by am.asset_id
				) as cnt`, "cnt.asset_id = a.id")
		},
		Scan: func(rows *sql.Rows, record *Asset) error {
			var containers sql.Null[string]
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
				&record.BundleVersion,
				&containers,
			); err != nil {
				return err
			}
			if containers.Valid {
				if err := json.Unmarshal([]byte(containers.V), &record.Containers); err != nil {
					return fmt.Errorf("containers unmarshal: %w", err)
				}
			}
			record.OIDStr = types.OIDFromString(record.OID).String()
			record.SizeStr = humanize.Bytes(uint64(record.Size))
			if record.FileType.Text.V == "image/png" {
				record.Icon, _ = url.JoinPath(s.deliveryURL, record.CDNID)
			}
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) UpdateAssets(ctx context.Context, req w2.SaveGridRequest[Asset]) error {
	const op = "api.Service.UpdateAssets"
	_, err := w2db.SaveGridContext(ctx, s.store.DB(), req, w2db.SaveGridOptions[Asset]{
		BuildOptions: func(change Asset) w2db.UpdateOptions {
			return w2db.UpdateOptions{
				Update: "asset",
				Values: map[string]any{
					"asset_type_id":  change.AssetType.ID,
					"asset_group_id": change.AssetGroup.ID,
					"res_name":       change.ResName,
				},
				Where: map[string]any{
					"id": change.ID,
				},
			}
		},
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteAssets(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteAssets"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "asset",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
