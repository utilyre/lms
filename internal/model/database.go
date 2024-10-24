package model

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel

	ID       int32 `bun:",pk,autoincrement"`
	Name     string
	Email    string `bun:",unique"`
	Password []byte
	Role     string
}

type Book struct {
	bun.BaseModel

	ID                 int32 `bun:",pk,autoincrement"`
	Title              string
	Author             string
	ISBN               string
	AvailabilityStatus string
}

type Loan struct {
	bun.BaseModel

	ID         int32 `bun:",pk,autoincrement"`
	UserID     int32
	BookID     int32
	LoanDate   time.Time
	DueDate    time.Time
	ReturnDate sql.NullTime
}

type Reservation struct {
	bun.BaseModel

	ID     int32 `bun:",pk,autoincrement"`
	UserID int32
	BookID int32
}
