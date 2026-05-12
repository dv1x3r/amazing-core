package avatar

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type Avatar struct {
	ID         int              `json:"id"`
	Name       w2.Field[string] `json:"name"`
	MaxOutfits w2.Field[int]    `json:"max_outfits"`
	Container  w2.Dropdown      `json:"container"`
}

func (s *Service) GetAvatarGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Avatar], error) {
	const op = "avatar.Service.GetAvatarGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Avatar]{
		From: "avatar as av",
		Select: []string{
			"av.id",
			"av.name",
			"av.max_outfits",
			"av.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
		},
		WhereMapping: map[string]string{
			"id":          "av.id",
			"name":        "av.name",
			"max_outfits": "av.max_outfits",
			"container":   "ac.gsfoid || ' - ' || ac.name",
		},
		OrderByMapping: map[string]string{
			"id":          "av.id",
			"name":        "av.name",
			"max_outfits": "av.max_outfits",
			"container":   "ac.gsfoid",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = av.container_id")
		},
		Scan: func(rows *sql.Rows, record *Avatar) error {
			return rows.Scan(
				&record.ID,
				&record.Name,
				&record.MaxOutfits,
				&record.Container.ID,
				&record.Container.Text,
			)
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateAvatar(ctx context.Context, req w2.SaveFormRequest[Avatar]) (int, error) {
	const op = "avatar.Service.CreateAvatar"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into:   "avatar",
		Cols:   []string{"name", "max_outfits", "container_id"},
		Values: []any{req.Record.Name, req.Record.MaxOutfits, req.Record.Container.ID},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrAvatarExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateAvatar(ctx context.Context, req w2.SaveFormRequest[Avatar]) error {
	const op = "avatar.Service.UpdateAvatar"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update:  "avatar",
		Cols:    []string{"name", "max_outfits", "container_id"},
		Values:  []any{req.Record.Name, req.Record.MaxOutfits, req.Record.Container.ID},
		IDField: "id",
		IDValue: req.Record.ID,
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrAvatarExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteAvatars(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "avatar.Service.DeleteAvatars"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "avatar",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) GetGSFAvatar(ctx context.Context, platform gsf.Platform, avatarID int) (types.Avatar, error) {
	const op = "avatar.Service.GetGSFAvatar"

	row := s.store.DB().QueryRowContext(ctx, `
			select
				name,
				max_outfits,
				container_id
			from avatar
			where id = ?;
		`, avatarID)

	var avatar types.Avatar
	var containerID int

	if err := row.Scan(
		&avatar.Name,
		&avatar.MaxOutfits,
		&containerID,
	); err != nil {
		return avatar, wrap.IfErr(op, err)
	}

	container, err := s.assets.GetGSFAssetContainer(ctx, platform, containerID)
	if err != nil {
		return avatar, wrap.IfErr(op, err)
	}
	avatar.AssetContainer = container

	return avatar, nil
}
