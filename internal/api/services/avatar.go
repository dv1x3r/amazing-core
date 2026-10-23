package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

var ErrAvatarExists = errors.New("avatar with the same name or container already exists")

type Avatar struct {
	ID        int              `json:"id"`
	Container w2.Dropdown      `json:"container"`
	Name      w2.Field[string] `json:"name"`

	Icon        string `json:"icon"`
	HasWin3D    bool   `json:"has_win_3d"`
	HasOSX3D    bool   `json:"has_osx_3d"`
	HasProperty bool   `json:"has_property"`
	HasEmotes   bool   `json:"has_emotes"`
	HasIcon     bool   `json:"has_icon"`
	HasProfile  bool   `json:"has_profile"`
}

func (s *Service) GetAvatarGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Avatar], error) {
	const op = "api.Service.GetAvatarGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Avatar]{
		From: "avatar as av",
		Select: []string{
			"av.id",
			"av.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
			"av.name",
			"icon.cdnid",
			"coalesce(flags.has_win_3d, 0)",
			"coalesce(flags.has_osx_3d, 0)",
			"coalesce(flags.has_property, 0)",
			"coalesce(flags.has_emotes, 0)",
			"coalesce(flags.has_icon, 0)",
			"coalesce(flags.has_profile, 0)",
		},
		Where: map[string]string{
			"id":        "av.id",
			"container": "ac.gsfoid || ' - ' || ac.name",
			"name":      "av.name",
		},
		OrderBy: map[string]string{
			"id":        "av.id",
			"container": "ac.gsfoid",
			"name":      "av.name",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = av.container_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, containerIconJoin, "icon.container_id = av.container_id and icon.rn = 1")
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
					select
						am.container_id,
						max(iif(at.name = 'Prefab_Unity3D', 1, 0)) as has_win_3d,
						max(iif(at.name = 'Prefab_Unity3D' and am.osx_asset_id is not null, 1, 0)) as has_osx_3d,
						max(iif(at.name = 'Text' and a.res_name like '%\_Property.txt' escape '\', 1, 0)) as has_property,
						max(iif(at.name = 'Treenode', 1, 0)) as has_emotes,
						max(iif(at.name = 'Images' and a.res_name like '%icon.png', 1, 0)) as has_icon,
						max(iif(at.name = 'Images' and a.res_name like '%profile.png', 1, 0)) as has_profile
					from asset_container_assetmap as am
					join asset as a on a.id = am.win_asset_id
					left join asset_type as at on at.id = a.asset_type_id
					group by am.container_id
				) as flags`, "flags.container_id = av.container_id")
		},
		Scan: func(rows *sql.Rows, record *Avatar) error {
			var icon *string
			if err := rows.Scan(
				&record.ID,
				&record.Container.ID,
				&record.Container.Text,
				&record.Name,
				&icon,
				&record.HasWin3D,
				&record.HasOSX3D,
				&record.HasProperty,
				&record.HasEmotes,
				&record.HasIcon,
				&record.HasProfile,
			); err != nil {
				return err
			}
			record.Icon = s.iconURL(icon)
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateAvatar(ctx context.Context, req w2.SaveFormRequest[Avatar]) (int, error) {
	const op = "api.Service.CreateAvatar"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "avatar",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"name":         req.Record.Name,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrAvatarExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateAvatar(ctx context.Context, req w2.SaveFormRequest[Avatar]) error {
	const op = "api.Service.UpdateAvatar"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "avatar",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"name":         req.Record.Name,
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrAvatarExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteAvatars(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteAvatars"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "avatar",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
