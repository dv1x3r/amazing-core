package db

import "database/sql"

type Store interface {
	DB() *sql.DB
	DriverName() string
	IsErrConstraintUnique(error) bool
	IsErrConstraintTrigger(error) bool
	IsErrConstraintForeignKey(error) bool
}
