package dbs

import "database/sql"

type DBi interface {
	ReturnDB() *sql.DB
	Ping() error
}

type DBInitiater interface {
	Init() error
}
