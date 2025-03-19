// хендлеры для api-эндпойнтов связанных с TodoList
package handler

import (
	"net/http"
	"strconv"

	todoapp "github.com/ArikuWoW/todo-app"
	"github.com/gin-gonic/gin"
)

// Обработчик создания списка
func (h *Handler) createList(c *gin.Context) {
	// Получаем id пользователя из jwt обращаясь к middleware
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// Читаем json запрос и парсим в структуру TodoList
	var input todoapp.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Вызов бизнес логики создания
	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// Структура для ответа на запрос получения всех списков
type getAllListResponse struct {
	Data []todoapp.TodoList `json:"data"`
}

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

	c.JSON(http.StatusOK, getAllListResponse{
		Data: lists,
	})
}
func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// Читаю id из url и преобразуем в число
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// Вызов логики для получения по id
	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

// Обработчик запроса на обновление списков
func (h *Handler) updateList(c *gin.Context) {
	// Также получаем id из middleware
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// Читаю id из url и преобразуем в число
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// читаем json запрос с новыми данными для списка
	var input todoapp.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// После всех проверок вызываем бизнес логику
	if err := h.services.Update(userId, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// При успехе выводим Ok
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// Обработчик запроса на удаление
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	// Читаю id из url и преобразуем в число
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// Вызов логики для получения по id
	err = h.services.TodoList.Delete(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
