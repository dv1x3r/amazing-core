package asset

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
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

func (s *Service) ImportCacheItems(ctx context.Context, items []CacheItem) (*ImportResult, error) {
	return ImportCacheItems(ctx, s.logger, s.store.DB(), items)
}

func (s *Service) GetGSFAssetByCDNID(ctx context.Context, cdnid string) (types.Asset, error) {
	const op = "asset.Service.GetGSFAssetByCDNID"
	a := types.Asset{}
	row := s.store.DB().QueryRowContext(ctx, `
			select
				a.gsfoid,
				coalesce(at.name, 'Undefined') as asset_type,
				a.cdnid,
				coalesce(a.res_name, 'Undefined') as res_name,
				coalesce(ag.name, 'Undefined') as asset_group,
				a.size
			from asset as a
			left join asset_type as at on at.id = a.asset_type_id
			left join asset_group as ag on ag.id = a.asset_group_id
			where a.cdnid = ?;
		`, strings.TrimSpace(cdnid))
	err := row.Scan(&a.OID, &a.AssetTypeName, &a.CDNID, &a.ResName, &a.GroupName, &a.FileSize)
	if err != nil {
		return a, wrap.IfErr(op, fmt.Errorf("cdnid %s: %w", cdnid, err))
	}
	return a, nil
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
			"ft.name",
			"a.asset_type_id",
			"at.name",
			"a.asset_group_id",
			"ag.name",
			"a.res_name",
			"a.description",
			"a.hash",
			"a.size",
			"json(am.metadata)",
			"(am.metadata ->> '$.info.version_engine') || ' ' || (am.metadata ->> '$.assets[0].target_platform')",
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
			"metadata":    "json(am.metadata)",
		},
		OrderByMapping: map[string]string{
			"id":          "a.id",
			"oid":         "a.gsfoid",
			"cdnid":       "a.cdnid COLLATE BINARY",
			"file_type":   "ft.name",
			"asset_type":  "at.name",
			"asset_group": "ag.name",
			"res_name":    "a.res_name",
			"description": "a.description",
			"hash":        "a.hash",
			"size":        "a.size",
			"size_str":    "a.size",
			"version":     "(am.metadata ->> '$.info.version_engine') || ' ' || (am.metadata ->> '$.assets[0].target_platform')",
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
		From: "asset_container",
		Select: []string{
			"id",
			"gsfoid",
			"name",
		},
		WhereMapping: map[string]string{
			"id":   "id",
			"oid":  "gsfoid",
			"name": "name",
		},
		OrderByMapping: map[string]string{
			"id":   "id",
			"oid":  "gsfoid",
			"name": "name",
		},
		Scan: func(rows *sql.Rows, record *Container) error {
			if err := rows.Scan(&record.ID, &record.OID, &record.Name); err != nil {
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
		From: "asset_container_asset as ca",
		Select: []string{
			"ca.id",
			"ca.position",
			"ca.win_asset_id",
			"ca.osx_asset_id",
			"concat_ws(' - ', at.name, a.gsfoid, coalesce(a.res_name, '[NULL]'), (am.metadata ->> '$.info.version_engine') || ' ' || (am.metadata ->> '$.assets[0].target_platform'))",
			"iif(ax.id is null, null, concat_ws(' - ', axt.name, ax.gsfoid, coalesce(ax.res_name, '[NULL]'), (axm.metadata ->> '$.info.version_engine') || ' ' || (axm.metadata ->> '$.assets[0].target_platform')))",
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
			"c.name",
			"p.name",
			"p.ptag",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Where(sb.EQ("cp.container_id", id))
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_package as p", "p.id = cp.package_id")
			sb.Join("asset_container as c", "c.id = cp.container_id")
			sb.OrderByAsc("cp.position").OrderByDesc("cp.id")
		},
		Scan: func(rows *sql.Rows, record *ContainerPackage) error {
			return rows.Scan(
				&record.ID,
				&record.Position,
				&record.Container,
				&record.Name,
				&record.PTag,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetPackageGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Package], error) {
	const op = "asset.Service.GetPackageGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Package]{
		From: "asset_package as p",
		Select: []string{
			"p.id",
			"p.name",
			"p.ptag",
			"datetime(p.created_date, 'unixepoch', 'localtime')",
			"p.container_id",
			"c.gsfoid || ' - ' || c.name",
		},
		WhereMapping: map[string]string{
			"id":        "p.id",
			"name":      "p.name",
			"ptag":      "p.ptag",
			"container": "c.gsfoid || ' - ' || c.name",
		},
		OrderByMapping: map[string]string{
			"id":        "p.id",
			"name":      "p.name",
			"ptag":      "p.ptag",
			"container": "c.gsfoid",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as c", "c.id = p.container_id")
		},
		Scan: func(rows *sql.Rows, record *Package) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.PTag,
				&record.CreatedDate,
				&record.Container.ID,
				&record.Container.Text,
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
				(am.metadata ->> '$.info.version_engine') || ' ' || (am.metadata ->> '$.assets[0].target_platform')
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
		TextField:    "gsfoid || ' - ' || name",
		OrderByField: "gsfoid",
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetPackagesDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "asset.Service.GetPackagesDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "asset_package",
		IDField:      "id",
		TextField:    "ptag",
		OrderByField: "ptag",
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
		Into:   "asset_container",
		Cols:   []string{"gsfoid", "name"},
		Values: []any{req.Record.OID, req.Record.Name},
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
				ub.Where(ub.EQ("id", change.ID))
				return ub
			},
		})
		return err
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrContainerExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) AddContainerAsset(ctx context.Context, req w2.SaveFormRequest[ContainerAsset], id int) (int, error) {
	const op = "asset.Service.AddContainerAsset"
	id, err := w2db.InsertFormContext(ctx, s.store.DB(), req, w2db.InsertFormOptions{
		Into:   "asset_container_asset",
		Cols:   []string{"container_id", "win_asset_id", "osx_asset_id"},
		Values: []any{id, req.Record.WINAsset.ID, req.Record.OSXAsset.ID},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrContainerAssetExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateContainerAssets(ctx context.Context, req w2.SaveGridRequest[ContainerAsset], id int) error {
	const op = "asset.Service.UpdateContainerAssets"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.SaveGridContext(ctx, tx, req, w2db.SaveGridOptions[ContainerAsset]{
			BuildUpdate: func(change ContainerAsset) *sqlbuilder.UpdateBuilder {
				ub := sqlbuilder.Update("asset_container_asset")
				w2sql.Set(ub, change.WINAsset.ID, "win_asset_id")
				w2sql.Set(ub, change.OSXAsset.ID, "osx_asset_id")
				ub.Where(ub.EQ("id", change.ID))
				return ub
			},
		})
		return err
	})
	if s.store.IsErrConstraintUnique(err) {
		return wrap.IfErr(op, ErrContainerAssetExists)
	}
	return wrap.IfErr(op, err)
}

func (s *Service) ReorderContainerAssets(ctx context.Context, req w2.ReorderGridRequest) error {
	const op = "asset.Service.ReorderContainerAssets"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.ReorderGridContext(ctx, tx, req, w2db.ReorderGridOptions{
			Update:     "asset_container_asset",
			IDField:    "id",
			SetField:   "position",
			GroupField: "container_id",
		})
		return err
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteAssets(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "asset.Service.DeleteAssets"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "asset",
		IDField: "id",
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
		From:    "asset_container_asset",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
