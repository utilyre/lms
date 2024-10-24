package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/utilyre/lms/internal/handler"
)

func setupRoutes(
	e *echo.Echo,
	userHandler handler.UserHandler,
	bookHandler handler.BookHandler,
	reportHandler handler.ReportHandler,
) {
	e.GET("/helloworld", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	apiV1 := e.Group("/api/v1")

	users := apiV1.Group("/users")
	users.POST("/", userHandler.Create)
	users.PUT("/:id", userHandler.Update)
	users.GET("/:id", userHandler.Get)
	users.DELETE("/:id", userHandler.Delete)

	books := apiV1.Group("/books")
	books.DELETE("/:id", bookHandler.Delete)
	books.PUT("/:id", bookHandler.Update)
	books.GET("/:id", bookHandler.Get)
	books.POST("/", bookHandler.Create)

	loans := apiV1.Group("/loans")
	loans.POST("/", bookHandler.Borrow)
	loans.PUT("/:id", bookHandler.ReturnLoan)

	reservations := apiV1.Group("/reservations")
	reservations.POST("/", bookHandler.Reserve)
	reservations.DELETE("/:id", bookHandler.CancelReservation)

	reports := apiV1.Group("/reports")
	reports.GET("/overdue-loans", reportHandler.GetOverdueLoans)
	reports.GET("/popular-books", reportHandler.GetPopularBooks)
	reports.GET("/user-activity/:id", reportHandler.GetUserActivity)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(e.Routes()); err != nil {
		log.Fatal(err)
	}
}
