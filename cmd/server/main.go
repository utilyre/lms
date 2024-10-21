package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/utilyre/lms/internal/handler"
	"github.com/utilyre/lms/internal/repository"
	"github.com/utilyre/lms/internal/service"
)

var listenPort string

func init() {
	flag.StringVar(&listenPort, "port", "8080", "specify port to listen on")
	flag.Parse()
}

func main() {
	log.Printf("Connecting to %s\n", os.Getenv("DB_URL"))
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DB_URL"))))
	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if _, err := db.
		NewCreateTable().
		IfNotExists().
		Model((*repository.User)(nil)).
		Exec(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if _, err := db.
		NewCreateTable().
		IfNotExists().
		Model((*repository.Book)(nil)).
		Exec(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if _, err := db.
		NewCreateTable().
		IfNotExists().
		Model((*repository.Loan)(nil)).
		Exec(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if _, err := db.
		NewCreateTable().
		IfNotExists().
		Model((*repository.Reservation)(nil)).
		Exec(context.TODO()); err != nil {
		log.Fatal(err)
	}

	userSVC := service.UserService{DB: db}
	bookSVC := service.BookService{DB: db}
	reportSVC := service.ReportService{DB: db}

	userHandler := handler.UserHandler{UserSVC: userSVC}
	bookHandler := handler.BookHandler{BookSVC: bookSVC}
	reportHandler := handler.ReportHandler{ReportSVC: reportSVC}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/helloworld", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	e.POST("/users", userHandler.Create)
	e.PUT("/users/:id", userHandler.Update)
	e.GET("/users/:id", userHandler.Get)
	e.DELETE("/users/:id", userHandler.Delete)

	e.DELETE("/books/:id", bookHandler.Delete)
	e.PUT("/books/:id", bookHandler.Update)
	e.GET("/books/:id", bookHandler.Get)
	e.POST("/books", bookHandler.Create)
	e.POST("/loans", bookHandler.Borrow)
	e.PUT("/loans/:id", bookHandler.ReturnLoan)
	e.POST("/reservations", bookHandler.Reserve)
	e.DELETE("/reservations/:id", bookHandler.CancelReservation)

	e.GET("/reports/overdue-loans", reportHandler.GetOverdueLoans)

	e.Logger.Fatal(e.Start(":" + listenPort))
}
