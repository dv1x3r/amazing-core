package services

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type PlayerMaze struct {
	ID       int              `json:"id"`
	OID      w2.Field[string] `json:"oid"`
	OIDStr   string           `json:"oid_str"`
	House    string           `json:"house"`
	Interior string           `json:"interior"`
	Items    int              `json:"items"`
}

func (s *Service) GetPlayerMazeGrid(ctx context.Context, req w2.GetGridRequest, playerID int) (w2.GetGridResponse[PlayerMaze], error) {
	const op = "api.Service.GetPlayerMazeGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[PlayerMaze]{
		From: "player_maze as pm",
		Select: []string{
			"pm.id",
			"pm.gsfoid",
			"coalesce(hpi.gsfoid || ' - ' || hit.name, 'Yard') as house",
			"coalesce(ac.gsfoid || ' - ' || ac.name, '') as interior",
			"(select count(*) from player_item as x where x.player_maze_id = pm.id) as items",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, "player_item as hpi", "hpi.id = pm.player_item_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "item as hit", "hit.id = hpi.item_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_container as ac", "ac.id = pm.container_id")
			sb.Where(sb.EQ("pm.player_id", playerID))
			sb.OrderByAsc("pm.id")
		},
		Scan: func(rows *sql.Rows, record *PlayerMaze) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.House,
				&record.Interior,
				&record.Items,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromString(record.OID.V).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}
