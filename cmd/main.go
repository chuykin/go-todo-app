package main

import (
	"context"
	todo "github.com/IncubusX/go-todo-app"
	"github.com/IncubusX/go-todo-app/pkg/handler"
	"github.com/IncubusX/go-todo-app/pkg/repository"
	"github.com/IncubusX/go-todo-app/pkg/service"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

const serverClosed = "http: Server closed"

//	@title			Todo App API
//	@version		1.0
//	@description	API Server for TodoList Application

//	@host		localhost:8000
//	@BasePath	/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @authorizationurl			http://localhost:8000/sign-in
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка при чтении конфигурации: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Ошибка при чтении переменных окружения:%s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("Ошибка при инициализации БД: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil && err.Error() != serverClosed {
			logrus.Fatalf("Ошибка при запуске HTTP сервера: %s", err.Error())
		}
	}()

	logrus.Println("HTTP Сервер запущен!")

	gracefulShutdown(srv, db)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func gracefulShutdown(srv *todo.Server, db *sqlx.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("Ошибка во время остановки HTTP Сервера: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Fatalf("Ошибка во время остановки БД: %s", err.Error())
	}

	logrus.Println("HTTP Сервер завершил работу!")
}
