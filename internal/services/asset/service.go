package asset

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/graph"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
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

func (s *Service) DeliveryURL() string {
	return s.deliveryURL
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
		order by ca.position asc, ca.id desc;
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
		order by cp.position asc, cp.id desc;
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

func (s *Service) ImportCacheItems(ctx context.Context, items []CacheItem) (*ImportResult, error) {
	return ImportCacheItems(ctx, s.logger, s.store.DB(), items)
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
			"a.description",
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
			"description": "a.description",
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
			"description": "a.description",
			"hash":        "a.hash",
			"size":        "a.size",
			"size_str":    "a.size",
			"version":     "(am.metadata ->> '$.assets[0].target_platform') || ' ' || (am.metadata ->> '$.info.version_engine')",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
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
				&record.Description,
				&record.Hash,
				&record.Size,
				&record.Metadata,
				&record.Version,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromInt64(record.OID).String()
			record.SizeStr = humanize.Bytes(uint64(record.Size))
			record.URL, _ = url.JoinPath(s.deliveryURL, record.CDNID)
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
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
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.Name,
				&record.PTag,
				&record.Assets,
				&record.Packages,
				&record.CreatedAt,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromInt64(record.OID.V).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetContainerAssetGrid(ctx context.Context, req w2.GetGridRequest, id int) (w2.GetGridResponse[ContainerAsset], error) {
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
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Where(sb.EQ("ca.container_id", id))
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset as a", "a.id = ca.win_asset_id")
			sb.Join("asset_type as at", "at.id = a.asset_type_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_metadata as am", "am.asset_id = a.id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset as ax", "ax.id = ca.osx_asset_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_type as axt", "axt.id = ax.asset_type_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_metadata as axm", "axm.asset_id = ax.id")
			sb.OrderByAsc("ca.position").OrderByDesc("ca.id")
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

func (s *Service) GetContainerPackageGrid(ctx context.Context, req w2.GetGridRequest, id int) (w2.GetGridResponse[ContainerPackage], error) {
	const op = "asset.Service.GetContainerPackageGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[ContainerPackage]{
		From: "asset_container_package as cp",
		Select: []string{
			"cp.id",
			"cp.position",
			"cp.pkg_container_id",
			"concat(pc.gsfoid, ' - ', pc.name, ' (' || pc.ptag || ')') as package",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Where(sb.EQ("cp.container_id", id))
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as pc", "pc.id = cp.pkg_container_id")
			sb.OrderByAsc("cp.position").OrderByDesc("cp.id")
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

func (s *Service) GetAssetsDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "asset.Service.GetAssetsDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:    "asset as a",
		IDField: "a.id",
		TextField: `
			concat_ws(' - ',
				at.name,
				a.gsfoid,
				coalesce(a.res_name, '[NULL]'),
				(am.metadata ->> '$.assets[0].target_platform') || ' ' || (am.metadata ->> '$.info.version_engine')
			)`,
		OrderByField: "at.name, a.gsfoid",
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_type as at", "at.id = a.asset_type_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_metadata as am", "am.asset_id = a.id")
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetContainersDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "asset.Service.GetContainersDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "asset_container",
		IDField:      "id",
		TextField:    "concat(gsfoid, ' - ', name, ' (' || ptag || ')')",
		OrderByField: "gsfoid",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetFileTypesDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "asset.Service.GetFileTypesDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "file_type",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetAssetTypesDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "asset.Service.GetAssetTypesDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "asset_type",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetAssetGroupsDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "asset.Service.GetAssetGroupsDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "asset_group",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "name",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) UpdateAssets(ctx context.Context, req w2.SaveGridRequest[Asset]) error {
	const op = "asset.Service.UpdateAssets"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.SaveGridContext(ctx, tx, req, w2db.SaveGridOptions[Asset]{
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

func (s *Service) CreateContainer(ctx context.Context, req w2.SaveFormRequest[Container]) (int, error) {
	const op = "asset.Service.CreateContainer"
	id, err := w2db.InsertFormContext(ctx, s.store.DB(), req, w2db.InsertFormOptions{
		Into: "asset_container",
		Cols: []string{"gsfoid", "name", "ptag"},
		Values: []any{
			sqlbuilder.Raw("(select coalesce(max(gsfoid) + 1, 1) from asset_container)"),
			req.Record.Name,
			req.Record.PTag,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrContainerExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateContainers(ctx context.Context, req w2.SaveGridRequest[Container]) error {
	const op = "asset.Service.UpdateContainers"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.SaveGridContext(ctx, tx, req, w2db.SaveGridOptions[Container]{
			BuildUpdate: func(change Container) *sqlbuilder.UpdateBuilder {
				ub := sqlbuilder.Update("asset_container")
				w2sql.Set(ub, change.OID, "gsfoid")
				w2sql.Set(ub, change.Name, "name")
				w2sql.Set(ub, change.PTag, "ptag")
				ub.Where(ub.EQ("id", change.ID))
				return ub
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrContainerExists
		}
		return err
	})
	return wrap.IfErr(op, err)
}

func (s *Service) AddContainerAsset(ctx context.Context, req w2.SaveFormRequest[ContainerAsset], containerID int) error {
	const op = "asset.Service.AddContainerAsset"
	_, err := w2db.InsertFormContext(ctx, s.store.DB(), req, w2db.InsertFormOptions{
		Into:   "asset_container_assetmap",
		Cols:   []string{"container_id", "win_asset_id", "osx_asset_id"},
		Values: []any{containerID, req.Record.WINAsset.ID, req.Record.OSXAsset.ID},
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrContainerAssetExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) UpdateContainerAssets(ctx context.Context, req w2.SaveGridRequest[ContainerAsset]) error {
	const op = "asset.Service.UpdateContainerAssets"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.SaveGridContext(ctx, tx, req, w2db.SaveGridOptions[ContainerAsset]{
			BuildUpdate: func(change ContainerAsset) *sqlbuilder.UpdateBuilder {
				ub := sqlbuilder.Update("asset_container_assetmap")
				w2sql.Set(ub, change.WINAsset.ID, "win_asset_id")
				w2sql.Set(ub, change.OSXAsset.ID, "osx_asset_id")
				ub.Where(ub.EQ("id", change.ID))
				return ub
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrContainerAssetExists
		}
		return err
	})
	return wrap.IfErr(op, err)
}

func (s *Service) AddContainerPackage(ctx context.Context, req w2.SaveFormRequest[ContainerPackage], containerID int) error {
	const op = "asset.Service.AddContainerPackage"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.InsertFormContext(ctx, tx, req, w2db.InsertFormOptions{
			Into:   "asset_container_package",
			Cols:   []string{"container_id", "pkg_container_id"},
			Values: []any{containerID, req.Record.PkgContainer.ID},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrContainerPackageExists
		} else if err != nil {
			return err
		}
		if cycle, err := hasPackageCycle(ctx, tx, containerID); err != nil {
			return err
		} else if cycle {
			return ErrPackageCyclicDependency
		}
		return nil
	})
	return wrap.IfErr(op, err)
}

func (s *Service) UpdateContainerPackages(ctx context.Context, req w2.SaveGridRequest[ContainerPackage]) error {
	const op = "asset.Service.UpdateContainerPackages"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.SaveGridContext(ctx, tx, req, w2db.SaveGridOptions[ContainerPackage]{
			BuildUpdate: func(change ContainerPackage) *sqlbuilder.UpdateBuilder {
				ub := sqlbuilder.Update("asset_container_package")
				w2sql.Set(ub, change.PkgContainer.ID, "pkg_container_id")
				ub.Where(ub.EQ("id", change.ID))
				return ub
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
			if cycle, err := hasPackageCycle(ctx, tx, containerID); err != nil {
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
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.ReorderGridContext(ctx, tx, req, w2db.ReorderGridOptions{
			Update:     "asset_container_assetmap",
			IDField:    "id",
			SetField:   "position",
			GroupField: "container_id",
		})
		return err
	})
	return wrap.IfErr(op, err)
}

func (s *Service) ReorderContainerPackages(ctx context.Context, req w2.ReorderGridRequest) error {
	const op = "asset.Service.ReorderContainerPackages"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.ReorderGridContext(ctx, tx, req, w2db.ReorderGridOptions{
			Update:     "asset_container_package",
			IDField:    "id",
			SetField:   "position",
			GroupField: "container_id",
		})
		return err
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

func hasPackageCycle(ctx context.Context, tx *sql.Tx, updatedContainerID int) (bool, error) {
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
