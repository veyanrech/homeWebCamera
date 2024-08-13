package dbs

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

func NewLocalFB() *sql.DB {

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		panic(err)
	}

	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm
		db.Close()
	}()

	return db
}
