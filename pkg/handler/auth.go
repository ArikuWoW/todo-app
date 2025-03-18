package handler

import (
	"net/http"

	todoapp "github.com/ArikuWoW/todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	// Пытаюсь считать json из тела запроса и записать в input
	var input todoapp.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Вызов бизнес-логики
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// При успешном выполнении возвращаем id
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

// Отдельная структура для входа пользователя
type singInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// функция обработчика входа пользователя
// Получает от клиента логин и пароль, проверяет существование пользователя
// Создает jwt токен и возвращает токен клиенту
func (h *Handler) signIn(c *gin.Context) {
	var input singInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Вызов бизнес-логики(сервиса)
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// При успешном выполнении возвращаем токен клиенту
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
