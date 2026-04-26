package asset

import (
	"context"
	"fmt"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

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

func (s *Service) GetGSFAssetContainer(ctx context.Context, platform gsf.Platform, id int) (types.AssetContainer, error) {
	const op = "asset.Service.GetGSFAssetContainer"
	ac, err := s.getGSFAssetContainer(ctx, id, platform, map[int]struct{}{})
	return ac, wrap.IfErr(op, err)
}

func (s *Service) getGSFAssetContainer(ctx context.Context, id int, platform gsf.Platform, path map[int]struct{}) (types.AssetContainer, error) {
	ac := types.AssetContainer{
		AssetMap:      types.AssetMap{},
		AssetPackages: []types.AssetPackage{},
	}

	if _, ok := path[id]; ok {
		return ac, fmt.Errorf("circular dependency detected for container %d", id)
	}

	path[id] = struct{}{}
	defer delete(path, id)

	row := s.store.DB().QueryRowContext(ctx, "select gsfoid from asset_container where id = ?;", id)
	if err := row.Scan(&ac.OID); err != nil {
		return ac, err
	}

	useOSXAsset := 0
	if platform == gsf.PlatformOSX {
		useOSXAsset = 1
	}

	rows, err := s.store.DB().QueryContext(ctx, `
		select
			a.gsfoid,
			at.name as type_name,
			a.cdnid,
			coalesce(a.res_name, 'Undefined') as res_name,
			coalesce(ag.name, 'Undefined') as group_name,
			a.size
		from asset_container_assetmap as ca
		join asset as a on a.id = iif(? = 1 and ca.osx_asset_id is not null, ca.osx_asset_id, ca.win_asset_id)
		join asset_type as at on at.id = a.asset_type_id
		left join asset_group as ag on ag.id = a.asset_group_id
		where ca.container_id = ?
		order by ca.position;
	`, useOSXAsset, id)
	if err != nil {
		return ac, err
	}
	defer rows.Close()

	for rows.Next() {
		var asset types.Asset
		if err := rows.Scan(
			&asset.OID,
			&asset.AssetTypeName,
			&asset.CDNID,
			&asset.ResName,
			&asset.GroupName,
			&asset.FileSize,
		); err != nil {
			return ac, err
		}

		ac.AssetMap[asset.AssetTypeName] = append(ac.AssetMap[asset.AssetTypeName], asset)
	}
	if err := rows.Err(); err != nil {
		return ac, err
	}

	pkgRows, err := s.store.DB().QueryContext(ctx, `
		select
			cp.pkg_container_id,
			coalesce(c.ptag, '') as ptag,
			c.created_at
		from asset_container_package as cp
		join asset_container as c on c.id = cp.pkg_container_id
		where cp.container_id = ?
		order by cp.position;
	`, id)
	if err != nil {
		return ac, err
	}
	defer pkgRows.Close()

	for pkgRows.Next() {
		var pkgContainerID int
		var pkg types.AssetPackage
		if err := pkgRows.Scan(&pkgContainerID, &pkg.PTag, &pkg.CreatedDate); err != nil {
			return ac, err
		}
		pkgContainer, err := s.getGSFAssetContainer(ctx, pkgContainerID, platform, path)
		if err != nil {
			return ac, err
		}
		pkg.AssetContainer = pkgContainer
		ac.AssetPackages = append(ac.AssetPackages, pkg)
	}
	if err := pkgRows.Err(); err != nil {
		return ac, err
	}

	return ac, nil
}
