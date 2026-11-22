package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/oidtype"
	"github.com/dv1x3r/w2go/w2db"
)

var ErrMazeItemNotFound = errors.New("maze item not found")

// houseMazeName is the name of a house interior, the client does not read it:
// LoadMazeScene renames the entered interior to maze.interior.home, ZoneManager.AddMazes appends the OID.
const houseMazeName = "maze.house"

// The yard scenes, the client loads the HomeLot scenes from its build instead of downloading them.
const (
	sceneHomeLotSmall  = "OTYxMTQ4NDU5NDE5MA"
	sceneHomeLotWinter = "OTQ1MDc3NTY0MjEyNg"
	sceneEmpty         = "OTYwNzA1NTQwOTE2Ng"
)

func (s *Service) GetPlayerHome(ctx context.Context, platform gsf.Platform, ownerOID types.OID, winter bool) (*types.PlayerHome, error) {
	const op = "game.Service.GetPlayerHome"

	var yardOID types.OID
	err := s.store.DB().QueryRowContext(ctx, `
			select pm.gsfoid
			from player_maze as pm
			join player as pl on pl.id = pm.player_id
			where pl.gsfoid = ? and pm.player_item_id is null;
		`, ownerOID).Scan(&yardOID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	// The Home.PlayerMaze.HomeTheme.AssetMap["Scene_Unity3D"] asset drives scene loading via
	// LoadMazeCommand.cs -> LoadMainScene() -> AssetDownloadManager.cs -> LoadMainScene()
	cdnid := sceneHomeLotSmall
	if winter {
		cdnid = sceneHomeLotWinter
	}
	scene, err := s.GetAssetByCDNID(ctx, cdnid)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	// EmptyScene has the shared scene objects (AngleViewCamera, RadialMenuManager, MiniMap),
	// LoadMazeCommand loads it additively after the Main_Scene.
	emptyScene, err := s.GetAssetByCDNID(ctx, sceneEmpty)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	yardTheme := types.AssetContainer{AssetMap: types.AssetMap{"Scene_Unity3D": {scene, emptyScene}}}

	// ZoneManager takes the yard owner from PlayerHome.PlayerOID and LoadMazeCommand from PlayerMaze.PlayerOID,
	// ENTER_BUILDING then sends it as the location.
	// The yard maze OID is MazeManager.mazeId of the yard, the items are placed in it.
	yard := types.PlayerMaze{
		OID:        yardOID,
		Name:       "coremaze",
		PlayerOID:  ownerOID,
		MazePieces: []types.PlayerMazePiece{},
		HomeTheme:  yardTheme,
	}

	type row struct {
		mazeOID     types.OID
		houseOID    types.OID
		containerID sql.Null[int]
	}
	houseRows, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[row]{
		Query: `
			select pm.gsfoid, pi.gsfoid, pm.container_id
			from player_maze as pm
			join player as pl on pl.id = pm.player_id
			join player_item as pi on pi.id = pm.player_item_id
			where pl.gsfoid = ?
			order by pm.id;
		`,
		Args: []any{ownerOID},
		Scan: func(rows *sql.Rows, r *row) error { return rows.Scan(&r.mazeOID, &r.houseOID, &r.containerID) },
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	// The door of a house leads to the maze with the house as parent (ZoneManager.GetMazeOidByParentOid).
	mazes := []types.PlayerMaze{yard}
	for _, r := range houseRows {
		theme := types.AssetContainer{AssetMap: types.AssetMap{}, AssetPackages: []types.AssetPackage{}}
		if r.containerID.Valid {
			theme, err = s.GetAssetContainer(ctx, platform, r.containerID.V)
			if err != nil {
				return nil, wrap.IfErr(op, err)
			}
		}
		mazes = append(mazes, types.PlayerMaze{
			OID:        r.mazeOID,
			Name:       houseMazeName,
			PlayerOID:  ownerOID,
			MazePieces: []types.PlayerMazePiece{},
			HomeTheme:  theme,
			ParentOID:  r.houseOID,
		})
	}

	return &types.PlayerHome{
		PlayerMaze:  yard,
		HomeTheme:   yardTheme,
		PlayerOID:   ownerOID,
		PlayerMazes: mazes,
	}, nil
}

func (s *Service) GetMazeItems(ctx context.Context, platform gsf.Platform, mazeOID, ownerOID types.OID) ([]types.PlayerItem, error) {
	const op = "game.Service.GetMazeItems"

	type row struct {
		id     int
		itemID int
		item   types.PlayerItem
	}
	itemRows, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[row]{
		Query: `
			select pi.id, pi.item_id, pi.gsfoid, pi.quantity,
				coalesce(pi.position_x, 0), coalesce(pi.position_y, 0), coalesce(pi.position_z, 0), coalesce(pi.rotation, '')
			from player_item as pi
			join player_maze as pm on pm.id = pi.player_maze_id
			join player as pl on pl.id = pm.player_id
			where pm.gsfoid = ? and pl.gsfoid = ?
			order by pi.id;
		`,
		Args: []any{mazeOID, ownerOID},
		Scan: func(rows *sql.Rows, r *row) error {
			pos := &r.item.InventoryPosition
			return rows.Scan(&r.id, &r.itemID, &r.item.OID, &r.item.Quantity,
				&pos.ObjectPosition.X, &pos.ObjectPosition.Y, &pos.ObjectPosition.Z, &pos.Rotation)
		},
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	type childRow struct {
		parentID int
		itemID   int
		item     types.PlayerItem
	}
	children, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[childRow]{
		Query: `
			select c.parent_player_item_id, c.item_id, c.gsfoid, c.quantity, coalesce(c.ordinal, 0)
			from player_item as c
			join player_item as p on p.id = c.parent_player_item_id
			join player_maze as pm on pm.id = p.player_maze_id
			where pm.gsfoid = ?
			order by c.id;
		`,
		Args: []any{mazeOID},
		Scan: func(rows *sql.Rows, r *childRow) error {
			return rows.Scan(&r.parentID, &r.itemID, &r.item.OID, &r.item.Quantity, &r.item.Ordinal)
		},
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	items := make([]types.PlayerItem, 0, len(itemRows))
	for _, r := range itemRows {
		item, err := s.GetItem(ctx, platform, r.itemID)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		playerItem := r.item
		playerItem.Item = item
		playerItem.PlayerOID = ownerOID
		playerItem.PlayerMazeOID = &mazeOID
		for _, c := range children {
			if c.parentID != r.id {
				continue
			}
			childItem, err := s.GetItem(ctx, platform, c.itemID)
			if err != nil {
				return nil, wrap.IfErr(op, err)
			}
			child := c.item
			child.Item = childItem
			child.PlayerOID = ownerOID
			child.ParentPIOOID = &playerItem.OID
			playerItem.AttachedItems = append(playerItem.AttachedItems, child)
		}
		items = append(items, playerItem)
	}
	return items, nil
}

func (s *Service) PlaceMazeItem(ctx context.Context, playerOID, itemOID, mazeOID types.OID, position types.ObjectPosition, rotation string) error {
	const op = "game.Service.PlaceMazeItem"
	result, err := s.store.DB().ExecContext(ctx, `
			update player_item
			set player_maze_id = pm.id,
				position_x = ?, position_y = ?, position_z = ?, rotation = ?,
				parent_player_item_id = null, ordinal = null
			from player_maze as pm
			join player as pl on pl.id = pm.player_id
			where pm.gsfoid = ? and pl.gsfoid = ?
				and player_item.gsfoid = ? and player_item.player_id = pl.id
				and player_item.player_avatar_outfit_id is null;
		`, position.X, position.Y, position.Z, rotation, mazeOID, playerOID, itemOID)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return wrap.IfErr(op, err)
	}
	if rowsAffected == 0 {
		return wrap.IfErr(op, ErrMazeItemNotFound)
	}
	return nil
}

func (s *Service) AttachMazeItem(ctx context.Context, playerOID, parentOID, itemOID types.OID, ordinal int32) error {
	const op = "game.Service.AttachMazeItem"
	result, err := s.store.DB().ExecContext(ctx, `
			update player_item
			set parent_player_item_id = p.id, ordinal = ?,
				player_maze_id = null, position_x = null, position_y = null, position_z = null, rotation = null
			from player_item as p
			join player as pl on pl.id = p.player_id
			where p.gsfoid = ? and pl.gsfoid = ? and p.parent_player_item_id is null
				and player_item.gsfoid = ? and player_item.player_id = pl.id and player_item.id <> p.id
				and player_item.player_avatar_outfit_id is null;
		`, ordinal, parentOID, playerOID, itemOID)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return wrap.IfErr(op, err)
	}
	if rowsAffected == 0 {
		return wrap.IfErr(op, ErrMazeItemNotFound)
	}
	return nil
}

func (s *Service) RemoveMazeItem(ctx context.Context, playerOID, itemOID types.OID) error {
	const op = "game.Service.RemoveMazeItem"
	result, err := s.store.DB().ExecContext(ctx, `
			update player_item
			set player_maze_id = null, position_x = null, position_y = null, position_z = null, rotation = null,
				parent_player_item_id = null, ordinal = null
			where gsfoid = ? and player_id = (select id from player where gsfoid = ?);
		`, itemOID, playerOID)
	if err != nil {
		return wrap.IfErr(op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return wrap.IfErr(op, err)
	}
	if rowsAffected == 0 {
		return wrap.IfErr(op, ErrMazeItemNotFound)
	}
	return nil
}

func (s *Service) DetachMazeItems(ctx context.Context, playerOID, parentOID types.OID, itemOIDs []types.OID) error {
	const op = "game.Service.DetachMazeItems"
	for _, itemOID := range itemOIDs {
		_, err := s.store.DB().ExecContext(ctx, `
				update player_item
				set parent_player_item_id = null, ordinal = null
				where gsfoid = ? and player_id = (select id from player where gsfoid = ?)
					and parent_player_item_id = (select id from player_item where gsfoid = ?);
			`, itemOID, playerOID, parentOID)
		if err != nil {
			return wrap.IfErr(op, err)
		}
	}
	return nil
}

func (s *Service) insertPlayerYard(ctx context.Context, tx *sql.Tx, playerID int) error {
	const op = "game.Service.insertPlayerYard"
	_, err := tx.ExecContext(ctx, `
			insert into player_maze (gsfoid, player_id)
			values ((
				select coalesce(max(pm.gsfoid) + 1, t.min_oid)
				from oid_type as t
				left join player_maze as pm on pm.gsfoid between t.min_oid and t.max_oid
				where t.id = ?
			), ?);
		`, oidtype.PLAYER_MAZE, playerID)
	return wrap.IfErr(op, err)
}
