package service

import (
	todoapp "github.com/ArikuWoW/todo-app"
	"github.com/ArikuWoW/todo-app/pkg/repository"
)

// Эти интерфейсы описывают методы бизнес-логики, которые будут рабоать с методами репозитория
type Authorization interface {
	CreateUser(user todoapp.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

// Конструктор, который создает и инициализирует объякты структур с нужными зависимостями
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
