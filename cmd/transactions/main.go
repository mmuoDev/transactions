package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/mmuoDev/transactions/internal/app"
	pg "github.com/mmuoDev/transactions/pkg/postgres"
	"github.com/mmuoDev/transactions/internal/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

//getGRPCAddress returns a grpc address
func getGRPCAddress() string {
	const defaultServerAddress = "127.0.0.1:4444"
	serverAddress, present := os.LookupEnv("PORT")
	if present {
		return serverAddress
	}
	return defaultServerAddress
}

//getPostgreConfig returns postgres conn configs
func getPostgreConfig() pg.Config {
	cfg := pg.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
	return cfg
}

//migratePostgre migrates postgres migrations
func migratePostgre(dbConn *pg.Connector) {
	driver, err := postgres.WithInstance(dbConn.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to connect with error %s", err)
	}
	db.MigrateDB(dbConn.DB, driver, "postgre")
}

func main() {
	dbConn, err := pg.NewConnector(getPostgreConfig())
	if err != nil {
		log.Fatal(err)
	}
	//grpc
	opts := grpc.WithInsecure()
	connOpts := grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: 5 * time.Second,
	})
	addr := getGRPCAddress()
	clientConn, err := grpc.Dial(addr, opts, connOpts)
	if err != nil {
		log.Fatal(err)
	}
	a := app.New(dbConn, clientConn)
	log.Println(fmt.Sprintf("Starting server on port:%s", os.Getenv("APP_PORT")))
	log.Println(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), a.Handler()))
}
