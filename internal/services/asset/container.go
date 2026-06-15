package asset

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/dv1x3r/amazing-core/internal/lib/graph"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type Container struct {
	ID        int              `json:"id"`
	OID       w2.Field[string] `json:"oid"`
	OIDStr    string           `json:"oid_str"`
	Name      w2.Field[string] `json:"name"`
	PTag      w2.Field[string] `json:"ptag"`
	Icon      string           `json:"icon"`
	Assets    int              `json:"assets"`
	Packages  int              `json:"packages"`
	CreatedAt w2.UnixTime      `json:"created_at"`
}

type ContainerAsset struct {
	ID       int         `json:"id"`
	Position int         `json:"position"`
	WINAsset w2.Dropdown `json:"win_asset"`
	OSXAsset w2.Dropdown `json:"osx_asset"`
}

type ContainerPackage struct {
	ID           int         `json:"id"`
	Position     int         `json:"position"`
	PkgContainer w2.Dropdown `json:"pkg_container"`
	Assets       int         `json:"assets"`
	Packages     int         `json:"packages"`
}

func (s *Service) GetContainerGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Container], error) {
	const op = "asset.Service.GetContainerGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Container]{
		From: "asset_container as c",
		Select: []string{
			"c.id",
			"c.gsfoid",
			"c.name",
			"c.ptag",
			"coalesce(a.assets, 0) as assets",
			"coalesce(p.packages, 0) as packages",
			"c.created_at",
			"icon.cdnid",
		},
		WhereMapping: map[string]string{
			"id":   "c.id",
			"oid":  "c.gsfoid",
			"name": "c.name",
			"ptag": "c.ptag",
		},
		OrderByMapping: map[string]string{
			"id":         "c.id",
			"oid":        "c.gsfoid",
			"oid_str":    "c.gsfoid",
			"name":       "c.name",
			"ptag":       "c.ptag",
			"assets":     "a.assets",
			"packages":   "p.packages",
			"created_at": "c.created_at",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
					select
						am.container_id,
					    a.cdnid,
					    row_number() over (partition by am.container_id order by am.position) as rn
					from asset_container_assetmap as am
					join asset as a on a.id = am.win_asset_id
					join file_type as ft on ft.id = a.file_type_id
					where ft.name = 'image/png'
				) as icon`,
				"icon.container_id = c.id and icon.rn = 1",
			)
			sb.JoinWithOption(sqlbuilder.LeftJoin,
				"(select container_id, count(*) as assets from asset_container_assetmap group by container_id) as a",
				"a.container_id = c.id",
			)
			sb.JoinWithOption(sqlbuilder.LeftJoin,
				"(select container_id, count(*) as packages from asset_container_package group by container_id) as p",
				"p.container_id = c.id",
			)
		},
		Scan: func(rows *sql.Rows, record *Container) error {
			var icon *string
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.Name,
				&record.PTag,
				&record.Assets,
				&record.Packages,
				&record.CreatedAt,
				&icon,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromString(record.OID.V).String()
			if icon != nil {
				record.Icon, _ = url.JoinPath(s.deliveryURL, *icon)
			}
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetContainerAssetGrid(ctx context.Context, req w2.GetGridRequest, containerID int) (w2.GetGridResponse[ContainerAsset], error) {
	const op = "asset.Service.GetContainerAssetGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[ContainerAsset]{
		From: "asset_container_assetmap as ca",
		Select: []string{
			"ca.id",
			"ca.position",
			"ca.win_asset_id",
			"ca.osx_asset_id",
			"concat_ws(' - ', at.name, a.gsfoid, coalesce(a.res_name, '[NULL]'), (am.metadata ->> '$.assets[0].target_platform') || ' ' || (am.metadata ->> '$.info.version_engine')) as windows",
			"iif(ax.id is null, null, concat_ws(' - ', axt.name, ax.gsfoid, coalesce(ax.res_name, '[NULL]'), (axm.metadata ->> '$.assets[0].target_platform') || ' ' || (axm.metadata ->> '$.info.version_engine'))) as osx",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset as a", "a.id = ca.win_asset_id")
			sb.Join("asset_type as at", "at.id = a.asset_type_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_metadata as am", "am.asset_id = a.id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset as ax", "ax.id = ca.osx_asset_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_type as axt", "axt.id = ax.asset_type_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_metadata as axm", "axm.asset_id = ax.id")
			sb.Where(sb.EQ("ca.container_id", containerID))
			sb.OrderByAsc("ca.position")
		},
		Scan: func(rows *sql.Rows, record *ContainerAsset) error {
			return rows.Scan(
				&record.ID,
				&record.Position,
				&record.WINAsset.ID,
				&record.OSXAsset.ID,
				&record.WINAsset.Text,
				&record.OSXAsset.Text,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetContainerPackageGrid(ctx context.Context, req w2.GetGridRequest, containerID int) (w2.GetGridResponse[ContainerPackage], error) {
	const op = "asset.Service.GetContainerPackageGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[ContainerPackage]{
		From: "asset_container_package as cp",
		Select: []string{
			"cp.id",
			"cp.position",
			"cp.pkg_container_id",
			"concat(pc.gsfoid, ' - ', pc.name, ' (' || pc.ptag || ')') as package",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as pc", "pc.id = cp.pkg_container_id")
			sb.Where(sb.EQ("cp.container_id", containerID))
			sb.OrderByAsc("cp.position")
		},
		Scan: func(rows *sql.Rows, record *ContainerPackage) error {
			return rows.Scan(
				&record.ID,
				&record.Position,
				&record.PkgContainer.ID,
				&record.PkgContainer.Text,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateContainer(ctx context.Context, req w2.SaveFormRequest[Container]) (int, error) {
	const op = "asset.Service.CreateContainer"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "asset_container",
		Cols:   []string{"gsfoid", "name", "ptag"},
		Values: []any{req.Record.OID, req.Record.Name, req.Record.PTag},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrContainerExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateContainer(ctx context.Context, req w2.SaveFormRequest[Container]) error {
	const op = "asset.Service.UpdateContainer"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update:  "asset_container",
		Cols:    []string{"gsfoid", "name", "ptag"},
		Values:  []any{req.Record.OID, req.Record.Name, req.Record.PTag},
		IDField: "id",
		IDValue: req.Record.ID,
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrContainerExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) AddContainerAsset(ctx context.Context, req w2.SaveFormRequest[ContainerAsset], containerID int) (int, error) {
	const op = "asset.Service.AddContainerAsset"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "asset_container_assetmap",
		Cols: []string{"container_id", "win_asset_id", "osx_asset_id", "position"},
		Values: []any{
			containerID,
			req.Record.WINAsset.ID,
			req.Record.OSXAsset.ID,
			sqlbuilder.Buildf("(select coalesce(max(position) + 1, 1) from asset_container_assetmap where container_id = %v)", containerID),
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrContainerAssetExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateContainerAssets(ctx context.Context, req w2.SaveGridRequest[ContainerAsset]) error {
	const op = "asset.Service.UpdateContainerAssets"
	_, err := w2db.SaveGridContext(ctx, s.store.DB(), req, w2db.SaveGridOptions[ContainerAsset]{
		BuildOptions: func(change ContainerAsset) w2db.UpdateOptions {
			return w2db.UpdateOptions{
				Update:  "asset_container_assetmap",
				Cols:    []string{"win_asset_id", "osx_asset_id"},
				Values:  []any{change.WINAsset.ID, change.OSXAsset.ID},
				IDField: "id",
				IDValue: change.ID,
			}
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrContainerAssetExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) AddContainerPackage(ctx context.Context, req w2.SaveFormRequest[ContainerPackage], containerID int) (int, error) {
	const op = "asset.Service.AddContainerPackage"
	var id int
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		var err error
		id, err = w2db.InsertContext(ctx, tx, w2db.InsertOptions{
			Into: "asset_container_package",
			Cols: []string{"container_id", "pkg_container_id", "position"},
			Values: []any{
				containerID,
				req.Record.PkgContainer.ID,
				sqlbuilder.Buildf("(select coalesce(max(position) + 1, 1) from asset_container_package where container_id = %v)", containerID),
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrContainerPackageExists
		} else if err != nil {
			return err
		} else if cycle, err := s.hasPackageCycle(ctx, tx, containerID); err != nil {
			return err
		} else if cycle {
			return ErrPackageCyclicDependency
		}
		return nil
	})
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateContainerPackages(ctx context.Context, req w2.SaveGridRequest[ContainerPackage]) error {
	const op = "asset.Service.UpdateContainerPackages"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.SaveGridContext(ctx, tx, req, w2db.SaveGridOptions[ContainerPackage]{
			BuildOptions: func(change ContainerPackage) w2db.UpdateOptions {
				return w2db.UpdateOptions{
					Update:  "asset_container_package",
					Cols:    []string{"pkg_container_id"},
					Values:  []any{change.PkgContainer.ID},
					IDField: "id",
					IDValue: change.ID,
				}
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrContainerPackageExists
		} else if err != nil {
			return err
		}

		affectedContainerIDs := map[int]struct{}{}
		for _, change := range req.Changes {
			var containerID int
			row := tx.QueryRowContext(ctx, "select container_id from asset_container_package where id = ?;", change.ID)
			if err := row.Scan(&containerID); err != nil {
				return err
			}
			affectedContainerIDs[containerID] = struct{}{}
		}

		s.logger.Debug("affected", "affected containers", affectedContainerIDs)

		for containerID := range affectedContainerIDs {
			if cycle, err := s.hasPackageCycle(ctx, tx, containerID); err != nil {
				return err
			} else if cycle {
				return ErrPackageCyclicDependency
			}
		}

		return nil
	})
	return wrap.IfErr(op, err)
}

func (s *Service) ReorderContainerAssets(ctx context.Context, req w2.ReorderGridRequest) error {
	const op = "asset.Service.ReorderContainerAssets"
	_, err := w2db.ReorderGridContext(ctx, s.store.DB(), req, w2db.ReorderGridOptions{
		Update:     "asset_container_assetmap",
		IDField:    "id",
		SetField:   "position",
		GroupField: "container_id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) ReorderContainerPackages(ctx context.Context, req w2.ReorderGridRequest) error {
	const op = "asset.Service.ReorderContainerPackages"
	_, err := w2db.ReorderGridContext(ctx, s.store.DB(), req, w2db.ReorderGridOptions{
		Update:     "asset_container_package",
		IDField:    "id",
		SetField:   "position",
		GroupField: "container_id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteContainers(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "asset.Service.DeleteContainers"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "asset_container",
		IDField: "id",
	})
	if s.store.IsErrConstraintTrigger(err) {
		return wrap.IfErr(op, ErrContainerInUse)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteContainerAssets(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "asset.Service.DeleteContainerAssets"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "asset_container_assetmap",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteContainerPackages(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "asset.Service.DeleteContainerPackages"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "asset_container_package",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) hasPackageCycle(ctx context.Context, tx *sql.Tx, updatedContainerID int) (bool, error) {
	// build complete adjacency map: { container_id: [pkg_container_id, ...] }
	rows, err := tx.QueryContext(ctx, "select container_id, pkg_container_id from asset_container_package;")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	nodes := graph.AdjacencyMap{}

	for rows.Next() {
		var nodeID, childID int
		if err := rows.Scan(&nodeID, &childID); err != nil {
			return false, err
		}
		nodes.Add(nodeID, childID)
	}

	if err := rows.Err(); err != nil {
		return false, err
	}

	return nodes.HasCycleFrom(updatedContainerID), nil
}

func (s *Service) GetGSFAssetContainer(ctx context.Context, platform gsf.Platform, containerID int) (types.AssetContainer, error) {
	const op = "asset.Service.GetGSFAssetContainer"
	ac, err := s.getGSFAssetContainer(ctx, containerID, platform, map[int]struct{}{})
	return ac, wrap.IfErr(op, err)
}

func (s *Service) getGSFAssetContainer(ctx context.Context, containerID int, platform gsf.Platform, path map[int]struct{}) (types.AssetContainer, error) {
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
			at.name as type_name,
			a.cdnid,
			coalesce(a.res_name, '') as res_name,
			coalesce(ag.name, '') as group_name,
			a.size
		from asset_container_assetmap as ca
		join asset as a on a.id = iif(? = 1 and ca.osx_asset_id is not null, ca.osx_asset_id, ca.win_asset_id)
		join asset_type as at on at.id = a.asset_type_id
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
