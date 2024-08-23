package bot

import (
	"database/sql"
	"fmt"

	"github.com/veyanrech/homeWebCamera/bot/dbs"
	"github.com/veyanrech/homeWebCamera/imagecapture/config"
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

func NewDB(c config.Config) *DBOps {

	//create postgres db driver

	var sqldb dbs.DBi

	if c.GetString("DB_TYPE") == "sqlite" {
		sqldb = dbs.NewLocalFB()
	} else {
		sqldb = dbs.NewPostgres(c)
	}

	err := sqldb.Ping()
	if err != nil {
		panic(err)
	}

	res := &DBOps{
		db:     sqldb.ReturnDB(),
		logger: utils.NewFileLogger("db_info.log", "db_error.log"),
	}

	if c.GetString("LOGS_DISABLED") == "true" {
		res.logger.Disable()
	}

	err = res.Init(sqldb)
	if err != nil {
		panic(err)
	}

	return res
}

func (db *DBOps) Ping() error {
	return db.db.Ping()
}

func (db *DBOps) DeactivateChatID(chatID int64) error {
	sqlq := "UPDATE registeredchats SET active = FALSE WHERE chat_id = $1"

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sqlq, chatID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *DBOps) ActivateChatID(chatID int64) error {
	sqlq := "UPDATE registeredchats SET active = TRUE WHERE chat_id = $1"

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sqlq, chatID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type chatInfo struct {
	id     int
	chatID int64
	token  string
	active bool
}

func (c *chatInfo) GetID() int {
	return c.id
}

func (c *chatInfo) GetChatID() int64 {
	return c.chatID
}

func (c *chatInfo) GetToken() string {
	return c.token
}

func (c *chatInfo) GetActive() bool {
	return c.active
}

func (db *DBOps) FindChatIDByToken(token string) (chatInfo, error) {
	sqlq := "SELECT id, chat_id, token, active FROM registeredchats WHERE token = $1 AND active = TRUE"

	row := db.db.QueryRow(sqlq, token)

	var chatID int64
	var chatToken string
	var chatActive bool
	var id int

	err := row.Scan(&id, &chatID, &chatToken, &chatActive)

	if err != nil {
		return chatInfo{}, err
	}

	return chatInfo{
		id:     id,
		chatID: chatID,
		token:  chatToken,
		active: chatActive,
	}, nil
}

func (db *DBOps) FindChatID(chatID int64) (chatInfo, error) {
	sqlq := "SELECT id, chat_id, token, active FROM registeredchats WHERE chat_id = $1 AND active = TRUE"

	row := db.db.QueryRow(sqlq, chatID)

	var reschatID int64
	var reschatToken string
	var reschatActive bool
	var resid int

	err := row.Scan(&resid, &reschatID, &reschatToken, &reschatActive)

	if err != nil {

		if err == sql.ErrNoRows {
			return chatInfo{}, nil
		}

		return chatInfo{}, err
	}

	return chatInfo{
		id:     resid,
		chatID: reschatID,
		token:  reschatToken,
		active: reschatActive,
	}, nil
}

func (db *DBOps) RegisterChatID(chatID int64, token string) error {
	sqlq := "INSERT INTO registeredchats (chat_id, token) VALUES ($1, $2)"

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sqlq, chatID, token)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *DBOps) Init(dbinst dbs.DBi) error {
	initable, ok := dbinst.(dbs.DBInitiater)

	if !ok {
		db.logger.Error("Failed to create table chat_ids")
		return fmt.Errorf("Failed to create table chat_ids")
	}

	err := initable.Init()

	if err != nil {
		db.logger.Error(err.Error())
	}

	db.test()

	return err
}

func (db *DBOps) test() {
	ee := db.RegisterChatID(123, "test")
	fmt.Println(ee)
	r, e := db.FindChatID(123)
	fmt.Println(r, e)
}
