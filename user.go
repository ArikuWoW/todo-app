package todoapp

// Структура, которая представляет пользователя в приложении
// Принимает данные из http-запроса, хранит данные в базе
// Передавать данные между слоями handler -> service -> repository
type User struct {
	Id int `json:"-" db:"id"`
	// binding используется gin для валидации запроса
	// Если поля ниже не переданы пользователем, тогда gin вернет ошибку 400
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
