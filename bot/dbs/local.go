package dbs

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

type LocalFB struct {
	db *sql.DB
}

func (l *LocalFB) ReturnDB() *sql.DB {
	return l.db
}

func (l *LocalFB) Init() error {
	sqlq := `CREATE TABLE IF NOT EXISTS registeredchats (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		chat_id BIGINT NOT NULL,
		token TEXT NOT NULL,
		active BOOLEAN DEFAULT TRUE
	);`

	_, err := l.db.Exec(sqlq)

	return err
}

func (l *LocalFB) Ping() error {
	return l.db.Ping()
}

func NewLocalFB() DBi {

	var res *LocalFB = nil

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		panic(err)
	}

	res = &LocalFB{
		db: db,
	}

	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm
		db.Close()
	}()

	return res
}
