package handler

import (
	"github.com/ArikuWoW/todo-app/pkg/service"
	"github.com/gin-gonic/gin"
)

// Может содержать зависимости
// Через нее обработчики смогут использовать сервисы для выполнения логики
type Handler struct {
	services *service.Service
}

// Сохраняю переданный сервис, что бы в методах хендлера вызывать методы сервиса
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// Метод для инициализации всех наших эндпойнтов(маршрутов)
// тоесть определяет какие url обрабатвыаются какими функциями
func (h *Handler) InitRoutes() *gin.Engine {
	// Слушает порт, принимает http запросы и передает их обработчикам
	router := gin.New()

	// ГРуппа маршрутов с префиксом "/auth"
	// Все адреса внутри начинаются с "/auth/..."
	// Маршруты авторизации
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api")
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
				items.GET("/:item_id", h.getItemById)
				items.PUT("/:item_id", h.updateItem)
				items.DELETE("/:item_id", h.deleteItem)
			}
		}
	}

	// Возвращаем настроенный роутер
	// Потом этот роутер запускается в main, что бы сервер слушал запросы
	return router
}
