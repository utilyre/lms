package repository

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type Loan struct {
	bun.BaseModel

	ID         int32 `bun:",pk,autoincrement"`
	UserID     int32
	BookID     int32
	LoanDate   time.Time
	DueDate    time.Time
	ReturnDate sql.NullTime
}
