// Реализует метод сохранения пользователя в БД
package repository

import (
	"fmt"

	todoapp "github.com/ArikuWoW/todo-app"
	"github.com/jmoiron/sqlx"
)

// СТруктура которая имеет подключение к БД
type AuthPostgres struct {
	db *sqlx.DB
}

// Конструктор, который получает подключение к БД и возвращает новый объект
// Готов выполнять запросы
func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// Формирует sql запрос с представленным user
func (r *AuthPostgres) CreateUser(user todoapp.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	// Получаем id в переменную
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// Функция которая ищет пользователя по логину и паролю
func (r *AuthPostgres) GetUser(username, password string) (todoapp.User, error) {
	var user todoapp.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
