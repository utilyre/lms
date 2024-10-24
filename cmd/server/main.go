package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/utilyre/lms/internal/handler"
	"github.com/utilyre/lms/internal/service"
)

var listenPort string

func init() {
	flag.StringVar(&listenPort, "port", "8080", "specify port to listen on")
	flag.Parse()
}

func main() {
	log.Printf("Connecting to database: %s\n", os.Getenv("DB_URL"))
	db := bun.NewDB(
		sql.OpenDB(
			pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DB_URL"))),
		),
		pgdialect.New(),
	)
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	log.Printf("Connecting to cache: %s\n", os.Getenv("CACHE_URL"))
	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("CACHE_URL")})

	userSVC := service.UserService{DB: db}
	bookSVC := service.BookService{DB: db}
	reportSVC := service.ReportService{DB: db, RDB: rdb}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	setupRoutes(
		e,
		handler.UserHandler{UserSVC: userSVC},
		handler.BookHandler{BookSVC: bookSVC},
		handler.ReportHandler{ReportSVC: reportSVC},
	)

	log.Fatal(e.Start(":" + listenPort))
}
