package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Структура, которую будем отправлять при ошибке
type errorResponse struct {
	Message string `json:"message"`
}

// Логируем ошибку
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
