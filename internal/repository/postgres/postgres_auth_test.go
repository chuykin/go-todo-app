package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IncubusX/go-todo-app/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuth_CreateUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewAuth(sqlxDB)

	type mockBehavior func(args entity.User)

	tt := []struct {
		name         string
		mockBehavior mockBehavior
		args         entity.User
		id           int
		wantErr      bool
	}{
		{
			name: "Ok",
			args: entity.User{
				Name:     "name",
				Username: "username",
				Password: "password",
			},
			id: 1,
			mockBehavior: func(args entity.User) {
				mock.ExpectQuery("INSERT INTO users (.+)").WithArgs(args.Name, args.Username, args.Password).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
		},
		{
			name: "Empty Fields",
			args: entity.User{
				Name:     "Name",
				Username: "",
				Password: "password",
			},
			mockBehavior: func(args entity.User) {
				mock.ExpectQuery("INSERT INTO users").WithArgs(args.Username, args.Password).
					WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.args)

			got, err := r.CreateUser(tc.args)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				fmt.Println(err, got)
				assert.NoError(t, err)
				assert.Equal(t, tc.id, got)
			}
		})
	}
}

func TestAuth_GetUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewAuth(sqlxDB)

	type mockBehavior func(args entity.User)

	tt := []struct {
		name             string
		mockBehavior     mockBehavior
		args             entity.User
		expectedResponse entity.User
		wantErr          bool
	}{
		{
			name: "Ok",
			mockBehavior: func(args entity.User) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)
				mock.ExpectQuery(`SELECT id FROM users WHERE (.+)`).WillReturnRows(rows)
			},
			expectedResponse: entity.User{
				Id: 1,
			},
		},
		{
			name: "Empty Items",
			mockBehavior: func(args entity.User) {
				mock.ExpectQuery(`SELECT id FROM users WHERE (.+)`).WillReturnError(sql.ErrNoRows)
			},
			expectedResponse: entity.User{},
			wantErr:          true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.args)

			got, err := r.GetUser(tc.args.Username, tc.args.Password)
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
