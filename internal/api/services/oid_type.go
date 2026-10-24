package services

import (
	"context"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"
	"github.com/dv1x3r/amazing-core/internal/network/gsf/types/oidtype"
	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"

	"github.com/huandu/go-sqlbuilder"
)

func (s *Service) GetOIDTypesDropdown(ctx context.Context, req w2.GetDropdownRequest) (w2.GetDropdownResponse[w2.Dropdown], error) {
	const op = "api.Service.GetOIDTypesDropdown"
	res, err := w2db.GetDropdownContext(ctx, s.store.DB(), req, w2db.GetDropdownOptions{
		From:         "oid_type",
		IDField:      "id",
		TextField:    "name",
		OrderByField: "id",
	})
	return res, wrap.IfErr(op, err)
}

// oidOrNext returns the entered OID, or the next free OID of the OIDTypeID in the table.
func oidOrNext[T comparable](oid w2.Field[T], table string, oidTypeID any) any {
	var zero T
	if !oid.Provided || (oid.Valid && oid.V != zero) {
		return oid
	}
	return sqlbuilder.Buildf(`(
		select coalesce(max(x.gsfoid) + 1, t.min_oid)
		from oid_type as t
		left join %v as x on x.gsfoid between t.min_oid and t.max_oid
		where t.id = coalesce(%v, %v)
	)`, sqlbuilder.Raw(table), oidTypeID, oidtype.UNTYPED)
}
