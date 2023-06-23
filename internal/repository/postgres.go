package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	usersTable      = "users"
	tripsTable      = "trips"
	usersTripsTable = "users_trips"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	fmt.Printf("host: %s\nusername: %s\nport: %s\ndbname: %s\nsslmode: %s\n", viper.GetString("db.host"), viper.GetString("db.username"), viper.GetString("db.port"), viper.GetString("db.dbname"), viper.GetString("db.sslmode"))
	for i := 0; i < 10 && err != nil; i++ {
		log.Println("Trying to connect to DB...")
		err = db.Ping()
		time.Sleep(time.Second * 2)
	}

	if err != nil {
		return nil, err
	}

	return db, err
}
