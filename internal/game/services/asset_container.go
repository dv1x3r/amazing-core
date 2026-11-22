package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2db"
)

func (s *Service) GetAssetContainer(ctx context.Context, platform gsf.Platform, containerID int) (types.AssetContainer, error) {
	const op = "game.Service.GetAssetContainer"
	ac, err := s.getAssetContainer(ctx, containerID, platform, map[int]struct{}{})
	return ac, wrap.IfErr(op, err)
}

func (s *Service) GetAssetContainersByOIDs(ctx context.Context, platform gsf.Platform, oids []types.OID) ([]types.AssetContainer, error) {
	const op = "game.Service.GetAssetContainersByOIDs"
	containers := make([]types.AssetContainer, 0, len(oids))
	for _, oid := range oids {
		var containerID int
		err := s.store.DB().QueryRowContext(ctx, "select id from asset_container where gsfoid = ?;", oid).Scan(&containerID)
		if errors.Is(err, sql.ErrNoRows) {
			containers = append(containers, types.AssetContainer{
				OID:           oid,
				AssetMap:      types.AssetMap{},
				AssetPackages: []types.AssetPackage{},
			})
			continue
		}
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		container, err := s.GetAssetContainer(ctx, platform, containerID)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		containers = append(containers, container)
	}
	return containers, nil
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

	assets, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[types.Asset]{
		Query: `
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
		`,
		Args: []any{useOSXAsset, containerID},
		Scan: func(rows *sql.Rows, asset *types.Asset) error {
			return rows.Scan(
				&asset.OID,
				&asset.CDNID,
				&asset.ResName,
				&asset.AssetTypeName,
				&asset.GroupName,
				&asset.FileSize,
			)
		},
	})
	if err != nil {
		return ac, err
	}
	for _, asset := range assets {
		ac.AssetMap[asset.AssetTypeName] = append(ac.AssetMap[asset.AssetTypeName], asset)
	}

	// An asset package is a container too, its packages are read after the rows are closed.
	type pkgRow struct {
		containerID int
		pkg         types.AssetPackage
	}
	pkgRows, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[pkgRow]{
		Query: `
			select
				cp.pkg_container_id,
				coalesce(c.ptag, '') as ptag,
				c.created_at
			from asset_container_package as cp
			join asset_container as c on c.id = cp.pkg_container_id
			where cp.container_id = ?
			order by cp.position;
		`,
		Args: []any{containerID},
		Scan: func(rows *sql.Rows, r *pkgRow) error {
			return rows.Scan(&r.containerID, &r.pkg.PTag, &r.pkg.CreatedDate)
		},
	})
	if err != nil {
		return ac, err
	}
	for _, r := range pkgRows {
		pkgContainer, err := s.getAssetContainer(ctx, r.containerID, platform, path)
		if err != nil {
			return ac, err
		}
		pkg := r.pkg
		pkg.AssetContainer = pkgContainer
		ac.AssetPackages = append(ac.AssetPackages, pkg)
	}

	return ac, nil
}
