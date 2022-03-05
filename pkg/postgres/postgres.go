package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

//Config presents data required to establish a connection to postgres
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

//Connector wraps a connection to postgres
type Connector struct {
	DB *sql.DB
}

//NewConnector creates a connector to postgres
func NewConnector(cfg Config) (*Connector, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return &Connector{}, err
	}
	return &Connector{DB: db}, nil
}
