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

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user entity.User)

	tt := []struct {
		name                string
		inputBody           string
		inputUser           entity.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name":"Test", "username":"test", "password":"qwerty"}`,
			inputUser: entity.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"username":"test", "password":"qwerty"}`,
			inputUser:           entity.User{},
			mockBehavior:        func(s *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"name":"Test", "username":"test", "password":"qwerty"}`,
			inputUser: entity.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedRequestBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user entity.User)

	tt := []struct {
		name                string
		inputBody           string
		inputUser           entity.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username":"test", "password":"qwerty"}`,
			inputUser: entity.User{
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("token", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"token":"token"}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"username":"", "password":""}`,
			inputUser:           entity.User{},
			mockBehavior:        func(s *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"name":"Test", "username":"test", "password":"qwerty"}`,
			inputUser: entity.User{
				Name:     "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("", errors.New(ErrServiceFailure))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedRequestBody)
		})
	}
}
