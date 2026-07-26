package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
)

func (s *Service) GetItemAcceptableSlotsByItemID(ctx context.Context, itemID int) ([]types.OID, error) {
	const op = "game.Service.GetItemAcceptableSlotsByItemID"

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
