package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/oidtype"
	"github.com/dv1x3r/w2go/w2db"
)

// StoreOID is the store of every NPC, it sells every item for free.
// The type is the store type of WorldShopSVCS.WORLD_SHOP_STORE_ID.
var StoreOID = types.OID{Class: 4, Type: 72, Server: 0, Number: 1}

// WorldShopOID is WorldShopSVCS.WORLD_SHOP_STORE_ID, it sells the items of the requested subcategory.
var WorldShopOID = types.OID{Class: 4, Type: 72, Server: 0, Number: 2946993}

// ErrStoreItemNotFound is a purchase of an item the store does not have.
var ErrStoreItemNotFound = errors.New("store item not found")

// GetCurrencies returns the Coinz currency, the store prices are in Coinz.
func (s *Service) GetCurrencies(ctx context.Context, platform gsf.Platform) ([]types.Currency, error) {
	const op = "game.Service.GetCurrencies"
	var containerID int
	err := s.store.DB().QueryRowContext(ctx, "select id from asset_container where name = 'Currency Coinz';").Scan(&containerID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	container, err := s.GetAssetContainer(ctx, platform, containerID)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	return []types.Currency{{AssetContainer: container, IsDefault: gsf.NewNullable(true)}}, nil
}

// GetStoreThemes returns the World Shop categories and subcategories.
func (s *Service) GetStoreThemes(ctx context.Context, platform gsf.Platform) ([]types.StoreTheme, error) {
	const op = "game.Service.GetStoreThemes"
	type row struct {
		containerID int
		parent      *types.OID
		ordinal     int32
	}
	themeRows, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[row]{
		Query: `
			select st.container_id, pc.gsfoid, st.ordinal
			from store_theme as st
			left join store_theme as p on p.id = st.parent_id
			left join asset_container as pc on pc.id = p.container_id
			order by st.id;
		`,
		Scan: func(rows *sql.Rows, r *row) error { return rows.Scan(&r.containerID, &r.parent, &r.ordinal) },
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	themes := make([]types.StoreTheme, 0, len(themeRows))
	for _, r := range themeRows {
		container, err := s.GetAssetContainer(ctx, platform, r.containerID)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		theme := types.StoreTheme{ParentStoreThemeOID: r.parent, Ordinal: r.ordinal}
		theme.AssetContainer = container
		themes = append(themes, theme)
	}
	return themes, nil
}

// GetStoreItems returns the items as free store products, of one subcategory when themeOID is set.
// The house upgrades are sold by their house only.
func (s *Service) GetStoreItems(ctx context.Context, platform gsf.Platform, themeOID *types.OID) ([]types.StoreItem, error) {
	const op = "game.Service.GetStoreItems"

	currencies, err := s.GetCurrencies(ctx, platform)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	currency := currencies[0]

	var theme any
	if themeOID != nil {
		theme = themeOID.Int64()
	}
	itemIDs, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[int]{
		Query: `
			select i.id
			from item as i
			left join store_theme as st on st.id = i.store_theme_id
			left join asset_container as ac on ac.id = st.container_id
			where i.house_item_id is null and (? is null or ac.gsfoid = ?)
			order by i.id;
		`,
		Args: []any{theme, theme},
		Scan: func(rows *sql.Rows, itemID *int) error { return rows.Scan(itemID) },
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	storeItems := make([]types.StoreItem, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		item, err := s.GetItem(ctx, platform, itemID)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		storeItem := types.StoreItem{
			ObjectOID:     item.OID,
			StoreOID:      StoreOID,
			Quantity:      -1, // GSFStoreItem.UNLIMITED
			PurchaseLimit: -1,
			Price:         0,
			CurrencyOID:   currency.OID,
			Item:          item,
			Currency:      currency,
		}
		// The client buys a product by its OID, it is the item OID as every item is sold once.
		storeItem.OID = item.OID
		storeItem.AssetMap = types.AssetMap{}
		storeItem.AssetPackages = []types.AssetPackage{}
		storeItems = append(storeItems, storeItem)
	}
	return storeItems, nil
}

// BuyStoreItems gives the player one of every requested item and returns the new player items,
// and the interiors of the bought houses, WorldShopSVCS adds them with ZoneManager.AddMazes.
func (s *Service) BuyStoreItems(ctx context.Context, platform gsf.Platform, playerOID types.OID, storeItems []types.StoreItem) ([]types.PlayerItem, []types.PlayerMaze, error) {
	const op = "game.Service.BuyStoreItems"

	tx, err := s.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, wrap.IfErr(op, err)
	}
	defer tx.Rollback()

	type house struct {
		mazeOID     types.OID
		houseOID    types.OID
		containerID int
	}
	var itemIDs []int
	var playerItemOIDs []types.OID
	var houses []house
	for _, storeItem := range storeItems {
		var playerItemID, itemID int
		var playerItemOID types.OID
		err := tx.QueryRowContext(ctx, `
				insert into player_item (gsfoid, player_id, item_id, quantity)
				select (select coalesce(max(gsfoid), 0) + 1 from player_item), pl.id, i.id, 1
				from player as pl
				join item as i on i.container_id = (select id from asset_container where gsfoid = ?)
				where pl.gsfoid = ?
				returning id, item_id, gsfoid;
			`, storeItem.ObjectOID, playerOID).Scan(&playerItemID, &itemID, &playerItemOID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, wrap.IfErr(op, fmt.Errorf("item %v: %w", storeItem.ObjectOID, ErrStoreItemNotFound))
		}
		if err != nil {
			return nil, nil, wrap.IfErr(op, err)
		}
		itemIDs = append(itemIDs, itemID)
		playerItemOIDs = append(playerItemOIDs, playerItemOID)

		// A house gets its interior maze, the door of the house leads to the maze with the house as parent.
		h := house{houseOID: playerItemOID}
		err = tx.QueryRowContext(ctx, `
				insert into player_maze (gsfoid, player_id, player_item_id, container_id)
				select (
					select coalesce(max(pm.gsfoid) + 1, t.min_oid)
					from oid_type as t
					left join player_maze as pm on pm.gsfoid between t.min_oid and t.max_oid
					where t.id = ?
				), pi.player_id, pi.id, i.interior_container_id
				from player_item as pi
				join item as i on i.id = pi.item_id
				where pi.id = ? and i.interior_container_id is not null
				returning gsfoid, container_id;
			`, oidtype.PLAYER_MAZE, playerItemID).Scan(&h.mazeOID, &h.containerID)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return nil, nil, wrap.IfErr(op, err)
		}
		houses = append(houses, h)
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, wrap.IfErr(op, err)
	}

	playerItems := make([]types.PlayerItem, 0, len(itemIDs))
	for i, itemID := range itemIDs {
		item, err := s.GetItem(ctx, platform, itemID)
		if err != nil {
			return nil, nil, wrap.IfErr(op, err)
		}
		playerItems = append(playerItems, types.PlayerItem{
			OID:       playerItemOIDs[i],
			Item:      item,
			PlayerOID: playerOID,
			Quantity:  1,
		})
	}

	mazes := make([]types.PlayerMaze, 0, len(houses))
	for _, h := range houses {
		theme, err := s.GetAssetContainer(ctx, platform, h.containerID)
		if err != nil {
			return nil, nil, wrap.IfErr(op, err)
		}
		mazes = append(mazes, types.PlayerMaze{
			OID:        h.mazeOID,
			Name:       houseMazeName,
			PlayerOID:  playerOID,
			MazePieces: []types.PlayerMazePiece{},
			HomeTheme:  theme,
			ParentOID:  h.houseOID,
		})
	}
	return playerItems, mazes, nil
}

// GetHouseUpgrades returns the upgrades of the house of the player above its stage, the next one first,
// the house is the store (House.UpdateUpgradeItems), the stage is the number of the attached upgrades.
func (s *Service) GetHouseUpgrades(ctx context.Context, platform gsf.Platform, playerOID, houseOID types.OID) ([]types.StoreItem, error) {
	const op = "game.Service.GetHouseUpgrades"

	currencies, err := s.GetCurrencies(ctx, platform)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	currency := currencies[0]

	itemIDs, err := w2db.SelectContext(ctx, s.store.DB(), w2db.SelectOptions[int]{
		Query: `
			select u.id
			from player_item as h
			join player as pl on pl.id = h.player_id
			join item as u on u.house_item_id = h.item_id
			where h.gsfoid = ? and pl.gsfoid = ?
				and u.house_stage > (
					select count(*)
					from player_item as a
					join item as ai on ai.id = a.item_id
					where a.parent_player_item_id = h.id and ai.house_item_id = h.item_id
				)
			order by u.house_stage;
		`,
		Args: []any{houseOID, playerOID},
		Scan: func(rows *sql.Rows, itemID *int) error { return rows.Scan(itemID) },
	})
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}

	storeItems := make([]types.StoreItem, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		item, err := s.GetItem(ctx, platform, itemID)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		storeItem := types.StoreItem{
			ObjectOID:     item.OID,
			StoreOID:      houseOID,
			Quantity:      -1, // GSFStoreItem.UNLIMITED
			PurchaseLimit: -1,
			Price:         0,
			CurrencyOID:   currency.OID,
			Item:          item,
			Currency:      currency,
		}
		storeItem.OID = item.OID
		storeItem.AssetMap = types.AssetMap{}
		storeItem.AssetPackages = []types.AssetPackage{}
		storeItems = append(storeItems, storeItem)
	}
	return storeItems, nil
}
