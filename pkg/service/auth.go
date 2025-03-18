// Реализует бизнес-логику авторизации и обработку пароля, не работая на прямую с БД
package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	todoapp "github.com/ArikuWoW/todo-app"
	"github.com/ArikuWoW/todo-app/pkg/repository"
	"github.com/golang-jwt/jwt"
)

const (
	// Доп строка, добавляемая к паролю перед хешированием
	salt = "foahnfoafjpfafsdcvf"
	// Ключ для подписания jwt-токенов
	signingKey = "ffgergerger498gfg#dsf"
	// Время жизни токена
	tokenTTL = 12 * time.Hour
)

// Структура данных внутри jwt токена
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

// Сервис обрабатывающий авторизацию
type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// Метод создания пользователя
// Принимает объект user из хендлера и хеширует пароль
// Дальше передает пользователя в метод репозитория
// ВОзвращает id нового пользователя или же ошибку
func (s *AuthService) CreateUser(user todoapp.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

// Получает логин и пароль от хендлера (клиента)
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	// Хешируем пароль
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			// Когда токен истекает
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			// Когда был создан
			IssuedAt: time.Now().Unix(),
		},
		user.Id,
	})
	// Подпись токена
	return token.SignedString([]byte(signingKey))
}

// Функция парсинга и валидации jwt-токена
// Нужно чтобы понимать, кто делает запрос и разрешщат доступ защищенным маршрутам
// Получает jwt токен от клиента
// Проверяем подлинность клиента
// Извлекаем id пользователя
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	// Принимаем токен, функцией parseWithClaims
	// Парсим его используя структуру tokenClaims
	// Внутренняя функция-подпись проверяет что токен подписан
	// Если нет ошибка что токен подделан

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		// Если все хорошо возвращаю секретный ключ
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	// Получение данных из токена
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

// Функция хеширования пароля
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
