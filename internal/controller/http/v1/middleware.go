package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AuthorizationHeader  = "Authorization"
	userCtx              = "userId"
	ErrEmptyAuthHeader   = "auth header is empty"
	ErrEmptyToken        = "token is empty"
	ErrInvalidAuthHeader = "invalid auth header"
	ErrFailedParseToken  = "failed to parse token"
	ErrUserNotFound      = "user id not found"
	ErrUserInvalidType   = "user id is of invalid type"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(AuthorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, ErrEmptyAuthHeader)
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader)
		return
	}

	if headerParts[1] == "" {
		newErrorResponse(c, http.StatusUnauthorized, ErrEmptyToken)
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, ErrFailedParseToken)
		return
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, ErrUserNotFound)
		return 0, errors.New(ErrUserNotFound)
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, ErrUserInvalidType)
		return 0, errors.New(ErrUserInvalidType)
	}
	return idInt, nil
}
