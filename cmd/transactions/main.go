package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mmuoDev/transactions/internal/app"
	pg "github.com/mmuoDev/transactions/pkg/postgres"
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
