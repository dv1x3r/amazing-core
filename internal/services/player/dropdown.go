package player

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

func (s *Service) GetPlayerAvatarsDropdown(ctx context.Context, req w2.GetDropdownRequest, playerID int) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "player.Service.GetPlayerAvatarsDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "player_avatar as pa",
		IDField:      "pa.id",
		TextField:    "(pa.gsfoid || ' - ' || pa.name)",
		OrderByField: "pa.gsfoid",
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Where(sb.EQ("pa.player_id", playerID))
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetPlayerItemsDropdown(ctx context.Context, req w2.GetDropdownRequest, playerID int) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "player.Service.GetPlayerItemsDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "player_item as pi",
		IDField:      "pi.id",
		TextField:    "(pi.gsfoid || ' - ' || it.name || ' x' || pi.quantity)",
		OrderByField: "pi.gsfoid",
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("item as it", "it.id = pi.item_id")
			sb.Where(sb.EQ("pi.player_id", playerID))
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetPlayerOutfitsDropdown(ctx context.Context, req w2.GetDropdownRequest, playerID int) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "player.Service.GetPlayerOutfitsDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "player_avatar_outfit as po",
		IDField:      "po.id",
		TextField:    "(pa.name || ' / Outfit #' || po.outfit_no)",
		OrderByField: "pa.name, po.outfit_no",
		BuildSelect: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("player_avatar as pa", "pa.id = po.player_avatar_id")
			sb.Where(sb.EQ("pa.player_id", playerID))
		},
	})
	return res, wrap.IfErr(op, err)
}
