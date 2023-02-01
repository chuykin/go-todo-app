package repository

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IncubusX/go-todo-app/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTodoItem_Create(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewTodoItem(sqlxDB)

	type args struct {
		listId int
		item   entity.TodoItem
	}
	type mockBehavior func(args args, id int)

	tt := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "Ok",
			args: args{
				listId: 1,
				item: entity.TodoItem{
					Title:       "test title",
					Description: "test desc",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO list_items").WithArgs(args.listId, id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Failed start Tx",
			args: args{
				listId: 1,
				item: entity.TodoItem{
					Title:       "test title",
					Description: "test desc",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin().WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "Empty Fields",
			args: args{
				listId: 1,
				item: entity.TodoItem{
					Title:       "",
					Description: "test desc",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd Insert Rollback",
			args: args{
				listId: 1,
				item: entity.TodoItem{
					Title:       "test title",
					Description: "test desc",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO list_items").WithArgs(args.listId, id).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.args, tc.id)

			got, err := r.Create(tc.args.listId, tc.args.item)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.id, got)
			}
		})
	}
}

func TestTodoItem_GetAll(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewTodoItem(sqlxDB)

	type mockBehavior func(userId, listId int)

	tt := []struct {
		name             string
		mockBehavior     mockBehavior
		userId           int
		listId           int
		expectedResponse []entity.TodoItem
		wantErr          bool
	}{
		{
			name: "Ok",
			mockBehavior: func(userId, listId int) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "test title 1", "test desc 1", false).
					AddRow(2, "test title 2", "test desc 2", false).
					AddRow(3, "test title 3", "test desc 3", false)
				mock.ExpectQuery(`SELECT (.+) FROM todo_items AS ti 
												INNER JOIN list_items AS li ON li.item_id = ti.id 
												INNER JOIN user_lists AS ul ON ul.list_id = li.list_id 
												WHERE (.+)`).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			userId: 1,
			listId: 1,
			expectedResponse: []entity.TodoItem{
				{Id: 1, Title: "test title 1", Description: "test desc 1", Done: false},
				{Id: 2, Title: "test title 2", Description: "test desc 2", Done: false},
				{Id: 3, Title: "test title 3", Description: "test desc 3", Done: false},
			},
		},
		{
			name: "Empty Items",
			mockBehavior: func(userId, listId int) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})
				mock.ExpectQuery(`SELECT (.+) FROM todo_items AS ti 
												INNER JOIN list_items AS li ON li.item_id = ti.id 
												INNER JOIN user_lists AS ul ON ul.list_id = li.list_id 
												WHERE (.+)`).
					WithArgs(1, -1).WillReturnRows(rows).WillReturnError(sql.ErrNoRows)
			},
			userId:           1,
			listId:           -1,
			expectedResponse: []entity.TodoItem(nil),
			wantErr:          true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.userId, tc.listId)

			got, err := r.GetAll(tc.userId, tc.listId)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedResponse, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItem_GetById(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewTodoItem(sqlxDB)

	type args struct {
		userId int
		itemId int
	}
	type mockBehavior func()

	tt := []struct {
		name             string
		mockBehavior     mockBehavior
		args             args
		expectedResponse entity.TodoItem
		wantErr          bool
	}{
		{
			name: "Ok",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "test title 1", "test desc 1", false)
				mock.ExpectQuery(`SELECT (.+) FROM todo_items AS ti 
												INNER JOIN list_items AS li ON li.item_id = ti.id 
												INNER JOIN user_lists AS ul ON ul.list_id = li.list_id 
												WHERE (.+)`).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			args: args{
				userId: 1,
				itemId: 1,
			},
			expectedResponse: entity.TodoItem{Id: 1, Title: "test title 1", Description: "test desc 1", Done: false},
		},
		{
			name: "Empty List",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery(`SELECT (.+) FROM todo_items AS ti 
												INNER JOIN list_items AS li ON li.item_id = ti.id 
												INNER JOIN user_lists AS ul ON ul.list_id = li.list_id 
												WHERE (.+)`).
					WithArgs(1, -1).WillReturnRows(rows).WillReturnError(errors.New("some error"))
			},
			args: args{
				userId: 1,
				itemId: -1,
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			got, err := r.GetById(tc.args.userId, tc.args.itemId)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItem_Delete(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewTodoItem(sqlxDB)

	type args struct {
		userId int
		itemId int
	}
	type mockBehavior func()

	tt := []struct {
		name             string
		mockBehavior     mockBehavior
		args             args
		expectedResponse entity.TodoItem
		wantErr          bool
	}{
		{
			name: "Ok",
			mockBehavior: func() {
				mock.ExpectExec(`DELETE FROM todo_items AS ti USING user_lists as ul, list_items as li 
												WHERE  ti.id = li.item_id AND 
														li.list_id = ul.list_id AND 
														(.+);`).
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				userId: 1,
				itemId: 1,
			},
		},
		{
			name: "Bad Connection",
			mockBehavior: func() {

				mock.ExpectExec(`DELETE FROM todo_items AS ti USING user_lists as ul, list_items as li 
												WHERE  ti.id = li.item_id AND 
														li.list_id = ul.list_id AND 
														(.+);`).
					WithArgs(1, -1).WillReturnError(driver.ErrBadConn)
			},
			args: args{
				userId: 1,
				itemId: -1,
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			err := r.Delete(tc.args.userId, tc.args.itemId)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItem_Update(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewTodoItem(sqlxDB)

	type args struct {
		userId int
		itemId int
		input  entity.UpdateItemInput
	}
	type mockBehavior func()
	var (
		testTitle = "Title test1"
		testDesc  = "Desc test1"
		testDone  = true
	)

	tt := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      bool
	}{
		{
			name: "Ok_All",
			mockBehavior: func() {
				mock.ExpectExec(`UPDATE todo_items AS ti SET title=(.+), description=(.+), done=(.+) 
												FROM user_lists AS ul, list_items AS li 
												WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = (.+) AND ti.id = (.+)`).
					WithArgs(testTitle, testDesc, testDone, 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				userId: 1,
				itemId: 1,
				input:  entity.UpdateItemInput{Title: &testTitle, Description: &testDesc, Done: &testDone},
			},
		},
		{
			name: "Ok_Title",
			mockBehavior: func() {
				mock.ExpectExec(`UPDATE todo_items AS ti SET title=(.+) 
												FROM user_lists AS ul, list_items AS li 
                        						WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = (.+) AND ti.id = (.+)`).
					WithArgs(testTitle, 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				userId: 1,
				itemId: 1,
				input:  entity.UpdateItemInput{Title: &testTitle},
			},
		},
		{
			name: "Ok_Description",
			mockBehavior: func() {
				mock.ExpectExec(`UPDATE todo_items AS ti SET description=(.+) 
												FROM user_lists AS ul, list_items AS li 
                        						WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = (.+) AND ti.id = (.+)`).
					WithArgs(testDesc, 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				userId: 1,
				itemId: 1,
				input:  entity.UpdateItemInput{Description: &testDesc},
			},
		},
		{
			name: "Ok_Done",
			mockBehavior: func() {
				mock.ExpectExec(`UPDATE todo_items AS ti SET done=(.+) 
												FROM user_lists AS ul, list_items AS li 
                        						WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = (.+) AND ti.id = (.+)`).
					WithArgs(testDone, 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				userId: 1,
				itemId: 1,
				input:  entity.UpdateItemInput{Done: &testDone},
			},
		},
		{
			name: "Bad Connection",
			mockBehavior: func() {
				mock.ExpectExec(`UPDATE todo_items AS ti SET title=(.+), description=(.+), done=(.+) 
												FROM user_lists AS ul, list_items AS li 
												WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = (.+) AND ti.id = (.+)`).
					WithArgs(testTitle, testDesc, testDone, 1, 1).WillReturnError(driver.ErrBadConn)
			},
			args: args{
				userId: 1,
				itemId: 1,
				input:  entity.UpdateItemInput{Title: &testTitle, Description: &testDesc, Done: &testDone},
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			err := r.Update(tc.args.userId, tc.args.itemId, tc.args.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
