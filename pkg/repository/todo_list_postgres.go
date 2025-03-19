// Реализует всю работу со списками задач в БД
package repository

import (
	"fmt"
	"strings"

	todoapp "github.com/ArikuWoW/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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

// Функция удаления списка
func (r *TodoListPostgres) Delete(userId, listId int) error {
	// Собираем sql запрос для удаления из бд
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2", todoListTable, usersListTable)
	// Выполняем запрос подставляя параметры
	_, err := r.db.Exec(query, userId, listId)
	return err
}

// Функция обновления существующего списка
func (r *TodoListPostgres) Update(userId, listId int, input todoapp.UpdateListInput) error {
	// Хранит sql Set выражения: title=$1, description=$2
	setValues := make([]string, 0)
	// Список значений для подстановки в запрос
	args := make([]interface{}, 0)
	// Счетчик для параметров $1, $2
	argId := 1

	// Динамическая сборка запроса
	// Проверяем наличие title и description
	if input.Title != nil {
		// Если есть добавсляем в set
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	// СБорка финального запроса
	// Объединяем слайс в строку title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListTable, setQuery, usersListTable, argId, argId+1)
	args = append(args, listId, userId)

	// Логируем сформировавшийся запрос
	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
