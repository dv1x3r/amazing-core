package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2db"
)

func (s *Service) GetAssetByCDNID(ctx context.Context, cdnid string) (types.Asset, error) {
	const op = "game.Service.GetAssetByCDNID"
	asset, found, err := w2db.SelectRowContext(ctx, s.store.DB(), w2db.SelectRowOptions[types.Asset]{
		Query: `
			select
				a.gsfoid,
				a.cdnid,
				coalesce(a.res_name, '') as res_name,
				coalesce(at.name, '') as type_name,
				coalesce(ag.name, '') as group_name,
				a.size
			from asset as a
			left join asset_type as at on at.id = a.asset_type_id
			left join asset_group as ag on ag.id = a.asset_group_id
			where a.cdnid = ?;
			`,
		Args: []any{cdnid},
		Scan: func(row *sql.Row, record *types.Asset) error {
			return row.Scan(
				&record.OID,
				&record.CDNID,
				&record.ResName,
				&record.AssetTypeName,
				&record.GroupName,
				&record.FileSize,
			)
		},
	})
	if err != nil {
		return asset, wrap.IfErr(op, fmt.Errorf("cdnid %s: %w", cdnid, err))
	} else if !found {
		return asset, wrap.IfErr(op, fmt.Errorf("cdnid %s: not found", cdnid))
	}
	return asset, nil
}
