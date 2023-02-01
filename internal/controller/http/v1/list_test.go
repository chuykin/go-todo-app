package v1

import (
	"bytes"
	"errors"
	"github.com/IncubusX/go-todo-app/internal/entity"
	"github.com/IncubusX/go-todo-app/internal/service"
	mock_service "github.com/IncubusX/go-todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestTodoListHandler_createList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId int, inputList entity.TodoList)

	tt := []struct {
		name                string
		setCtx              func(c *gin.Context)
		userId              int
		listId              int
		inputBody           string
		inputList           entity.TodoList
		url                 string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Ok",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"List 1", "description":"Desc 1"}`,
			inputList: entity.TodoList{
				Title:       "List 1",
				Description: "Desc 1",
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists",
			mockBehavior: func(s *mock_service.MockTodoList, userId int, inputList entity.TodoList) {
				s.EXPECT().Create(userId, inputList).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:      "BindJSON",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"List 1", "description":"Desc 1"`,
			inputList: entity.TodoList{
				Title:       "List 1",
				Description: "Desc 1",
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists",
			mockBehavior: func(s *mock_service.MockTodoList, userId int, inputList entity.TodoList) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"List 1", "description":"Desc 1"}`,
			inputList: entity.TodoList{
				Title:       "List 1",
				Description: "Desc 1",
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists",
			mockBehavior: func(s *mock_service.MockTodoList, userId int, inputList entity.TodoList) {
				s.EXPECT().Create(userId, inputList).Return(0, errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:      "Bad Ctx",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"List 1", "description":"Desc 1"}`,
			inputList: entity.TodoList{
				Title:       "List 1",
				Description: "Desc 1",
			},
			setCtx: func(c *gin.Context) {
			},
			url: "/api/v1/lists",
			mockBehavior: func(s *mock_service.MockTodoList, userId int, inputList entity.TodoList) {
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.userId, tc.inputList)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			r := gin.New()
			r.POST("/api/v1/lists", tc.setCtx, handler.createList)

			req := httptest.NewRequest("POST", tc.url, bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTodoListHandler_getAllLists(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId, listId int)

	tt := []struct {
		name                string
		setCtx              func(c *gin.Context)
		userId              int
		listId              int
		url                 string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "Ok",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
				s.EXPECT().GetAll(userId).Return([]entity.TodoList{}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[]}`,
		},
		{
			name:   "Service failure",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
				s.EXPECT().GetAll(userId).Return([]entity.TodoList{}, errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:   "Bad Ctx",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
			},
			url: "/api/v1/lists",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.userId, tc.listId)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			r := gin.New()
			r.GET("/api/v1/lists", tc.setCtx, handler.getAllLists)

			req := httptest.NewRequest("GET", tc.url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTodoListHandler_getListById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId, listId int)

	tt := []struct {
		name                string
		setCtx              func(c *gin.Context)
		userId              int
		listId              int
		url                 string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "Ok",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/1",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
				s.EXPECT().GetById(userId, listId).Return(entity.TodoList{Id: 1, Title: "Title 1", Description: "Desc 1"}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"title":"Title 1","description":"Desc 1"}`,
		},
		{
			name:   "Bad Request",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/WrongPath",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "Service failure",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/1",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
				s.EXPECT().GetById(userId, listId).Return(entity.TodoList{}, errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:   "Bad Ctx",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
			},
			url: "/api/v1/lists/1",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.userId, tc.listId)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			r := gin.New()
			r.GET("/api/v1/lists/:id", tc.setCtx, handler.getListById)

			req := httptest.NewRequest("GET", tc.url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTodoListHandler_updateList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId, listId int, inputList entity.UpdateListInput)
	var testString = "test"

	tt := []struct {
		name                string
		setCtx              func(c *gin.Context)
		userId              int
		listId              int
		inputBody           string
		inputList           entity.UpdateListInput
		url                 string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Ok",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"test", "description":"test"}`,
			inputList: entity.UpdateListInput{
				Title:       &testString,
				Description: &testString,
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/1",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int, inputList entity.UpdateListInput) {
				s.EXPECT().Update(userId, listId, inputList).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},
		{
			name:      "Bad Request",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"test", "description":"test"}`,
			inputList: entity.UpdateListInput{
				Title:       &testString,
				Description: &testString,
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/WrongPath",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int, inputList entity.UpdateListInput) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "BindJSON",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"test", "description":"test"`,
			inputList: entity.UpdateListInput{
				Title:       &testString,
				Description: &testString,
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/1",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int, inputList entity.UpdateListInput) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"test", "description":"test"}`,
			inputList: entity.UpdateListInput{
				Title:       &testString,
				Description: &testString,
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/1",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int, inputList entity.UpdateListInput) {
				s.EXPECT().Update(userId, listId, inputList).Return(errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:      "Bad Ctx",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"List 1", "description":"Desc 1"}`,
			inputList: entity.UpdateListInput{
				Title:       &testString,
				Description: &testString,
			},
			setCtx: func(c *gin.Context) {
			},
			url: "/api/v1/lists/1",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int, inputList entity.UpdateListInput) {
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.userId, tc.listId, tc.inputList)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			r := gin.New()
			r.PUT("/api/v1/lists/:id", tc.setCtx, handler.updateList)

			req := httptest.NewRequest("PUT", tc.url, bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTodoListHandler_deleteList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId, listId int)

	tt := []struct {
		name                string
		setCtx              func(c *gin.Context)
		userId              int
		listId              int
		url                 string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "Ok",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/1",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
				s.EXPECT().Delete(userId, listId)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},
		{
			name:   "Bad Request",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/WrongPath",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "Service failure",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/1",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
				s.EXPECT().Delete(userId, listId).Return(errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:   "Bad Ctx",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
			},
			url: "/api/v1/lists/1",
			mockBehavior: func(s *mock_service.MockTodoList, userId, listId int) {
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mock_service.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.userId, tc.listId)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			r := gin.New()
			r.DELETE("/api/v1/lists/:id", tc.setCtx, handler.deleteList)

			req := httptest.NewRequest("DELETE", tc.url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}
