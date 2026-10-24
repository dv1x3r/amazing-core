package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/graph"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/oidtype"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

var (
	ErrCategoryExists           = errors.New("category with the same name already exists")
	ErrCategoryCyclicDependency = errors.New("circular dependency detected (A -> B -> A)")
)

type Category struct {
	ID         int              `json:"id"`
	OID        w2.Field[string] `json:"oid"`
	OIDStr     string           `json:"oid_str"`
	Name       w2.Field[string] `json:"name"`
	Parent     w2.Dropdown      `json:"parent"`
	IsPublic   w2.Field[bool]   `json:"is_public"`
	IsOutdoor  w2.Field[bool]   `json:"is_outdoor"`
	IsWalkover w2.Field[bool]   `json:"is_walkover"`
	ShowInDock w2.Field[bool]   `json:"show_in_dock"`
}

func (s *Service) GetCategoryGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Category], error) {
	const op = "api.Service.GetCategoryGrid"
	res, err := w2db.GetGridContext(ctx, s.store.DB(), req, w2db.GetGridOptions[Category]{
		From: "item_category as ic",
		Select: []string{
			"ic.id",
			"ic.gsfoid",
			"ic.name",
			"ic.is_public",
			"ic.is_outdoor",
			"ic.is_walkover",
			"ic.show_in_dock",
			"ic.parent_id",
			"icp.name as parent_name",
		},
		Where: map[string]string{
			"id":        "ic.id",
			"oid":       "ic.gsfoid",
			"name":      "ic.name",
			"container": "ac.gsfoid || ' - ' || ac.name",
		},
		OrderBy: map[string]string{
			"id":        "ic.id",
			"oid":       "ic.gsfoid",
			"name":      "ic.name",
			"container": "ac.gsfoid",
		},
		Build: func(sb *sqlbuilder.SelectBuilder) {
			sb.JoinWithOption(sqlbuilder.LeftJoin, "item_category as icp", "icp.id = ic.parent_id")
		},
		Scan: func(rows *sql.Rows, record *Category) error {
			if err := rows.Scan(
				&record.ID,
				&record.OID,
				&record.Name,
				&record.IsPublic,
				&record.IsOutdoor,
				&record.IsWalkover,
				&record.ShowInDock,
				&record.Parent.ID,
				&record.Parent.Text,
			); err != nil {
				return err
			}
			record.OIDStr = types.OIDFromString(record.OID.V).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateCategory(ctx context.Context, req w2.SaveFormRequest[Category]) (int, error) {
	const op = "api.Service.CreateCategory"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "item_category",
		Values: map[string]any{
			"gsfoid":       oidOrNext(req.Record.OID, "item_category", oidtype.UNTYPED),
			"name":         req.Record.Name,
			"parent_id":    req.Record.Parent.ID,
			"is_public":    req.Record.IsPublic,
			"is_outdoor":   req.Record.IsOutdoor,
			"is_walkover":  req.Record.IsWalkover,
			"show_in_dock": req.Record.ShowInDock,
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrCategoryExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateCategory(ctx context.Context, req w2.SaveFormRequest[Category]) error {
	const op = "api.Service.UpdateCategory"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.UpdateContext(ctx, tx, w2db.UpdateOptions{
			Update: "item_category",
			Values: map[string]any{
				"gsfoid":       oidOrNext(req.Record.OID, "item_category", oidtype.UNTYPED),
				"name":         req.Record.Name,
				"parent_id":    req.Record.Parent.ID,
				"is_public":    req.Record.IsPublic,
				"is_outdoor":   req.Record.IsOutdoor,
				"is_walkover":  req.Record.IsWalkover,
				"show_in_dock": req.Record.ShowInDock,
			},
			Where: map[string]any{
				"id": req.Record.ID,
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrCategoryExists
		} else if err != nil {
			return err
		}

		if cycle, err := s.hasCategoryCycle(ctx, tx, req.Record.ID); err != nil {
			return err
		} else if cycle {
			return ErrCategoryCyclicDependency
		}

		return nil
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteCategories(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "api.Service.DeleteCategories"
	_, err := w2db.RemoveGridContext(ctx, s.store.DB(), req, w2db.RemoveGridOptions{
		From:    "item_category",
		IDField: "id",
	})
	return wrap.IfErr(op, err)
}

func (s *Service) hasCategoryCycle(ctx context.Context, tx *sql.Tx, updatedCategoryID int) (bool, error) {
	// build complete adjacency map: { parent_category_id: [child_category_id, ...] }
	rows, err := tx.QueryContext(ctx, "select parent_id, id from item_category where parent_id is not null;")
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

	return nodes.HasCycleFrom(updatedCategoryID), nil
}

func (s *Service) setItemCategories(ctx context.Context, tx *sql.Tx, itemID int, categories []w2.Dropdown) error {
	var categoryIDs []int
	for _, category := range categories {
		categoryIDs = append(categoryIDs, category.ID.V)
	}

	// 1. delete not presented
	dlb := sqlbuilder.DeleteFrom("item_category_map")
	dlb.Where(dlb.EQ("item_id", itemID))
	if len(categories) > 0 {
		dlb.Where(dlb.NotIn("category_id", sqlbuilder.List(categoryIDs)))
	}

	query, args := dlb.BuildWithFlavor(sqlbuilder.SQLite)
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("delete categories: %w", err)
	}

	if len(categories) == 0 {
		return nil
	}

	// 2. add if not exists
	ib := sqlbuilder.NewInsertBuilder()
	ib.SetFlavor(sqlbuilder.SQLite)
	ib.InsertIgnoreInto("item_category_map")
	ib.Cols("item_id", "category_id")
	for _, categoryID := range categoryIDs {
		ib.Values(itemID, categoryID)
	}
	query, args = ib.BuildWithFlavor(sqlbuilder.SQLite)
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("insert categories: %w", err)
	}

	return nil
}
