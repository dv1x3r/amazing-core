package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dv1x3r/amazing-core/internal/lib/graph"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

var (
	ErrStoreThemeExists           = errors.New("store theme with the same container already exists")
	ErrStoreThemeCyclicDependency = errors.New("circular dependency detected (A -> B -> A)")
)

type StoreTheme struct {
	ID        int           `json:"id"`
	Icon      string        `json:"icon"`
	Container w2.Dropdown   `json:"container"`
	Parent    w2.Dropdown   `json:"parent"`
	Ordinal   w2.Field[int] `json:"ordinal"`
	Items     int           `json:"items"`
}

func (s *Service) GetStoreThemeGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[StoreTheme], error) {
	const op = "api.Service.GetStoreThemeGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[StoreTheme]{
		From: "store_theme as st",
		Select: []string{
			"st.id",
			"icon.cdnid",
			"st.container_id",
			"(ac.gsfoid || ' - ' || ac.name) as container",
			"st.parent_id",
			"pac.name as parent_name",
			"st.ordinal",
			"(select count(*) from item as it where it.store_theme_id = st.id) as items",
		},
		Where: map[string]string{
			"id":        "st.id",
			"container": "ac.gsfoid || ' - ' || ac.name",
		},
		OrderBy: map[string]string{
			"id":        "st.id",
			"container": "ac.name",
			"parent":    "pac.name",
			"ordinal":   "st.ordinal",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = st.container_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "store_theme as p", "p.id = st.parent_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, "asset_container as pac", "pac.id = p.container_id")
			sb.JoinWithOption(sqlbuilder.LeftJoin, containerIconJoin, "icon.container_id = st.container_id and icon.rn = 1")
		},
		Scan: func(rows *sql.Rows, record *StoreTheme) error {
			var icon *string
			if err := rows.Scan(
				&record.ID,
				&icon,
				&record.Container.ID,
				&record.Container.Text,
				&record.Parent.ID,
				&record.Parent.Text,
				&record.Ordinal,
				&record.Items,
			); err != nil {
				return err
			}
			record.Icon = s.iconURL(icon)
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) GetStoreThemesDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetStoreThemesDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "store_theme as st",
		IDField:      "st.id",
		TextField:    "ac.name",
		OrderByField: "ac.name",
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.Join("asset_container as ac", "ac.id = st.container_id")
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateStoreTheme(ctx context.Context, req w2.SaveFormRequest[StoreTheme]) (int, error) {
	const op = "api.Service.CreateStoreTheme"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "store_theme",
		Values: map[string]any{
			"container_id": req.Record.Container.ID,
			"parent_id":    req.Record.Parent.ID,
			"ordinal":      req.Record.Ordinal,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrStoreThemeExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateStoreTheme(ctx context.Context, req w2.SaveFormRequest[StoreTheme]) error {
	const op = "api.Service.UpdateStoreTheme"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.UpdateContext(ctx, tx, w2db.UpdateOptions{
			Update: "store_theme",
			Values: map[string]any{
				"container_id": req.Record.Container.ID,
				"parent_id":    req.Record.Parent.ID,
				"ordinal":      req.Record.Ordinal,
			},
			Where: map[string]any{
				"id": req.Record.ID,
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrStoreThemeExists
		} else if err != nil {
			return err
		}

		if cycle, err := s.hasStoreThemeCycle(ctx, tx, req.Record.ID); err != nil {
			return err
		} else if cycle {
			return ErrStoreThemeCyclicDependency
		}

		return nil
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteStoreThemes(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteStoreThemes"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "store_theme",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) hasStoreThemeCycle(ctx context.Context, tx *sql.Tx, updatedStoreThemeID int) (bool, error) {
	// build complete adjacency map: { parent_store_theme_id: [child_store_theme_id, ...] }
	rows, err := tx.QueryContext(ctx, "select parent_id, id from store_theme where parent_id is not null;")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	nodes := graph.AdjacencyMap{}

	for rows.Next() {
		var nodeID, childID int
		if err := rows.Scan(&nodeID, &childID); err != nil {
			return false, err
		}
		nodes.Add(nodeID, childID)
	}

	if err := rows.Err(); err != nil {
		return false, err
	}

	return nodes.HasCycleFrom(updatedStoreThemeID), nil
}
