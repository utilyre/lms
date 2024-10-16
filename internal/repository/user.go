package repository

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel

	ID       int32 `bun:",pk,autoincrement"`
	Name     string
	Email    string
	Password []byte
	Role     string
}
