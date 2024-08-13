package bot

import (
	"database/sql"
	"os"

	"github.com/veyanrech/homeWebCamera/bot/dbs"
	"github.com/veyanrech/homeWebCamera/imagecapture/utils"
)

type DBActionsI interface {
	RegisterChatID(chatID int64, token string) error
	FindChatIDByToken(token string) (int64, error)
}

type DBOps struct {
	db     *sql.DB
	logger utils.Logger
}

func NewDB() *DBOps {

	//create postgres db driver

	var sqldb *sql.DB

	environment := "prod"

	if os.Getenv("webcamerabot_env") == "dev" {
		environment = "dev"
	}

	if environment == "dev" {
		sqldb = dbs.NewLocalFB()
	} else {
		sqldb = dbs.NewPostgres()
	}

	return &DBOps{
		db:     sqldb,
		logger: utils.NewFileLogger("db_info.log", "db_error.log"),
	}

}

func connect() (*sql.DB, error) {
	return sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
}

func (db *DBOps) RegisterChatID(chatID int64, token string) {
	sqlq := "INSERT INTO  (chat_id, token) VALUES (?, ?)"
	_, err := db.db.Exec(sqlq, chatID, token)
	if err != nil {

	}
}
