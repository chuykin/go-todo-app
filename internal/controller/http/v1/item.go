package v1

import (
	"github.com/IncubusX/go-todo-app/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Create item
// @Security		ApiKeyAuth
// @Tags			items
// @Description	Создание задачи
// @ID				create-item
// @Accept			json
// @Produce		json
// @Param			id		path		int				true	"List ID"
// @Param			input	body		entity.TodoItem	true	"item info"
// @Success		200		{object}	idResponse
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/v1/lists/{id}/items [post]
func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input entity.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.TodoItem.Create(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, idResponse{
		Id: id,
	})
}

type getAllItemsResponse struct {
	Data []entity.TodoItem `json:"data"`
}

// @Summary		Get All list item
// @Security		ApiKeyAuth
// @Tags			items
// @Description	Получение списка задач
// @ID				get-all-list-items
// @Accept			json
// @Produce		json
// @Param			id		path		int	true	"List ID"
// @Success		200		{object}	getAllItemsResponse
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/v1/lists/{id}/items [get]
func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllItemsResponse{
		Data: items,
	})
}

// @Summary		Get list item By ID
// @Security		ApiKeyAuth
// @Tags			items
// @Description	Получение конкретной задачи по ИД
// @ID				get-list-item-by-id
// @Accept			json
// @Produce		json
// @Param			id		path		int	true	"Item ID"
// @Success		200		{object}	entity.TodoItem
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/v1/items/{id} [get]
func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.services.TodoItem.GetById(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary		Update list item
// @Security		ApiKeyAuth
// @Tags			items
// @Description	Обновление задачи
// @ID				update-item
// @Accept			json
// @Produce		json
// @Param			id		path		int				true	"List ID"
// @Param			input	body		entity.TodoItem	true	"item info"
// @Success		200		{object}	statusResponse
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/v1/items/{id} [put]
func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input entity.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.TodoItem.Update(userId, itemId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

}

// @Summary		Delete list item
// @Security		ApiKeyAuth
// @Tags			items
// @Description	Удаление задачи
// @ID				delete-item
// @Accept			json
// @Produce		json
// @Param			id		path		int	true	"Item ID"
// @Success		200		{object}	statusResponse
// @Failure		400,401	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/v1/items/{id} [delete]
func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.TodoItem.Delete(userId, itemId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

}
