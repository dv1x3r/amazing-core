package services

import (
	"context"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetAssetContainer(ctx context.Context, platform gsf.Platform, containerID int) (types.AssetContainer, error) {
	const op = "game.Service.GetAssetContainer"
	ac, err := s.getAssetContainer(ctx, containerID, platform, map[int]struct{}{})
	return ac, wrap.IfErr(op, err)
}

func (s *Service) getAssetContainer(ctx context.Context, containerID int, platform gsf.Platform, path map[int]struct{}) (types.AssetContainer, error) {
	ac := types.AssetContainer{
		AssetMap:      types.AssetMap{},
		AssetPackages: []types.AssetPackage{},
	}

	if _, ok := path[containerID]; ok {
		return ac, fmt.Errorf("circular dependency detected for container %d", containerID)
	}

	path[containerID] = struct{}{}
	defer delete(path, containerID)

	row := s.store.DB().QueryRowContext(ctx, "select gsfoid from asset_container where id = ?;", containerID)
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
			a.cdnid,
			coalesce(a.res_name, '') as res_name,
			coalesce(at.name, '') as type_name,
			coalesce(ag.name, '') as group_name,
			a.size
		from asset_container_assetmap as ca
		join asset as a on a.id = iif(? = 1 and ca.osx_asset_id is not null, ca.osx_asset_id, ca.win_asset_id)
		left join asset_type as at on at.id = a.asset_type_id
		left join asset_group as ag on ag.id = a.asset_group_id
		where ca.container_id = ?
		order by ca.position;
	`, useOSXAsset, containerID)
	if err != nil {
		return ac, err
	}
	defer rows.Close()

	for rows.Next() {
		var asset types.Asset
		if err := rows.Scan(
			&asset.OID,
			&asset.CDNID,
			&asset.ResName,
			&asset.AssetTypeName,
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
	`, containerID)
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
		pkgContainer, err := s.getAssetContainer(ctx, pkgContainerID, platform, path)
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
