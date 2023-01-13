package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "user_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "list_items"

	ReconnectCount    = 5
	ReconnectCooldown = 5 * time.Second
)

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	connectString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	db, err := sqlx.Open("postgres", connectString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		for i := 0; i < ReconnectCount; i++ {
			time.Sleep(ReconnectCooldown)
			logrus.Printf("Повторная попытка подключения к БД #%d", i+1)
			err = db.Ping()
			if err == nil {
				logrus.Println("Успешное подключение к БД")
				break
			}
			if i == ReconnectCount-1 {
				return nil, err
			}
		}
	}
	return db, nil
}
