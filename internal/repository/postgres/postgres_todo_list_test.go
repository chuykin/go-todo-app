package repository

import (
	"database/sql"
	"errors"
	"github.com/IncubusX/go-todo-app/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestTodoList_Create(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewTodoList(sqlxDB)

	type args struct {
		userId int
		list   entity.TodoList
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
				userId: 1,
				list: entity.TodoList{
					Title:       "test title",
					Description: "test desc",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_lists").WithArgs(args.list.Title, args.list.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO user_lists").WithArgs(args.userId, id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{
				userId: 1,
				list: entity.TodoList{
					Title:       "",
					Description: "test desc",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO todo_lists").WithArgs(args.list.Title, args.list.Description).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd Insert Rollback",
			args: args{
				userId: 1,
				list: entity.TodoList{
					Title:       "test title",
					Description: "test desc",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_lists").WithArgs(args.list.Title, args.list.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO user_lists").WithArgs(args.userId, id).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.args, tc.id)

			got, err := r.Create(tc.args.userId, tc.args.list)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.id, got)
			}
		})
	}
}

func TestTodoList_GetAll(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewTodoList(sqlxDB)

	type mockBehavior func(userId int)

	tt := []struct {
		name             string
		mockBehavior     mockBehavior
		userId           int
		expectedResponse []entity.TodoList
		wantErr          bool
	}{
		{
			name: "Ok",
			mockBehavior: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "test title 1", "test desc 1").
					AddRow(2, "test title 2", "test desc 2").
					AddRow(3, "test title 3", "test desc 3")
				mock.ExpectQuery("SELECT (.+) FROM todo_lists AS tl INNER JOIN user_lists AS ul ON (.+) WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			userId: 1,
			expectedResponse: []entity.TodoList{
				{Id: 1, Title: "test title 1", Description: "test desc 1"},
				{Id: 2, Title: "test title 2", Description: "test desc 2"},
				{Id: 3, Title: "test title 3", Description: "test desc 3"},
			},
		},
		{
			name: "Empty List",
			mockBehavior: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "title", "description"})
				mock.ExpectQuery("SELECT (.+) FROM todo_lists AS tl INNER JOIN user_lists AS ul ON (.+) WHERE (.+)").
					WithArgs(-1).WillReturnRows(rows).WillReturnError(errors.New("sql: no rows in result set"))
			},
			userId:           -1,
			expectedResponse: []entity.TodoList(nil),
			wantErr:          true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.userId)

			got, err := r.GetAll(tc.userId)
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

func TestTodoList_GetById(t *testing.T) {
	mockDB, mock, _ := sqlmock.Newx()
	defer func(mockDB *sqlx.DB) {
		_ = mockDB.Close()
	}(mockDB)
	//sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewTodoList(mockDB)

	type args struct {
		userId int
		listId int
	}
	type mockBehavior func(userId, listId int)

	tt := []struct {
		name             string
		mockBehavior     mockBehavior
		args             args
		expectedResponse entity.TodoList
		wantErr          bool
	}{
		{
			name: "Ok",
			mockBehavior: func(userId, listId int) {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "test title 1", "test desc 1")
				mock.ExpectQuery("SELECT (.+) FROM todo_lists AS tl INNER JOIN user_lists AS ul ON tl.id = ul.list_id WHERE (.+);").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			args: args{
				userId: 1,
				listId: 1,
			},
			expectedResponse: entity.TodoList{Id: 1, Title: "test title 1", Description: "test desc 1"},
		},
		{
			name: "Empty List",
			mockBehavior: func(userId, listId int) {
				//rows := sqlmock.NewRows([]string{"id", "title", "description"})

				mock.ExpectQuery("SELECT (.+) FROM todo_lists tl INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, -1).WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.args.userId, tc.args.listId)

			got, err := r.GetById(tc.args.userId, tc.args.listId)
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
