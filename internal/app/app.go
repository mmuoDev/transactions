package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mmuoDev/transactions/internal/db"
	pg "github.com/mmuoDev/transactions/pkg/postgres"
)

//App has handlers for this app
type App struct {
	InsertTransactionHandler http.HandlerFunc
}

//Handler returns the main handler for this application
func (a App) Handler() http.HandlerFunc {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/transactions", a.InsertTransactionHandler)
	return http.HandlerFunc(router.ServeHTTP)
}

// Options is a type for application options to modify the app
type Options func(o *OptionalArgs)

// /OptionalArgs optional arguments for this application
type OptionalArgs struct {
	InsertTransaction db.InsertTransactionFunc
}

//New creates a new instance of the App
func New(dbConnector *pg.Connector, options ...Options) App {
	o := OptionalArgs{
		InsertTransaction: db.InsertTransaction(*dbConnector),
	}
	for _, option := range options {
		option(&o)
	}
	insertTransaction := InsertTransactionHandler(o.InsertTransaction)
	return App{
		InsertTransactionHandler: insertTransaction,
	}
}
