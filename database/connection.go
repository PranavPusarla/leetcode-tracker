package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
)

import _ "github.com/lib/pq"

import "github.com/computer-geek64/leetcode-tracker/config"

func Connect(conf config.Config) *sql.DB {
	var password string
	if conf.Database.Password == nil {
		password = os.Getenv("POSTGRES_PASSWORD")
	} else {
		password = *conf.Database.Password
	}

	var postgresInfo = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.Database.Username, password, conf.Database.Host, conf.Database.Port, conf.Database.Name)
	var db, err = sql.Open("postgres", postgresInfo)
	if err != nil {
		slog.Error("Failed to create database connection")
		panic(err)
	}
	if err := db.Ping(); err != nil {
		slog.Error("Failed to connect to database")
		panic(err)
	}
	return db
}
