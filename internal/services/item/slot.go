package item

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"

	"github.com/huandu/go-sqlbuilder"
)

func (s *Service) setItemSlots(ctx context.Context, tx *sql.Tx, itemID int, slots []w2.Dropdown) error {
	var slotIDs []int
	for _, slot := range slots {
		slotIDs = append(slotIDs, slot.ID.V)
	}

	// 1. delete not presented
	dlb := sqlbuilder.DeleteFrom("item_acceptable_slot")
	dlb.Where(dlb.EQ("item_id", itemID))
	if len(slots) > 0 {
		dlb.Where(dlb.NotIn("avatar_slot_id", sqlbuilder.List(slotIDs)))
	}

	query, args := dlb.BuildWithFlavor(sqlbuilder.SQLite)
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("delete slots: %w", err)
	}

	if len(slots) == 0 {
		return nil
	}

	// 2. add if not exists
	ib := sqlbuilder.NewInsertBuilder()
	ib.SetFlavor(sqlbuilder.SQLite)
	ib.InsertIgnoreInto("item_acceptable_slot")
	ib.Cols("item_id", "avatar_slot_id")
	for _, slotID := range slotIDs {
		ib.Values(itemID, slotID)
	}
	query, args = ib.BuildWithFlavor(sqlbuilder.SQLite)
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("insert slots: %w", err)
	}

	return nil
}

func (s *Service) GetGSFItemAcceptableSlotsByItemID(ctx context.Context, itemID int) ([]types.OID, error) {
	const op = "item.Service.GetGSFItemAcceptableSlotsByItemID"

	rows, err := s.store.DB().QueryContext(ctx, `
			select avs.gsfoid
			from item_acceptable_slot as ias
			join avatar_slot as avs on avs.id = ias.avatar_slot_id
			where ias.item_id = ?
			order by avs.id;
		`, itemID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	var slots []types.OID
	for rows.Next() {
		var slot types.OID
		if err := rows.Scan(&slot); err != nil {
			return slots, wrap.IfErr(op, err)
		}
		slots = append(slots, slot)
	}

	return slots, wrap.IfErr(op, rows.Err())
}
