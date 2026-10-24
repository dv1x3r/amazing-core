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

var ErrNPCExists = errors.New("npc with the same name or container already exists")

type NPC struct {
	ID        int              `json:"id"`
	Container w2.Dropdown      `json:"container"`
	Zone      w2.Dropdown      `json:"zone"`
	Name      w2.Field[string] `json:"name"`
	IsLocked  w2.Field[bool]   `json:"is_locked"`

	Icon            string `json:"icon"`
	HasWin3D        bool   `json:"has_win_3d"`
	HasOSX3D        bool   `json:"has_osx_3d"`
	HasLocation     bool   `json:"has_location"`
	HasText         bool   `json:"has_text"`
	HasRelationship bool   `json:"has_relationship"`
}

func (s *Service) GetNPCGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[NPC], error) {
	const op = "api.Service.GetNPCGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[NPC]{
		From: "npc as n",
		Select: []string{
			"n.id",
			"n.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
			"n.zone_id",
			"(zc.gsfoid || ' - ' || zc.name) as zone",
			"n.name",
			"n.is_locked",
			"icon.cdnid",
			"coalesce(flags.has_win_3d, 0)",
			"coalesce(flags.has_osx_3d, 0)",
			"coalesce(flags.has_location, 0)",
			"coalesce(flags.has_text, 0)",
			"coalesce(flags.has_relationship, 0)",
		},
		Where: map[string]string{
			"id":        "n.id",
			"container": "ac.gsfoid || ' - ' || ac.name",
			"zone":      "zc.gsfoid || ' - ' || zc.name",
			"name":      "n.name",
		},
		OrderBy: map[string]string{
			"id":        "n.id",
			"container": "ac.gsfoid",
			"zone":      "zc.gsfoid",
			"name":      "n.name",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = n.container_id")
			sb.Join("zone as z", "z.id = n.zone_id")
			sb.Join("asset_container as zc", "zc.id = z.container_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, containerIconJoin, "icon.container_id = n.container_id and icon.rn = 1")
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
					select
						am.container_id,
						max(iif(at.name = 'Prefab_Unity3D' and ag.name = '3D Prefab', 1, 0)) as has_win_3d,
						max(iif(at.name = 'Prefab_Unity3D' and ag.name = '3D Prefab' and am.osx_asset_id is not null, 1, 0)) as has_osx_3d,
						max(iif(at.name = 'Location_Unity3D', 1, 0)) as has_location,
						max(iif(at.name = 'Text', 1, 0)) as has_text,
						max(iif(at.name = 'Relationship_Text', 1, 0)) as has_relationship
					from asset_container_assetmap as am
					join asset as a on a.id = am.win_asset_id
					left join asset_type as at on at.id = a.asset_type_id
					left join asset_group as ag on ag.id = a.asset_group_id
					group by am.container_id
				) as flags`, "flags.container_id = n.container_id")
		},
		Scan: func(rows *sql.Rows, record *NPC) error {
			var icon *string
			if err := rows.Scan(
				&record.ID,
				&record.Container.ID,
				&record.Container.Text,
				&record.Zone.ID,
				&record.Zone.Text,
				&record.Name,
				&record.IsLocked,
				&icon,
				&record.HasWin3D,
				&record.HasOSX3D,
				&record.HasLocation,
				&record.HasText,
				&record.HasRelationship,
			); err != nil {
				return err
			}
			record.Icon = s.iconURL(icon)
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateNPC(ctx context.Context, req w2.SaveFormRequest[NPC]) (int, error) {
	const op = "api.Service.CreateNPC"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "npc",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"zone_id":      req.Record.Zone.ID,
			"name":         req.Record.Name,
			"is_locked":    req.Record.IsLocked,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrNPCExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateNPC(ctx context.Context, req w2.SaveFormRequest[NPC]) error {
	const op = "api.Service.UpdateNPC"
	_, err := w2db.UpdateContext(ctx, s.store.DB(), w2db.UpdateOptions{
		Update: "npc",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"zone_id":      req.Record.Zone.ID,
			"name":         req.Record.Name,
			"is_locked":    req.Record.IsLocked,
		},
		Where: map[string]any{
			"id": req.Record.ID,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return ErrNPCExists
	}
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteNPCs(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteNPCs"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "npc",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
