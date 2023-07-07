package main

import (
	"net/http"
	"os"

	_ "github.com/itoqsky/InnoCoTravel-backend/docs"
	"github.com/itoqsky/InnoCoTravel-backend/internal/repository"
	"github.com/itoqsky/InnoCoTravel-backend/internal/service"
	"github.com/itoqsky/InnoCoTravel-backend/internal/transport/http"
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
		logrus.Fatalf("error occured while connecting to db: %s", err.Error())
		return
	}

	repos := repository.NewAuthPostgres(db)
	services := service.NewAuthService(repos)
	handlers := http.NewHandler(services)

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
