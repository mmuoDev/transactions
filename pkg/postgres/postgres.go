package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
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
	db *sql.DB
}

//NewConnector creates a connector to postgres
func NewConnector(cfg Config) (*Connector, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/test", cfg.Username, cfg.Password)
	// connStr := fmt.Sprintf("host=%s port=%s user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return &Connector{}, err
	}
	return &Connector{db: db}, nil
}

//Insert inserts a new record in a table
//Expects a table name and a map of columns and its respective values
func (c *Connector) Insert(tableName string, values map[string]interface{}) (int64, error) {
	columns := make([]string, 0, len(values))
	rows := []interface{}{}
	for k, v := range values {
		columns = append(columns, k)
		rows = append(rows, v)
	}
	cols := strings.Join(columns, ",")
	ll := len(columns)
	pp := make([]string, 0, ll)
	for i := 1; i <= ll; i++ {
		pp = append(pp, "?")
	}
	pStr := strings.Join(pp, ",")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", tableName, cols, pStr)
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to prepare statement")
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, rows...)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "Unable to get last inserted ID")
	}
	return id, nil
}
