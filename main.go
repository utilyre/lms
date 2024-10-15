package main

import (
	"database/sql"
	"net/http"
	"flag"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var listenPort string

func init() {
	flag.StringVar(&listenPort, "port", "8080", "specify port to listen on")
	flag.Parse()
}

func main() {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	_ = db

	e := echo.New()

	e.GET("/helloworld", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	e.Logger.Fatal(e.Start(":"+listenPort))
}
