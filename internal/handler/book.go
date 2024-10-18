package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utilyre/lms/internal/service"
)

type BookHandler struct {
	BookSVC service.BookService
}

func (bh BookHandler) Delete(c echo.Context) error {
	type Req struct {
		ID int32 `param:"id"`
	}

	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	err := bh.BookSVC.DeleteByID(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Book deleted successfully",
	})
}

func (bh BookHandler) Update(c echo.Context) error {
	type Req struct {
		ID                 int32  `param:"id"`
		Title              string `json:"title"`
		Author             string `json:"author"`
		ISBN               string `json:"isbn"`
		AvailabilityStatus string `json:"availability_status"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	book, err := bh.BookSVC.UpdateByID(c.Request().Context(), req.ID, service.BookUpdateByIDParams{
		Title:              req.Title,
		Author:             req.Author,
		ISBN:               req.ISBN,
		AvailabilityStatus: req.AvailabilityStatus,
	})
	if err != nil {
		return err
	}

	type Resp struct {
		ID                 int32  `json:"id"`
		Title              string `json:"title"`
		Author             string `json:"author"`
		ISBN               string `json:"isbn"`
		AvailabilityStatus string `json:"availability_status"`
	}
	return c.JSON(http.StatusCreated, Resp{
		ID:                 book.ID,
		Title:              book.Title,
		Author:             book.Author,
		ISBN:               book.ISBN,
		AvailabilityStatus: book.AvailabilityStatus,
	})
}

func (bh BookHandler) Get(c echo.Context) error {
	type Req struct {
		ID int32 `param:"id"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	book, err := bh.BookSVC.GetByID(c.Request().Context(), req.ID)
	if err != nil {
		return err
	}

	type Resp struct {
		ID                 int32  `json:"id"`
		Title              string `json:"title"`
		Author             string `json:"author"`
		ISBN               string `json:"isbn"`
		AvailabilityStatus string `json:"availability_status"`
	}
	return c.JSON(http.StatusCreated, Resp{
		ID:                 book.ID,
		Title:              book.Title,
		Author:             book.Author,
		ISBN:               book.ISBN,
		AvailabilityStatus: book.AvailabilityStatus,
	})
}

func (bh BookHandler) Create(c echo.Context) error {
	type Req struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		ISBN   string `json:"isbn"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	book, err := bh.BookSVC.Create(c.Request().Context(), service.BookCreateParams{
		Title:  req.Title,
		Author: req.Author,
		ISBN:   req.ISBN,
	})
	if err != nil {
		return err
	}

	type Resp struct {
		ID                 int32  `json:"id"`
		Title              string `json:"title"`
		Author             string `json:"author"`
		ISBN               string `json:"isbn"`
		AvailabilityStatus string `json:"availability_status"`
	}
	return c.JSON(http.StatusCreated, Resp{
		ID:                 book.ID,
		Title:              book.Title,
		Author:             book.Author,
		ISBN:               book.ISBN,
		AvailabilityStatus: book.AvailabilityStatus,
	})
}
