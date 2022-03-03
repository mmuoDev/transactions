package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mmuoDev/transactions/internal/app"
	pg "github.com/mmuoDev/transactions/pkg/postgres"
)

func main() {
	cfg := pg.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
	dbConn, err := pg.NewConnector(cfg)
	if err != nil {
		log.Fatal(err)
	}
	a := app.New(dbConn)
	log.Println(fmt.Sprintf("Starting server on port:%s", os.Getenv("APP_PORT")))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), a.Handler()))
}
