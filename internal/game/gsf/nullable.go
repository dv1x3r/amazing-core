package gsf

import (
	"database/sql"
	"encoding/json"
)

type Null[T any] struct {
	sql.Null[T]
}

func NewNullable[T any](value T) Null[T] {
	return Null[T]{
		Null: sql.Null[T]{
			V:     value,
			Valid: true,
		},
	}
}

func (n *Null[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		var zero T
		n.V = zero
		n.Valid = false
		return nil
	}

	err := json.Unmarshal(data, &n.V)
	n.Valid = err == nil
	return err
}

func (n Null[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.V)
}
