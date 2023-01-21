package v1

import (
	"github.com/IncubusX/go-todo-app/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary		SignUp
// @Tags			auth
// @Description	Создание аккаунта
// @ID				create-account
// @Accept			json
// @Produce		json
// @Param			input	body		entity.User	true	"account info"
// @Success		200		{object}	idResponse
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input entity.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, idResponse{
		Id: id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary		SignIn
// @Tags			auth
// @Description	Вход
// @ID				login
// @Accept			json
// @Produce		json
// @Param			input	body		signInInput	true	"credentials"
// @Success		200		{object}	signInResponse
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token: token,
	})
}
