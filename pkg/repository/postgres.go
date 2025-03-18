// Файл реализует функцию подключения к БД
// Хранит константы имет таблицы и структуру настроек
package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	// Драйвер для работы с sql
	_ "github.com/lib/pq"
)

// Константы таблиц
const (
	usersTable      = "users"
	todoListTable   = "todo_lists"
	usersListTable  = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

// Структура с параметрами подключения к PostgreSql
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// Функция создает подключение к БД
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	// Проверка соединения(отправляет простой запрос к БД)
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
