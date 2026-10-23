package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

var ErrItemExists = errors.New("item with the same name or container already exists")

type Item struct {
	ID         int              `json:"id"`
	Container  w2.Dropdown      `json:"container"`
	Name       w2.Field[string] `json:"name"`
	Categories []w2.Dropdown    `json:"categories"`
	Slots      []w2.Dropdown    `json:"slots"`

	Icon             string `json:"icon"`
	HasWin3D         bool   `json:"has_win_3d"`
	HasOSX3D         bool   `json:"has_osx_3d"`
	HasText          bool   `json:"has_text"`
	HasInventoryIcon bool   `json:"has_inventory_icon"`
	HasShopIcon      bool   `json:"has_shop_icon"`
}

func (s *Service) GetItemGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Item], error) {
	const op = "api.Service.GetItemGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Item]{
		From: "item as it",
		Select: []string{
			"it.id",
			"it.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
			"it.name",
			"cat.categories",
			"slt.slots",
			"icon.cdnid",
			"coalesce(flags.has_win_3d, 0)",
			"coalesce(flags.has_osx_3d, 0)",
			"coalesce(flags.has_text, 0)",
			"coalesce(flags.has_inventory_icon, 0)",
			"coalesce(flags.has_shop_icon, 0)",
		},
		Where: map[string]string{
			"id":        "it.id",
			"container": "ac.gsfoid || ' - ' || ac.name",
			"name":      "it.name",
		},
		OrderBy: map[string]string{
			"id":        "it.id",
			"container": "ac.gsfoid",
			"name":      "it.name",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = it.container_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
					select
						icm.item_id,
						json_group_array(
							json_object('id', ic.id, 'text', ic.name)
							order by ic.name
						) as categories
					from item_category_map as icm
					join item_category as ic on ic.id = icm.category_id
					group by icm.item_id
				) as cat`, "cat.item_id = it.id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
					select
						ias.item_id,
						json_group_array(
							json_object('id', avs.id, 'text', avs.slot)
							order by avs.slot
						) as slots
					from item_acceptable_slot as ias
					join avatar_slot as avs on avs.id = ias.avatar_slot_id
					group by ias.item_id
				) as slt`, "slt.item_id = it.id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, containerIconJoin, "icon.container_id = it.container_id and icon.rn = 1")
			sb.JoinWithOption(sqlbuilder.LeftJoin, `(
					select
						am.container_id,
						max(iif(at.name = 'Prefab_Unity3D' and a.asset_group_id is null, 1, 0)) as has_win_3d,
						max(iif(at.name = 'Prefab_Unity3D' and a.asset_group_id is null and am.osx_asset_id is not null, 1, 0)) as has_osx_3d,
						max(iif(at.name = 'Text', 1, 0)) as has_text,
						max(iif(at.name = 'Images' and ag.name = 'Inventory Icon', 1, 0)) as has_inventory_icon,
						max(iif(at.name = 'Images' and (a.asset_group_id is null or (ag.name in ('Exterior Icon', 'Interior Icon') and a.res_name like '%tier1%')), 1, 0)) as has_shop_icon
					from asset_container_assetmap as am
					join asset as a on a.id = am.win_asset_id
					left join asset_type as at on at.id = a.asset_type_id
					left join asset_group as ag on ag.id = a.asset_group_id
					group by am.container_id
				) as flags`, "flags.container_id = it.container_id")
		},
		Scan: func(rows *sql.Rows, record *Item) error {
			var categories, slots sql.Null[string]
			var icon *string
			if err := rows.Scan(
				&record.ID,
				&record.Container.ID,
				&record.Container.Text,
				&record.Name,
				&categories,
				&slots,
				&icon,
				&record.HasWin3D,
				&record.HasOSX3D,
				&record.HasText,
				&record.HasInventoryIcon,
				&record.HasShopIcon,
			); err != nil {
				return fmt.Errorf("scan: %w", err)
			}
			record.Icon = s.iconURL(icon)
			if categories.Valid {
				if err := json.Unmarshal([]byte(categories.V), &record.Categories); err != nil {
					return fmt.Errorf("categories unmarshal: %w", err)
				}
			}
			if slots.Valid {
				if err := json.Unmarshal([]byte(slots.V), &record.Slots); err != nil {
					return fmt.Errorf("slots unmarshal: %w", err)
				}
			}
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateItem(ctx context.Context, req w2.SaveFormRequest[Item]) (int, error) {
	const op = "api.Service.CreateItem"
	var itemID int
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		var err error
		itemID, err = w2db.InsertContext(ctx, tx, w2db.InsertOptions{
			Into: "item",
			Values: map[string]any{
				"container_id": req.Record.Container.ID,
				"name":         req.Record.Name,
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrItemExists
		} else if err != nil {
			return err
		} else if err := s.setItemCategories(ctx, tx, itemID, req.Record.Categories); err != nil {
			return err
		} else if err := s.setItemSlots(ctx, tx, itemID, req.Record.Slots); err != nil {
			return err
		}
		return nil
	})
	return itemID, wrap.IfErr(op, err)
}

func (s *Service) UpdateItem(ctx context.Context, req w2.SaveFormRequest[Item]) error {
	const op = "api.Service.UpdateItem"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.UpdateContext(ctx, tx, w2db.UpdateOptions{
			Update: "item",
			Values: map[string]any{
				"container_id": req.Record.Container.ID,
				"name":         req.Record.Name,
			},
			Where: map[string]any{
				"id": req.Record.ID,
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrItemExists
		} else if err != nil {
			return err
		} else if err := s.setItemCategories(ctx, tx, req.Record.ID, req.Record.Categories); err != nil {
			return err
		} else if err := s.setItemSlots(ctx, tx, req.Record.ID, req.Record.Slots); err != nil {
			return err
		}
		return nil
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteItems(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteItems"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "item",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}
