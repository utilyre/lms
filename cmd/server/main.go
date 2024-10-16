package main

import (
	"database/sql"
	"flag"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var listenPort string

func init() {
	flag.StringVar(&listenPort, "port", "8080", "specify port to listen on")
	flag.Parse()
}

func main() {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DB_URL"))))
	db := bun.NewDB(sqldb, pgdialect.New())
	_ = db

	e := echo.New()

	e.GET("/helloworld", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	e.Logger.Fatal(e.Start(":" + listenPort))
}
