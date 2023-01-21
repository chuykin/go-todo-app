package v1

import (
	"github.com/IncubusX/go-todo-app/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Create list
// @Security		ApiKeyAuth
// @Tags			lists
// @Description	Создание списка задач
// @ID				create-list
// @Accept			json
// @Produce		json
// @Param			input	body		entity.TodoList	true	"list info"
// @Success		200		{object}	idResponse
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/v1/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input entity.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, idResponse{
		Id: id,
	})
}

type getAllListsResponse struct {
	Data []entity.TodoList `json:"data"`
}

// @Summary		Get all lists
// @Security		ApiKeyAuth
// @Tags			lists
// @Description	Вывод всех задач
// @ID				get-all-lists
// @Accept			json
// @Produce		json
// @Success		200		{object}	getAllListsResponse
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/v1/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})
}

// @Summary		Get list by ID
// @Security		ApiKeyAuth
// @Tags			lists
// @Description	Получение списка по ИД
// @ID				get-list-by-id
// @Accept			json
// @Produce		json
// @Param			id		path		int	true	"List ID"
// @Success		200		{object}	entity.TodoList
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/v1/lists/{id} [get]
func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	list, err := h.services.TodoList.GetById(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

// @Summary		Update list
// @Security		ApiKeyAuth
// @Tags			lists
// @Description	Обновление списка задач
// @ID				update-list
// @Accept			json
// @Produce		json
// @Param			id		path		int				true	"List ID"
// @Param			input	body		entity.TodoList	true	"list info"
// @Success		200		{object}	statusResponse
// @Success		200		{object}	statusResponse
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/v1/lists/{id} [put]
func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input entity.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.TodoList.Update(userId, listId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary		Delete list
// @Security		ApiKeyAuth
// @Tags			lists
// @Description	Удаление списка задач
// @ID				delete-list
// @Accept			json
// @Produce		json
// @Param			id		path		int	true	"List ID"
// @Success		200		{object}	statusResponse
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/v1/lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.TodoList.Delete(userId, listId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
