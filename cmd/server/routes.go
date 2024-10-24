package main

import (
	"net/http"

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
	e.GET("/reports/popular-books", reportHandler.GetPopularBooks)
	e.GET("/reports/user-activity/:id", reportHandler.GetUserActivity)
}
