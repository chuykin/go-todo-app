package handler

import (
	"github.com/IncubusX/go-todo-app"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Create todo list
// @Security		ApiKeyAuth
// @Tags			lists
// @Description	Создание списка задач
// @ID				create-list
// @Accept			json
// @Produce		json
// @Param			input	body		todo.TodoList	true	"list info"
// @Success		200		{object}	idResponse
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input todo.TodoList
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
	Data []todo.TodoList `json:"data"`
}

// @Summary		Get all todo lists
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
// @Router			/api/lists [get]
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

// @Summary		Get todo list by ID
// @Security		ApiKeyAuth
// @Tags			lists
// @Description	Получение списка по ИД
// @ID				get-list-by-id
// @Accept			json
// @Produce		json
// @Param			id		path		int	true	"List ID"
// @Success		200		{object}	todo.TodoList
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/lists/{id} [get]
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

// @Summary		Update todo list
// @Security		ApiKeyAuth
// @Tags			lists
// @Description	Обновление списка задач
// @ID				update-list
// @Accept			json
// @Produce		json
// @Param			id		path		int				true	"List ID"
// @Param			input	body		todo.TodoList	true	"list info"
// @Success		200		{object}	statusResponse
// @Success		200		{object}	statusResponse
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/lists/{id} [put]
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

	var input todo.UpdateListInput
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

// @Summary		Delete todo list
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
// @Router			/api/lists/{id} [delete]
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
