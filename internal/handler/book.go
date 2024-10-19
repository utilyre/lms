package handler

import (
	"encoding/json"
	"net/http"
	"time"

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

type DateOnly struct{ time.Time }

func (do DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(do.Format(time.DateOnly))
}

func (do *DateOnly) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	t, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return err
	}

	do.Time = t
	return nil
}

func (bh BookHandler) Borrow(c echo.Context) error {
	type Req struct {
		UserID int32 `json:"user_id"`
		BookID int32 `json:"book_id"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	loan, err := bh.BookSVC.Borrow(c.Request().Context(), service.BookBorrowParams{
		UserID: req.UserID,
		BookID: req.BookID,
	})
	if err != nil {
		return err
	}

	type Resp struct {
		ID       int32    `json:"id"`
		UserID   int32    `json:"user_id"`
		BookID   int32    `json:"book_id"`
		LoanDate DateOnly `json:"loan_date"`
		DueDate  DateOnly `json:"due_date"`
	}
	return c.JSON(http.StatusCreated, Resp{
		ID:       loan.ID,
		UserID:   loan.UserID,
		BookID:   loan.BookID,
		LoanDate: DateOnly{loan.LoanDate},
		DueDate:  DateOnly{loan.DueDate},
	})
}

func (bh BookHandler) ReturnLoan(c echo.Context) error {
	type Req struct {
		ID         int32    `param:"id"`
		ReturnDate DateOnly `json:"return_date"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	loan, err := bh.BookSVC.ReturnLoan(c.Request().Context(), service.BookReturnLoanParams{
		LoanID:     req.ID,
		ReturnDate: req.ReturnDate.Time,
	})
	if err != nil {
		return err
	}

	type Resp struct {
		ID         int32    `json:"id"`
		UserID     int32    `json:"user_id"`
		BookID     int32    `json:"book_id"`
		LoanDate   DateOnly `json:"loan_date"`
		DueDate    DateOnly `json:"due_date"`
		ReturnDate DateOnly `json:"return_date"`
	}
	return c.JSON(http.StatusOK, Resp{
		ID:         loan.ID,
		UserID:     loan.UserID,
		BookID:     loan.BookID,
		LoanDate:   DateOnly{loan.LoanDate},
		DueDate:    DateOnly{loan.DueDate},
		ReturnDate: DateOnly{loan.ReturnDate.Time},
	})
}

func (bh BookHandler) Reserve(c echo.Context) error {
	type Req struct {
		UserID int32 `json:"user_id"`
		BookID int32 `json:"book_id"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	reservation, err := bh.BookSVC.Reserve(c.Request().Context(), service.BookReserveParams{
		UserID: req.UserID,
		BookID: req.BookID,
	})
	if err != nil {
		return err
	}

	type Resp struct {
		ID     int32 `json:"id"`
		UserID int32 `json:"user_id"`
		BookID int32 `json:"book_id"`
	}
	return c.JSON(http.StatusCreated, Resp{
		ID:     reservation.ID,
		UserID: reservation.UserID,
		BookID: reservation.BookID,
	})
}

func (bh BookHandler) CancelReservation(c echo.Context) error {
	type Req struct {
		ID int32 `param:"id"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := bh.BookSVC.CancelReservation(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Reservation canceled successfully",
	})
}
