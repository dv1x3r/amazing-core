package item

import (
	"context"
	"database/sql"

	"github.com/dv1x3r/amazing-core/internal/lib/graph"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

type Category struct {
	ID         int              `json:"id"`
	OID        w2.Field[int64]  `json:"oid"`
	OIDStr     string           `json:"oid_str"`
	Name       w2.Field[string] `json:"name"`
	Parent     w2.Dropdown      `json:"parent"`
	IsPublic   w2.Field[bool]   `json:"is_public"`
	IsOutdoor  w2.Field[bool]   `json:"is_outdoor"`
	IsWalkover w2.Field[bool]   `json:"is_walkover"`
	ShowInDock w2.Field[bool]   `json:"show_in_dock"`
}

func (s *Service) GetCategoryGrid(ctx context.Context, req w2.GetGridRequest) (w2.GetGridResponse[Category], error) {
	const op = "item.Service.GetCategoryGrid"
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
		WhereMapping: map[string]string{
			"id":        "ic.id",
			"oid":       "ic.gsfoid",
			"name":      "ic.name",
			"container": "ac.gsfoid || ' - ' || ac.name",
		},
		OrderByMapping: map[string]string{
			"id":        "ic.id",
			"oid":       "ic.gsfoid",
			"name":      "ic.name",
			"container": "ac.gsfoid",
		},
		BuildBase: func(sb *sqlbuilder.SelectBuilder) {
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
			record.OIDStr = types.OIDFromInt64(record.OID.V).String()
			return nil
		},
	})
	return res, wrap.IfErr(op, err)
}

func (s *Service) CreateCategory(ctx context.Context, req w2.SaveFormRequest[Category]) (int, error) {
	const op = "item.Service.CreateCategory"
	id, err := w2db.InsertContext(ctx, s.store.DB(), w2db.InsertOptions{
		Into: "item_category",
		Cols: []string{"name", "parent_id", "is_public", "is_outdoor", "is_walkover", "show_in_dock", "gsfoid"},
		Values: []any{
			req.Record.Name,
			req.Record.Parent.ID,
			req.Record.IsPublic,
			req.Record.IsOutdoor,
			req.Record.IsWalkover,
			req.Record.ShowInDock,
			sqlbuilder.Raw("(select coalesce(max(gsfoid) + 1, 1) from item_category)"),
		},
	})
	if s.store.IsErrConstraintUnique(err) {
		return 0, wrap.IfErr(op, ErrCategoryExists)
	}
	return id, wrap.IfErr(op, err)
}

func (s *Service) UpdateCategories(ctx context.Context, req w2.SaveGridRequest[Category]) error {
	const op = "item.Service.UpdateCategory"
	err := w2db.WithinTransactionContext(ctx, s.store.DB(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := w2db.SaveGridContext(ctx, tx, req, w2db.SaveGridOptions[Category]{
			BuildOptions: func(change Category) w2db.UpdateOptions {
				return w2db.UpdateOptions{
					Update: "item_category",
					Cols:   []string{"name", "gsfoid", "parent_id", "is_public", "is_outdoor", "is_walkover", "show_in_dock"},
					Values: []any{
						change.Name,
						change.OID,
						change.Parent.ID,
						change.IsPublic,
						change.IsOutdoor,
						change.IsWalkover,
						change.ShowInDock,
					},
					IDField: "id",
					IDValue: change.ID,
				}
			},
		})
		if s.store.IsErrConstraintUnique(err) {
			return ErrCategoryExists
		} else if err != nil {
			return err
		}

		for _, category := range req.Changes {
			if cycle, err := s.hasCategoryCycle(ctx, tx, category.ID); err != nil {
				return err
			} else if cycle {
				return ErrCategoryCyclicDependency
			}
		}

		return nil
	})
	return wrap.IfErr(op, err)
}

func (s *Service) DeleteCategories(ctx context.Context, req w2.RemoveGridRequest) error {
	const op = "item.Service.DeleteCategories"
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

func (s *Service) GetGSFItemCategories(ctx context.Context, publicOnly bool) ([]types.ItemCategory, error) {
	const op = "item.Service.GetGSFItemCategories"

	rows, err := s.store.DB().QueryContext(ctx, `
			select
				ic.gsfoid,
				icp.gsfoid,
				ic.name,
				ic.is_outdoor,
				ic.is_walkover,
				ic.show_in_dock
			from item_category as ic
			left join item_category as icp on icp.id = ic.parent_id
			where (? = 0 or ic.is_public = 1);
		`, publicOnly)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	var categories []types.ItemCategory

	for rows.Next() {
		var category types.ItemCategory
		if err := rows.Scan(
			&category.OID,
			&category.ParentID,
			&category.Name,
			&category.IsOutdoor,
			&category.IsWalkover,
			&category.ShowInDock,
		); err != nil {
			return categories, wrap.IfErr(op, err)
		}
		categories = append(categories, category)
	}

	return categories, wrap.IfErr(op, rows.Err())
}
