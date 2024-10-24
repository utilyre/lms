package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
	"github.com/utilyre/lms/internal/model"
)

var (
	ErrBookNotFound = errors.New("book not found")
	ErrBookReserved = errors.New("book reserved")
)

type BookService struct {
	DB bun.IDB
}

type BookCreateParams struct {
	Title  string
	Author string
	ISBN   string
}

func (bs BookService) Create(ctx context.Context, params BookCreateParams) (*model.Book, error) {
	if len(params.Title) == 0 {
		return nil, ValidationError{
			Field: "title",
			Err:   ErrRequired,
		}
	}
	if len(params.Author) == 0 {
		return nil, ValidationError{
			Field: "author",
			Err:   ErrRequired,
		}
	}
	if len(params.ISBN) == 0 {
		return nil, ValidationError{
			Field: "isbn",
			Err:   ErrRequired,
		}
	}

	book := model.Book{
		Title:              params.Title,
		Author:             params.Author,
		ISBN:               params.ISBN,
		AvailabilityStatus: "available",
	}

	_, err := bs.DB.NewInsert().Model(&book).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (bs BookService) GetByID(ctx context.Context, id int32) (*model.Book, error) {
	if id < 1 {
		return nil, ValidationError{
			Field: "id",
			Err:   ErrInvalidID,
		}
	}

	var book model.Book
	if err := bs.DB.
		NewSelect().
		Model(&book).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}

		return nil, err
	}

	return &book, nil
}

type BookUpdateByIDParams struct {
	Title              string
	Author             string
	ISBN               string
	AvailabilityStatus string
}

func (bs BookService) UpdateByID(ctx context.Context, id int32, params BookUpdateByIDParams) (*model.Book, error) {
	if id < 1 {
		return nil, ValidationError{
			Field: "id",
			Err:   ErrInvalidID,
		}
	}
	if len(params.Title) == 0 {
		return nil, ValidationError{
			Field: "title",
			Err:   ErrRequired,
		}
	}
	if len(params.Author) == 0 {
		return nil, ValidationError{
			Field: "author",
			Err:   ErrRequired,
		}
	}
	if len(params.ISBN) == 0 {
		return nil, ValidationError{
			Field: "isbn",
			Err:   ErrRequired,
		}
	}

	book := model.Book{
		Title:              params.Title,
		Author:             params.Author,
		ISBN:               params.ISBN,
		AvailabilityStatus: params.AvailabilityStatus,
	}

	if _, err := bs.DB.
		NewUpdate().
		Model(&book).
		OmitZero().
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return nil, err
	}
	if err := bs.DB.
		NewSelect().
		Model(&book).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return nil, err
	}

	return &book, nil
}

func (bs BookService) DeleteByID(ctx context.Context, id int32) error {
	if id < 1 {
		return ValidationError{
			Field: "id",
			Err:   ErrInvalidID,
		}
	}

	if _, err := bs.DB.
		NewDelete().
		Model((*model.Book)(nil)).
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

type BookBorrowParams struct {
	UserID int32
	BookID int32
}

func (bs BookService) Borrow(ctx context.Context, params BookBorrowParams) (*model.Loan, error) {
	if params.UserID < 1 {
		return nil, ValidationError{
			Field: "user_id",
			Err:   ErrInvalidID,
		}
	}
	if params.BookID < 1 {
		return nil, ValidationError{
			Field: "book_id",
			Err:   ErrInvalidID,
		}
	}

	reserved := true
	var reservation model.Reservation
	if err := bs.DB.
		NewSelect().
		Model(&reservation).
		Where("book_id = ?", params.BookID).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			reserved = false
		} else {
			return nil, err
		}
	}

	if reserved && reservation.UserID != params.UserID {
		return nil, ErrBookReserved
	}

	now := time.Now()
	loan := model.Loan{
		UserID:   params.UserID,
		BookID:   params.BookID,
		LoanDate: now,
		DueDate:  now.Add(14 * 24 * time.Hour),
	}

	if _, err := bs.DB.NewInsert().Model(&loan).Exec(ctx); err != nil {
		return nil, err
	}

	return &loan, nil
}

type BookReturnLoanParams struct {
	LoanID     int32
	ReturnDate time.Time
}

func (bs BookService) ReturnLoan(ctx context.Context, params BookReturnLoanParams) (*model.Loan, error) {
	if params.LoanID < 1 {
		return nil, ValidationError{
			Field: "loan_id",
			Err:   ErrInvalidID,
		}
	}

	loan := model.Loan{
		ID:         params.LoanID,
		ReturnDate: sql.NullTime{Time: params.ReturnDate, Valid: true},
	}

	if _, err := bs.DB.
		NewUpdate().
		Model(&loan).
		OmitZero().
		WherePK().
		Exec(ctx); err != nil {
		return nil, err
	}
	if err := bs.DB.NewSelect().Model(&loan).WherePK().Scan(ctx); err != nil {
		return nil, err
	}

	return &loan, nil
}

type BookReserveParams struct {
	UserID int32
	BookID int32
}

func (bs BookService) Reserve(ctx context.Context, params BookReserveParams) (*model.Reservation, error) {
	if params.UserID < 1 {
		return nil, ValidationError{
			Field: "user_id",
			Err:   ErrInvalidID,
		}
	}
	if params.BookID < 1 {
		return nil, ValidationError{
			Field: "book_id",
			Err:   ErrInvalidID,
		}
	}

	reservation := model.Reservation{
		UserID: params.UserID,
		BookID: params.BookID,
	}

	if _, err := bs.DB.NewInsert().Model(&reservation).Exec(ctx); err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (bs BookService) CancelReservation(ctx context.Context, id int32) error {
	if id < 1 {
		return ValidationError{
			Field: "id",
			Err:   ErrInvalidID,
		}
	}

	if _, err := bs.DB.
		NewDelete().
		Model((*model.Reservation)(nil)).
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return err
	}

	return nil
}
