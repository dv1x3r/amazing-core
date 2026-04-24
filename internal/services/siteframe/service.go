package siteframe

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/services/asset"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
	"github.com/dv1x3r/w2go/w2sql"

	"github.com/huandu/go-sqlbuilder"
)

var (
	ErrSiteFrameExists = errors.New("site frame with the same type value already exists")
)

type Service struct {
	logger *slog.Logger
	store  db.Store
	assets *asset.Service
}

func NewService(logger *slog.Logger, store db.Store, assets *asset.Service) *Service {
	return &Service{
		logger: logger,
		store:  store,
		assets: assets,
	}
}

func (s *Service) GetGSFSiteFrame(ctx context.Context, platform gsf.Platform, typeValue int32) (types.SiteFrame, error) {
	const op = "siteframe.Service.GetGSFSiteFrame"

	var sf types.SiteFrame
	sf.TypeValue = typeValue

	var containerID int
	row := s.store.DB().QueryRowContext(ctx, "select container_id from site_frame where type_value = ?;", typeValue)
	if err := row.Scan(&containerID); err != nil {
		return sf, wrap.IfErr(op, err)
	}

	container, err := s.assets.GetGSFAssetContainer(ctx, platform, containerID)
	if err != nil {
		return sf, wrap.IfErr(op, err)
	}
	sf.AssetContainer = container

	// TODO: should not be empty, but I do not know what kind of assets should be here
	if sf.AssetMap["Preload_PrefabUnity3D"] == nil {
		sf.AssetMap["Preload_PrefabUnity3D"] = []types.Asset{}
	}

	// TODO: should not be empty, but I do not know what kind of assets should be here
	if sf.AssetMap["Audio"] == nil {
		sf.AssetMap["Audio"] = []types.Asset{}
	}

	return sf, nil
}

type SiteFrame struct {
	ID        int           `json:"id"`
	TypeValue w2.Field[int] `json:"type_value"`
	Container w2.Dropdown   `json:"container"`
}

func (s *Service) GetSiteFrameGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[SiteFrame], error) {
	const op = "siteframe.Service.GetSiteFrameGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[SiteFrame]{
		From: "site_frame as sf",
		Select: []string{
			"sf.id",
			"sf.type_value",
			"sf.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
		},
		WhereMapping: map[string]string{
			"id":         "sf.id",
			"type_value": "sf.type_value",
			"container":  "ac.gsfoid || ' - ' || ac.name",
		},
		OrderByMapping: map[string]string{
			"id":         "sf.id",
			"type_value": "sf.type_value",
			"container":  "ac.gsfoid",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = sf.container_id")
		},
		Scan: func(rows *sql.Rows, record *SiteFrame) error {
			return rows.Scan(
				&record.ID,
				&record.TypeValue,
				&record.Container.ID,
				&record.Container.Text,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateSiteFrame(ctx context.Context, req w2.SaveFormRequest[SiteFrame]) (int, error) {
	const op = "siteframe.Service.CreateSiteFrame"
	id, err := w2db.InsertFormContext(ctx, s.store.DB(), req, w2db.InsertFormOptions{
		Into:   "site_frame",
		Cols:   []string{"type_value", "container_id"},
		Values: []any{req.Record.TypeValue, req.Record.Container.ID},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrSiteFrameExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateSiteFrames(ctx context.Context, req w2.SaveGridRequest[SiteFrame]) error {
	const op = "siteframe.Service.UpdateSiteFrames"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.SaveGridContext(ctx, tx, req, w2db.SaveGridOptions[SiteFrame]{
			BuildUpdate: func(change SiteFrame) *sqlbuilder.UpdateBuilder {
				ub := sqlbuilder.Update("site_frame")
				w2sql.Set(ub, change.TypeValue, "type_value")
				w2sql.Set(ub, change.Container.ID, "container_id")
				ub.Where(ub.EQ("id", change.ID))
				return ub
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrSiteFrameExists
		}
		return err
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteSiteFrames(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "siteframe.Service.DeleteSiteFrames"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "site_frame",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
