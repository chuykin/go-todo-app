package repository

import (
	repository "github.com/IncubusX/go-todo-app/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type (
	Repository struct {
		Authorization
		TodoList
		TodoItem
	}
)

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: repository.NewAuth(db),
		TodoList:      repository.NewTodoList(db),
		TodoItem:      repository.NewTodoItem(db),
	}
}
