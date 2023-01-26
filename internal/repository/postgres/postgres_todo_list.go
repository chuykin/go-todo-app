package repository

import (
	"fmt"
	"github.com/IncubusX/go-todo-app/internal/entity"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoList struct {
	db *sqlx.DB
}

func NewTodoList(db *sqlx.DB) *TodoList {
	return &TodoList{db: db}
}

func (r *TodoList) Create(userId int, input entity.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id;", todoListsTable)
	row := tx.QueryRow(createListQuery, input.Title, input.Description)
	if err = row.Scan(&id); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createUserLists := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2);", usersListsTable)
	_, err = tx.Exec(createUserLists, userId, id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoList) GetAll(userId int) ([]entity.TodoList, error) {
	var lists []entity.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s AS tl INNER JOIN %s AS ul ON tl.id = ul.list_id WHERE ul.user_id = $1;",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoList) GetById(userId, listId int) (entity.TodoList, error) {
	var list entity.TodoList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s AS tl 
								   INNER JOIN %s AS ul ON tl.id = ul.list_id 
								   WHERE ul.user_id = $1 AND tl.id = $2;`, todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoList) Update(userId, listId int, input entity.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s AS tl SET %s FROM %s AS ul WHERE tl.id = ul.list_id AND ul.user_id = $%d AND ul.list_id = $%d`,
		todoListsTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, userId, listId)
	_, err := r.db.Exec(query, args...)

	return err
}

func (r *TodoList) Delete(userId, listId int) error {
	query := fmt.Sprintf(`DELETE FROM %s AS tl USING %s as ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2;`,
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}
