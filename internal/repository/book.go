package repository

import "github.com/uptrace/bun"

type Book struct {
	bun.BaseModel

	ID                 int32 `bun:",pk,autoincrement"`
	Title              string
	Author             string
	ISBN               string
	AvailabilityStatus string
}
