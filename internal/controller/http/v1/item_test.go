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

func TestTodoItemHandler_createItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, userId, listId int, inputItem entity.TodoItem)

	tt := []struct {
		name                string
		setCtx              func(c *gin.Context)
		userId              int
		listId              int
		inputBody           string
		inputItem           entity.TodoItem
		url                 string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Ok",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"Item 1", "description":"Desc 1"}`,
			inputItem: entity.TodoItem{
				Title:       "Item 1",
				Description: "Desc 1",
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/1/createItem",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int, inputItem entity.TodoItem) {
				s.EXPECT().Create(userId, listId, inputItem).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:      "Bad Request",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"Item 1", "description":"Desc 1"}`,
			inputItem: entity.TodoItem{
				Title:       "Item 1",
				Description: "Desc 1",
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/WrongPath/createItem",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int, inputItem entity.TodoItem) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "BindJSON",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"Item 1", "description":"Desc 1"`,
			inputItem: entity.TodoItem{
				Title:       "Item 1",
				Description: "Desc 1",
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/1/createItem",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int, inputItem entity.TodoItem) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"Item 1", "description":"Desc 1"}`,
			inputItem: entity.TodoItem{
				Title:       "Item 1",
				Description: "Desc 1",
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/1/createItem",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int, inputItem entity.TodoItem) {
				s.EXPECT().Create(userId, listId, inputItem).Return(0, errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:      "Bad Ctx",
			userId:    1,
			listId:    1,
			inputBody: `{"title":"Item 1", "description":"Desc 1"}`,
			inputItem: entity.TodoItem{
				Title:       "Item 1",
				Description: "Desc 1",
			},
			setCtx: func(c *gin.Context) {
			},
			url: "/api/v1/lists/1/createItem",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int, inputItem entity.TodoItem) {
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mock_service.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.userId, tc.listId, tc.inputItem)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			r := gin.New()
			r.POST("/api/v1/lists/:id/createItem", tc.setCtx, handler.createItem)

			req := httptest.NewRequest("POST", tc.url, bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTodoItemHandler_getAllItems(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, userId, listId int)

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
			url: "/api/v1/lists/1/items",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int) {
				s.EXPECT().GetAll(userId, listId).Return([]entity.TodoItem{}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[]}`,
		},
		{
			name:   "Bad Request",
			userId: 1,
			listId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/lists/WrongPath/items",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int) {
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
			url: "/api/v1/lists/1/items",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int) {
				s.EXPECT().GetAll(userId, listId).Return([]entity.TodoItem{}, errors.New(ErrServiceFailure))
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
			url: "/api/v1/lists/1/items",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, listId int) {
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mock_service.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.userId, tc.listId)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			r := gin.New()
			r.GET("/api/v1/lists/:id/items", tc.setCtx, handler.getAllItems)

			req := httptest.NewRequest("GET", tc.url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTodoItemHandler_getItemById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, userId, listId int)

	tt := []struct {
		name                string
		setCtx              func(c *gin.Context)
		userId              int
		itemId              int
		url                 string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "Ok",
			userId: 1,
			itemId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/items/1",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int) {
				s.EXPECT().GetById(userId, itemId).Return(entity.TodoItem{Id: 1, Title: "Title 1", Description: "Desc 1", Done: true}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"title":"Title 1","description":"Desc 1","done":true}`,
		},
		{
			name:   "Bad Request",
			userId: 1,
			itemId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/items/WrongPath",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "Service failure",
			userId: 1,
			itemId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/items/1",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int) {
				s.EXPECT().GetById(userId, itemId).Return(entity.TodoItem{}, errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:   "Bad Ctx",
			userId: 1,
			itemId: 1,
			setCtx: func(c *gin.Context) {
			},
			url: "/api/v1/items/1",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int) {
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mock_service.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.userId, tc.itemId)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			r := gin.New()
			r.GET("/api/v1/items/:item_id", tc.setCtx, handler.getItemById)

			req := httptest.NewRequest("GET", tc.url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTodoItemHandler_updateItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, userId, itemId int, inputItem entity.UpdateItemInput)
	var testString = "test"
	var testBool = true

	tt := []struct {
		name                string
		setCtx              func(c *gin.Context)
		userId              int
		itemId              int
		inputBody           string
		inputItem           entity.UpdateItemInput
		url                 string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Ok",
			userId:    1,
			itemId:    1,
			inputBody: `{"title":"test", "description":"test","done":true}`,
			inputItem: entity.UpdateItemInput{
				Title:       &testString,
				Description: &testString,
				Done:        &testBool,
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/items/1",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int, inputItem entity.UpdateItemInput) {
				s.EXPECT().Update(userId, itemId, inputItem).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},
		{
			name:      "Bad Request",
			userId:    1,
			itemId:    1,
			inputBody: `{"title":"test", "description":"test","done":true}`,
			inputItem: entity.UpdateItemInput{
				Title:       &testString,
				Description: &testString,
				Done:        &testBool,
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/items/WrongPath",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int, inputItem entity.UpdateItemInput) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "BindJSON",
			userId:    1,
			itemId:    1,
			inputBody: `{"title":"test", "description":"test","done":true`,
			inputItem: entity.UpdateItemInput{
				Title:       &testString,
				Description: &testString,
				Done:        &testBool,
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/items/1",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int, inputItem entity.UpdateItemInput) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			userId:    1,
			itemId:    1,
			inputBody: `{"title":"test", "description":"test","done":true}`,
			inputItem: entity.UpdateItemInput{
				Title:       &testString,
				Description: &testString,
				Done:        &testBool,
			},
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/items/1",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int, inputItem entity.UpdateItemInput) {
				s.EXPECT().Update(userId, itemId, inputItem).Return(errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:      "Bad Ctx",
			userId:    1,
			itemId:    1,
			inputBody: `{"title":"Item 1", "description":"Desc 1"}`,
			inputItem: entity.UpdateItemInput{
				Title:       &testString,
				Description: &testString,
				Done:        &testBool,
			},
			setCtx: func(c *gin.Context) {
			},
			url: "/api/v1/items/1",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int, inputItem entity.UpdateItemInput) {
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mock_service.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.userId, tc.itemId, tc.inputItem)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			r := gin.New()
			r.PUT("/api/v1/items/:item_id", tc.setCtx, handler.updateItem)

			req := httptest.NewRequest("PUT", tc.url, bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}

func TestTodoItemHandler_deleteItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoItem, userId, listId int)

	tt := []struct {
		name                string
		setCtx              func(c *gin.Context)
		userId              int
		itemId              int
		url                 string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "Ok",
			userId: 1,
			itemId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/items/1",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int) {
				s.EXPECT().Delete(userId, itemId)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},
		{
			name:   "Bad Request",
			userId: 1,
			itemId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/items/WrongPath",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "Service failure",
			userId: 1,
			itemId: 1,
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			url: "/api/v1/items/1",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int) {
				s.EXPECT().Delete(userId, itemId).Return(errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:   "Bad Ctx",
			userId: 1,
			itemId: 1,
			setCtx: func(c *gin.Context) {
			},
			url: "/api/v1/items/1",
			mockBehavior: func(s *mock_service.MockTodoItem, userId, itemId int) {
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mock_service.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.userId, tc.itemId)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			w := httptest.NewRecorder()
			r := gin.New()
			r.DELETE("/api/v1/items/:item_id", tc.setCtx, handler.deleteItem)

			req := httptest.NewRequest("DELETE", tc.url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedRequestBody, w.Body.String())
		})
	}
}
