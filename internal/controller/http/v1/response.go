package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	ErrInvalidInputBody = "invalid input body"
	ErrServiceFailure   = "service failure"
)

type signInResponse struct {
	Token string `json:"token"`
}

type idResponse struct {
	Id int `json:"id"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
