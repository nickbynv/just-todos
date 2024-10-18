package main

import (
	"context"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"

	todo "just-todos"
	"just-todos/pkg/handler"
	"just-todos/pkg/repository"
	"just-todos/pkg/service"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Just todos API
// @version 1.0
// @description API server for todo-list app

// @host localhost:8800
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	err := initConfig()
	if err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	err = godotenv.Load()
	if err != nil {
		logrus.Fatalf("error loading .env file: %s", err.Error())
	}

	db, err := repository.NewPostgres(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error initializing DB: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	svc := service.NewService(repos)
	hdl := handler.NewHandler(svc)

	server := new(todo.Server)

	go func() {
		err := server.Run(
			viper.GetString("port"),
			hdl.InitRoutes(),
		)
		if err != nil {
			logrus.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	logrus.Println("app started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logrus.Println("app shutting down")

	err = server.Shutdown(context.Background())
	if err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	err = db.Close()
	if err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
