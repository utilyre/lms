package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utilyre/lms/internal/service"
)

type ReportHandler struct {
	ReportSVC service.ReportService
}

func (rh ReportHandler) GetOverdueLoans(c echo.Context) error {
	loans, err := rh.ReportSVC.GetOverdueLoans(c.Request().Context())
	if err != nil {
		return err
	}

	type RespElem struct {
		ID         int32    `json:"id"`
		UserID     int32    `json:"user_id"`
		BookID     int32    `json:"book_id"`
		LoanDate   DateOnly `json:"loan_date"`
		DueDate    DateOnly `json:"due_date"`
		ReturnDate DateOnly `json:"return_date,omitempty"`
	}
	resp := make([]RespElem, len(loans))
	for i, loan := range loans {
		resp[i].ID = loan.ID
		resp[i].UserID = loan.UserID
		resp[i].BookID = loan.BookID
		resp[i].LoanDate = DateOnly{loan.LoanDate}
		resp[i].DueDate = DateOnly{loan.DueDate}
		resp[i].ReturnDate = DateOnly{loan.ReturnDate.Time}
	}
	return c.JSON(http.StatusOK, resp)
}

func (rh ReportHandler) GetPopularBooks(c echo.Context) error {
	results, err := rh.ReportSVC.GetPopularBooks(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, results)
}
