package randomnames

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dv1x3r/amazing-core/internal/lib/db"
	"github.com/dv1x3r/amazing-core/internal/lib/wrap"

	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2sql/w2sqlbuilder"
	"github.com/huandu/go-sqlbuilder"
)

var (
	ErrNameNotFound = errors.New("name not found")
	ErrNameExists   = errors.New("name with the same type and name already exists")
)

type RandomName struct {
	ID       int    `json:"id"`
	PartType string `json:"part_type"`
	Name     string `json:"name"`
}

func (n *RandomName) ScanRow(scan func(dest ...any) error) error {
	return scan(&n.ID, &n.PartType, &n.Name)
}

type Service struct {
	store db.Store
}

func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetNStringsByType(ctx context.Context, namePartType string, amount int) ([]string, error) {
	const op = "randomnames.Service.GetNStringsByType"

	var records []string

	const query = "select name from random_name where part_type = ? order by random() limit ?;"
	rows, err := s.store.DB().QueryContext(ctx, query, namePartType, amount)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, wrap.IfErr(op, err)
		}
		records = append(records, name)
	}

	if err := rows.Err(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	return records, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (RandomName, error) {
	const op = "randomnames.Service.GetByID"

	var name RandomName

	const query = "select id, part_type, name from random_name where id = ?;"
	row := s.store.DB().QueryRowContext(ctx, query, id)
	if err := name.ScanRow(row.Scan); err == sql.ErrNoRows {
		return name, ErrNameNotFound
	} else if err != nil {
		return name, wrap.IfErr(op, err)
	}

	return name, nil
}

func (s *Service) Insert(ctx context.Context, name RandomName) (int, error) {
	const op = "randomnames.Service.Insert"

	const query = "insert into random_name (part_type, name) values (?, ?);"
	res, err := s.store.DB().ExecContext(ctx, query, name.PartType, name.Name)
	if s.store.IsErrConstraintUnique(err) {
		return 0, fmt.Errorf("%w: type=%s name=%s", ErrNameExists, name.PartType, name.Name)
	} else if err != nil {
		return 0, wrap.IfErr(op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, wrap.IfErr(op, err)
	}

	return int(id), nil
}

func (s *Service) UpdateByID(ctx context.Context, id int, name RandomName) error {
	const op = "randomnames.Service.UpdateByID"

	const query = "update random_name set part_type = ?, name = ? where id = ?;"
	res, err := s.store.DB().ExecContext(ctx, query, name.PartType, name.Name, id)
	if s.store.IsErrConstraintUnique(err) {
		return fmt.Errorf("%w: type=%s name=%s", ErrNameExists, name.PartType, name.Name)
	} else if err != nil {
		return wrap.IfErr(op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return wrap.IfErr(op, err)
	}

	if affected == 0 {
		return fmt.Errorf("%w: id=%d", ErrNameNotFound, id)
	}

	return nil
}

func (s *Service) GetList(ctx context.Context, r w2.GridDataRequest) ([]RandomName, int, error) {
	const op = "randomnames.Service.GetList"

	var total int
	var records []RandomName

	sb := sqlbuilder.Select("count(*)")
	sb.From("random_name")

	w2sqlbuilder.Where(sb, r, map[string]string{
		"part_type": "part_type",
		"name":      "name",
	})

	query, args := sb.BuildWithFlavor(sqlbuilder.SQLite)
	row := s.store.DB().QueryRowContext(ctx, query, args...)
	if err := row.Scan(&total); err != nil && err != sql.ErrNoRows {
		return nil, 0, wrap.IfErr(op, err)
	}

	sb.Select("id", "part_type", "name")

	w2sqlbuilder.OrderBy(sb, r, map[string]string{
		"id":        "id",
		"part_type": "part_type",
		"name":      "name",
	})

	w2sqlbuilder.Limit(sb, r)
	w2sqlbuilder.Offset(sb, r)

	query, args = sb.BuildWithFlavor(sqlbuilder.SQLite)
	rows, err := s.store.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, wrap.IfErr(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var name RandomName
		if err := name.ScanRow(rows.Scan); err != nil {
			return nil, 0, wrap.IfErr(op, err)
		}
		records = append(records, name)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, wrap.IfErr(op, err)
	}

	return records, total, nil
}

func (s *Service) Delete(ctx context.Context, ids []int) error {
	const op = "randomnames.Service.Delete"
	dlb := sqlbuilder.DeleteFrom("random_name")
	dlb.Where(dlb.In("id", sqlbuilder.List(ids)))
	query, args := dlb.BuildWithFlavor(sqlbuilder.SQLite)
	_, err := s.store.DB().ExecContext(ctx, query, args...)
	return wrap.IfErr(op, err)
}
