package service

import "github.com/IncubusX/go-todo-app/internal/entity"

type (
	Authorization interface {
		CreateUser(user entity.User) (int, error)
		GenerateToken(username, password string) (string, error)
		ParseToken(token string) (int, error)
	}

	TodoList interface {
		Create(userId int, input entity.TodoList) (int, error)
		GetAll(userId int) ([]entity.TodoList, error)
		GetById(userId, listId int) (entity.TodoList, error)
		Update(userId, listId int, input entity.UpdateListInput) error
		Delete(userId, listId int) error
	}

	TodoItem interface {
		Create(userId, listId int, input entity.TodoItem) (int, error)
		GetAll(userId, listId int) ([]entity.TodoItem, error)
		GetById(userId, itemId int) (entity.TodoItem, error)
		Update(userId, itemId int, input entity.UpdateItemInput) error
		Delete(userId, itemId int) error
	}
)
