package service

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/utilyre/lms/internal/repository"
)

type BookService struct {
	DB bun.IDB
}

type BookCreateParams struct {
	Title  string
	Author string
	ISBN   string
}

func (bs BookService) Create(ctx context.Context, params BookCreateParams) (*repository.Book, error) {
	book := repository.Book{
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

func (bs BookService) GetByID(ctx context.Context, id int32) (*repository.Book, error) {
	var book repository.Book
	if err := bs.DB.
		NewSelect().
		Model(&book).
		Where("id = ?", id).
		Scan(ctx); err != nil {
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

func (bs BookService) UpdateByID(ctx context.Context, id int32, params BookUpdateByIDParams) (*repository.Book, error) {
	book := repository.Book{
		Title:              params.Title,
		Author:             params.Author,
		ISBN:               params.ISBN,
		AvailabilityStatus: params.AvailabilityStatus,
	}

	// TODO: make these a singular tx
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
	if _, err := bs.DB.
		NewDelete().
		Model((*repository.Book)(nil)).
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return err
	}

	return nil
}
