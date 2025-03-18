package service

import "github.com/ArikuWoW/todo-app/pkg/repository"

// Эти интерфейсы описывают методы бизнес-логики, которые будут рабоать с методами репозитория
type Authorization interface {
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

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
