package db

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mmuoDev/transactions/internal"
)

const (
	migrationDirectory = "wallet/internal/db/migration/"
)

//InsertTransactionFunc provides a functionality to insert data into table
type InsertTransactionFunc func(req internal.TransactionDBRequest) (int64, error)

// MigrateDB applies migrations to a database
func MigrateDB(db *sql.DB, driver database.Driver, dbType string) {

	currentDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	migDir := strings.Split(currentDir, os.Getenv("SERVICE_NAME"))[0] + migrationDirectory

	var migrationDir = flag.String("migration files", migDir, "Directory where the migration file exists")
	flag.Parse()

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", *migrationDir), dbType, driver,
	)
	if err != nil {
		log.Fatalf("Error encountered in creating new db instance, %v", err)
	}
	err = m.Up()

	if err != nil && err == migrate.ErrNoChange {
		log.Println("No new migration file")
		return
	} else if err != nil {
		log.Fatalf("Error in migrating with error, %v", err)
	}
	log.Println("Migrated successfully")
}
