package repository

import "github.com/IncubusX/go-todo-app/internal/entity"

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type (
	Authorization interface {
		CreateUser(user entity.User) (int, error)
		GetUser(username, password string) (entity.User, error)
	}

	TodoList interface {
		Create(userId int, input entity.TodoList) (int, error)
		GetAll(userId int) ([]entity.TodoList, error)
		GetById(userId, listId int) (entity.TodoList, error)
		Update(userId, listId int, list entity.UpdateListInput) error
		Delete(userId, listId int) error
	}

	TodoItem interface {
		Create(listId int, input entity.TodoItem) (int, error)
		GetAll(userId, listId int) ([]entity.TodoItem, error)
		GetById(userId, itemId int) (entity.TodoItem, error)
		Update(userId, itemId int, input entity.UpdateItemInput) error
		Delete(userId, itemId int) error
	}
)
