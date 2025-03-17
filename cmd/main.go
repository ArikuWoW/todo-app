package main

import (
	"log"

	todoapp "github.com/ArikuWoW/todo-app"
	"github.com/ArikuWoW/todo-app/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)
	// Создаем сервер типа Server
	srv := new(todoapp.Server)

	// Пытаемся запустить сервер методом Server и обрабатываем возможную ошибку
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}
