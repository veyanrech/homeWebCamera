package dbs

import "database/sql"

func NewPostgres() *sql.DB {
	return &sql.DB{}
}
