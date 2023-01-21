package repository

import (
	"fmt"
	"github.com/IncubusX/go-todo-app/internal/entity"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoItem struct {
	db *sqlx.DB
}

func NewTodoItem(db *sqlx.DB) *TodoItem {
	return &TodoItem{db: db}
}

func (r *TodoItem) Create(listId int, input entity.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id;", todoItemsTable)
	row := tx.QueryRow(createItemQuery, input.Title, input.Description)
	if err = row.Scan(&itemId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2);", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItem) GetAll(userId, listId int) ([]entity.TodoItem, error) {
	var items []entity.TodoItem

	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s AS ti "+
		"INNER JOIN %s AS li ON li.item_id = ti.id "+
		"INNER JOIN %s AS ul ON ul.list_id = li.list_id "+
		"WHERE ul.user_id = $1 AND ul.list_id = $2;",
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Select(&items, query, userId, listId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItem) GetById(userId, itemId int) (entity.TodoItem, error) {
	var item entity.TodoItem

	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s AS ti "+
		"INNER JOIN %s AS li ON li.item_id = ti.id "+
		"INNER JOIN %s AS ul ON ul.list_id = li.list_id "+
		"WHERE ul.user_id = $1 AND ti.id = $2;",
		todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Get(&item, query, userId, itemId)

	return item, err
}

func (r *TodoItem) Update(userId, itemId int, input entity.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s AS ti SET %s FROM %s AS ul, %s AS li "+
		"WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d",
		todoItemsTable, setQuery, usersListsTable, listsItemsTable, argId, argId+1)

	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *TodoItem) Delete(userId, itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s AS ti USING %s as ul, %s as li "+
		"WHERE  ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2;",
		todoItemsTable, usersListsTable, listsItemsTable)
	_, err := r.db.Exec(query, userId, itemId)

	return err
}
