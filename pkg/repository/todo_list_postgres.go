// Реализует всю работу со списками задач в БД
package repository

import (
	"fmt"

	todoapp "github.com/ArikuWoW/todo-app"
	"github.com/jmoiron/sqlx"
)

// Структура содержит подключение к БД
type TodoListPostgres struct {
	db *sqlx.DB
}

// Конструктор который принимает подключение и  возвращает объект готовый к работе
func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

// Создание нового списка задач
func (r *TodoListPostgres) Create(userId int, list todoapp.TodoList) (int, error) {
	// Начинаю транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	// Вставляю title и description в таблицу todo_lists
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	// Получаем id новой записи и сохраняю в переменную id
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Привязка списка к пользователю
	// Вставляю userId and listId в таблицу users_list
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Фиксируем изменения в БД при успехе
	// Возвращаем id нового списка задач
	return id, tx.Commit()
}

// Получаем все списки пользователя
func (r *TodoListPostgres) GetAll(userId int) ([]todoapp.TodoList, error) {
	var lists []todoapp.TodoList // Сохраняем их в слайсе типа TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1", todoListTable, usersListTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

// Получаем конкретный список по id
func (r *TodoListPostgres) GetById(userId, listId int) (todoapp.TodoList, error) {
	var list todoapp.TodoList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`, todoListTable, usersListTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}
