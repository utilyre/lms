package service

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/utilyre/lms/internal/repository"
)

type ReportService struct {
	DB bun.IDB
}

func (rs ReportService) GetOverdueLoans(ctx context.Context) ([]repository.Loan, error) {
	var loans []repository.Loan

	if err := rs.DB.
		NewSelect().
		Model(&loans).
		// WHERE return_date IS NULL AND NOW() > due_date OR return_date > due_date
		Where("return_date IS NULL").
		WhereGroup("AND", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("NOW() > due_date")
		}).
		WhereOr("return_date > due_date").
		Scan(ctx); err != nil {
		return nil, err
	}

	return loans, nil
}

// func (rs ReportService) GetPopularBooks(ctx context.Context) ([]repository.Book, error) {
// 	var books []repository.Book

// 	rs.DB.NewSelect().Model(&books).Where()
// }
