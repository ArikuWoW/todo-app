package repository

import "github.com/jmoiron/sqlx"

// Это интерфейс, который описывает методы для работы с авторизацией в БД
type Authorization interface {
}

// Интерфейс для работы со списками задач в бд
type TodoList interface {
}

// Интерфейс для работы с элементами в списках задач
type TodoItem interface {
}

// Для простоты работы в коде, объединяем все интерфейсы в структуру
type Repository struct {
	Authorization
	TodoList
	TodoItem
}

//
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
