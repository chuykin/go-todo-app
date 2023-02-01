package v1

import (
	"errors"
	"fmt"
	"github.com/IncubusX/go-todo-app/internal/service"
	mock_service "github.com/IncubusX/go-todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

// TestHandler_userIdentity Unit-test middleware userIdentity
func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)
	tt := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  AuthorizationHeader,
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "Empty auth header",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"auth header is empty"}`,
		},
		{
			name:                 "Invalid Bearer",
			headerName:           AuthorizationHeader,
			headerValue:          "Bearr token",
			token:                "token",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "Invalid Token",
			headerName:           AuthorizationHeader,
			headerValue:          "Bearer ",
			token:                "",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"token is empty"}`,
		},
		{
			name:        "Service failure",
			headerName:  AuthorizationHeader,
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(0, errors.New(ErrFailedParseToken))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"failed to parse token"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.token)

			//Регистрируем только необходимый сервис, т.к. middleware зависит от сервиса авторизации
			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			//Регистрируем функцию, которая будет обернута в Middleware
			r.GET("/test", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set(tc.headerName, tc.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}

}

func TestHandler_getUserId(t *testing.T) {
	tt := []struct {
		name                 string
		setCtx               func(c *gin.Context)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Empty userCtx",
			setCtx: func(c *gin.Context) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"user id not found"}`,
		},
		{
			name: "Wrong userCtx type",
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, "wrong type")
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"user id is of invalid type"}`,
		},
		{
			name: "Ok",
			setCtx: func(c *gin.Context) {
				c.Set(userCtx, 1)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `1`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			gin.SetMode(gin.ReleaseMode)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)

			r := gin.New()

			//Регистрируем функцию, которая будет обернута в Middleware
			r.GET("/test", func(c *gin.Context) {
				tc.setCtx(c)

				id, err := getUserId(c)
				if err == nil {
					c.String(200, fmt.Sprintf("%d", id))
				}
			})

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}

}
