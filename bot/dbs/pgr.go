package dbs

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"
)

type PostgresDB struct {
	db *sql.DB
}

func (p *PostgresDB) ReturnDB() *sql.DB {
	return p.db
}

func (p *PostgresDB) Ping() error {
	return p.db.Ping()
}

func NewPostgres() DBi {

	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm
		db.Close()
	}()

	return &PostgresDB{
		db: db,
	}
}
