package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/utilyre/lms/internal/model"
)

type ReportService struct {
	DB  bun.IDB
	RDB *redis.Client
}

func (rs ReportService) GetOverdueLoans(ctx context.Context) ([]model.Loan, error) {
	var loans []model.Loan

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

type ReportGetPopularBooksResult struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Borrows int    `json:"borrows"`
}

var (
	keyPopularBooks = "popular-books"
)

func (rs ReportService) GetPopularBooks(ctx context.Context) ([]ReportGetPopularBooksResult, error) {
	exists, err := rs.RDB.Exists(ctx, keyPopularBooks).Result()
	if err != nil {
		return nil, err
	}
	if exists != 0 {
		objects, err := rs.RDB.LRange(ctx, keyPopularBooks, 0, -1).Result()
		if err != nil {
			return nil, err
		}

		results := make([]ReportGetPopularBooksResult, len(objects))
		for i, obj := range objects {
			if err := json.Unmarshal([]byte(obj), &results[i]); err != nil {
				return nil, err
			}
		}

		log.Println("Used cache to respond popular books")
		return results, nil
	}

	var results []ReportGetPopularBooksResult
	if err := rs.DB.
		NewSelect().
		Model((*model.Book)(nil)).
		ColumnExpr("book.id id, book.title title, COUNT(*) borrows").
		Join("JOIN loans loan ON loan.book_id = book.id").
		Group("book.id").
		Scan(ctx, &results); err != nil {
		return nil, err
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		raw := make([]any, len(results))
		for i, result := range results {
			data, err := json.Marshal(result)
			if err != nil {
				log.Println("Failed to marshal popular book:", err)
				return
			}

			raw[i] = data
		}

		if err := rs.RDB.LPush(ctx, keyPopularBooks, raw...).Err(); err != nil {
			log.Println("Failed to push popular books to cache:", err)
			return
		}
		if err := rs.RDB.Expire(ctx, keyPopularBooks, 24*time.Hour).Err(); err != nil {
			log.Println("Failed to set expiration time on popular books:", err)
			return
		}
	}()

	return results, nil
}

func (rs ReportService) GetUserActivity(ctx context.Context, id int32) ([]model.Loan, error) {
	if id < 1 {
		return nil, ValidationError{
			Field: "id",
			Err:   ErrInvalidID,
		}
	}

	var loans []model.Loan

	if err := rs.DB.
		NewSelect().
		Column("id", "book_id", "loan_date", "due_date", "return_date").
		Model(&loans).
		Where("user_id = ?", id).
		Scan(ctx); err != nil {
		return nil, err
	}

	return loans, nil
}
