package main

import (
	"log"
	"os"

	todoapp "github.com/ArikuWoW/todo-app"
	"github.com/ArikuWoW/todo-app/pkg/handler"
	"github.com/ArikuWoW/todo-app/pkg/repository"
	"github.com/ArikuWoW/todo-app/pkg/service"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	// Инициализируем базу передав нужные значения
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	//Создаю инстанты всех слоев

	// Работа с БД
	repos := repository.NewRepository(db)

	// Логика использует БД
	services := service.NewService(repos)

	handlers := handler.NewHandler(services)
	// Создаем сервер типа Server
	srv := new(todoapp.Server)

	// Пытаемся запустить сервер методом Server и обрабатываем возможную ошибку
	if err := srv.Run(viper.GetString("8080"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
