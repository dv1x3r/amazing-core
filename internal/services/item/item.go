package item

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type Item struct {
	ID         int              `json:"id"`
	Name       w2.Field[string] `json:"name"`
	Container  w2.Dropdown      `json:"container"`
	Categories []w2.Dropdown    `json:"categories"`
	Slots      []w2.Dropdown    `json:"slots"`
}

func (s *Service) GetItemGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Item], error) {
	const op = "item.Service.GetItemGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Item]{
		From: "item as it",
		Select: []string{
			"it.id",
			"it.name",
			"it.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
			"cat.categories",
			"slt.slots",
		},
		WhereMapping: map[string]string{
			"id":        "it.id",
			"name":      "it.name",
			"container": "ac.gsfoid || ' - ' || ac.name",
		},
		OrderByMapping: map[string]string{
			"id":        "it.id",
			"name":      "it.name",
			"container": "ac.gsfoid",
		},
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
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
		},
		Scan: func(rows *sql.Rows, record *Item) error {
			var categories, slots sql.Null[string]
			if err := rows.Scan(
				&record.ID,
				&record.Name,
				&record.Container.ID,
				&record.Container.Text,
				&categories,
				&slots,
			); err != nil {
				return fmt.Errorf("scan: %w", err)
			}
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
	const op = "item.Service.CreateItem"
	var itemID int
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		var err error
		itemID, err = w2db.InsertContext(ctx, tx, w2db.InsertOptions{
			Into:   "item",
			Cols:   []string{"name", "container_id"},
			Values: []any{req.Record.Name, req.Record.Container.ID},
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
	const op = "item.Service.UpdateItem"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.UpdateContext(ctx, tx, w2db.UpdateOptions{
			Update:  "item",
			Cols:    []string{"name", "container_id"},
			Values:  []any{req.Record.Name, req.Record.Container.ID},
			IDField: "id",
			IDValue: req.Record.ID,
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
	const op = "item.Service.DeleteItems"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "item",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) GetGSFItem(ctx context.Context, platform gsf.Platform, itemID int) (types.Item, error) {
	const op = "item.Service.GetGSFItem"

	row := s.store.DB().QueryRowContext(ctx, `
			select
				name,
				container_id
			from item
			where id = ?;
		`, itemID)

	var item types.Item
	var containerID int

	if err := row.Scan(
		&item.Name,
		&containerID,
	); err != nil {
		return item, wrap.IfErr(op, err)
	}

	container, err := s.assets.GetGSFAssetContainer(ctx, platform, containerID)
	if err != nil {
		return item, wrap.IfErr(op, err)
	}
	item.AssetContainer = container

	categories, err := s.GetGSFItemCategoriesByItemID(ctx, itemID)
	if err != nil {
		return item, wrap.IfErr(op, err)
	}
	item.ItemCategories = categories

	slots, err := s.GetGSFItemAcceptableSlotsByItemID(ctx, itemID)
	if err != nil {
		return item, wrap.IfErr(op, err)
	}
	item.AcceptableSlotOIDs = slots

	return item, nil
}
