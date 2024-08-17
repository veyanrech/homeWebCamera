package dbs

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/veyanrech/homeWebCamera/imagecapture/config"
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

func NewPostgres(c config.Config) DBi {

	user := c.GetString("POSTGRES_USER")
	dbname := c.GetString("POSTGRES_DBNAME")
	password := c.GetString("POSTGRES_PASSWORD")
	host := c.GetString("POSTGRES_HOST")
	// port := c.GetString("POSTGRES_PORT")

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s", user, dbname, password, host))
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
