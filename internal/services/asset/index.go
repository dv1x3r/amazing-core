package asset

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
)

type AssetIndex struct {
	Assets []AssetIndexEntry `json:"assets"`
}

type AssetIndexEntry struct {
	ID            int    `json:"id"`
	CDNID         string `json:"cdnid"`
	OID           string `json:"oid"`
	FileType      string `json:"file_type"`
	AssetType     string `json:"asset_type,omitempty"`
	AssetGroup    string `json:"asset_group,omitempty"`
	ResName       string `json:"res_name,omitempty"`
	Hash          string `json:"hash"`
	Size          int64  `json:"size"`
	BundleVersion string `json:"bundle_version,omitempty"`
	Scene         bool   `json:"scene"`
}

func (s *Service) GetAssetIndex(ctx context.Context) (AssetIndex, error) {
	const op = "asset.Service.GetAssetIndex"

	const query = `
		select
			a.id,
			a.cdnid,
			a.gsfoid,
			ft.name,
			coalesce(at.name, ''),
			coalesce(ag.name, ''),
			coalesce(a.res_name, ''),
			a.hash,
			a.size,
			coalesce(a.bundle_version, ''),
			iif(
				json_array_length(af.metadata, '$.bundle.scene') > 0 and
				af.metadata ->> '$.bundle.counts.types.Mesh' > 0,
				1, 0)
		from asset as a
		join file_type as ft on ft.id = a.file_type_id
		left join asset_type as at on at.id = a.asset_type_id
		left join asset_group as ag on ag.id = a.asset_group_id
		left join asset_file as af on af.cdnid = a.cdnid
		order by a.id;
	`

	rows, err := s.store.DB().QueryContext(ctx, query)
	if err != nil {
		return AssetIndex{}, wrap.IfErr(op, err)
	}
	defer rows.Close()

	index := AssetIndex{Assets: []AssetIndexEntry{}}
	for rows.Next() {
		entry := AssetIndexEntry{}
		if err := rows.Scan(
			&entry.ID,
			&entry.CDNID,
			&entry.OID,
			&entry.FileType,
			&entry.AssetType,
			&entry.AssetGroup,
			&entry.ResName,
			&entry.Hash,
			&entry.Size,
			&entry.BundleVersion,
			&entry.Scene,
		); err != nil {
			return AssetIndex{}, wrap.IfErr(op, err)
		}
		index.Assets = append(index.Assets, entry)
	}

	if err := rows.Err(); err != nil {
		return AssetIndex{}, wrap.IfErr(op, err)
	}

	return index, nil
}
