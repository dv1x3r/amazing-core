package services

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2db"
)

func (s *Service) GetItemAcceptableSlotsByItemID(ctx context.Context, itemID int) ([]types.OID, error) {
	const op = "game.Service.GetItemAcceptableSlotsByItemID"
	slots, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[types.OID]{
		Query: `
			select avs.gsfoid
			from item_acceptable_slot as ias
			join avatar_slot as avs on avs.id = ias.avatar_slot_id
			where ias.item_id = ?
			order by avs.id;
		`,
		Args: []any{itemID},
		Scan: func(rows *sql.Rows, slot *types.OID) error { return rows.Scan(slot) },
	})
	return slots, wrap.IfErr(op, err)
}
