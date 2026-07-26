package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetAssetByCDNID(ctx context.Context, cdnid string) (types.Asset, error) {
	const op = "game.Service.GetAssetByCDNID"
	a := types.Asset{}
	row := s.store.DB().QueryRowContext(ctx, `
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
		`, strings.TrimSpace(cdnid))
	if err := row.Scan(&a.OID, &a.CDNID, &a.ResName, &a.AssetTypeName, &a.GroupName, &a.FileSize); err != nil {
		return a, wrap.IfErr(op, fmt.Errorf("cdnid %s: %w", cdnid, err))
	}
	return a, nil
}
