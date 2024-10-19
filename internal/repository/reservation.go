package repository

import "github.com/uptrace/bun"

type Reservation struct {
	bun.BaseModel

	ID     int32 `bun:",pk,autoincrement"`
	UserID int32
	BookID int32
}
