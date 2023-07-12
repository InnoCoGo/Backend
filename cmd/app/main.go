package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/itoqsky/InnoCoTravel-backend/docs"
	"github.com/itoqsky/InnoCoTravel-backend/internal/kafka"
	"github.com/itoqsky/InnoCoTravel-backend/internal/repository"
	"github.com/itoqsky/InnoCoTravel-backend/internal/server"
	"github.com/itoqsky/InnoCoTravel-backend/internal/service"
	transport "github.com/itoqsky/InnoCoTravel-backend/internal/transport/http"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error occured while initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			logrus.Errorf("not found .env file: %s", err.Error())
		} else {
			logrus.Fatalf("error occured while loading .env file: %s", err.Error())
		}
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error occured while connecting  to db: %s", err.Error())
		return
	}

	kafka.InitProducer(os.Getenv("KAFKA_TOPIC"), os.Getenv("KAFKA_HOSTS"))
	kafka.InitConsumer(os.Getenv("KAFKA_HOSTS"))

	hub := server.NewHub()
	go hub.Run()

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := transport.NewHandler(services, hub)

	srv := server.NewServer(viper.GetString("port"), handlers.Init())

	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	// Gracefull Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Printf("\nServer shutting down...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("failed to shut down the server: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured while closing the database: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
