package gsf

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type UnixTime struct {
	time.Time
}

func (t *UnixTime) Scan(value any) error {
	var n sql.NullInt64
	if err := n.Scan(value); err != nil {
		return err
	}

	if n.Valid {
		t.Time = time.Unix(n.Int64, 0).UTC()
	} else {
		t.Time = time.Time{}
	}

	return nil
}

func (t UnixTime) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.UTC().Unix(), nil
}
