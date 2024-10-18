package service

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/utilyre/lms/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB bun.IDB
}

type UserCreateParams struct {
	Name     string
	Email    string
	Password []byte
	Role     string
}

func (us UserService) Create(ctx context.Context, params UserCreateParams) (*repository.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := repository.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: hash,
		Role:     params.Role,
	}

	_, err = us.DB.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us UserService) GetByID(ctx context.Context, id int32) (*repository.User, error) {
	var user repository.User
	if err := us.DB.
		NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

type UserUpdateByIDParams struct {
	Name  string
	Email string
	Role  string
}

func (us UserService) UpdateByID(ctx context.Context, id int32, params UserUpdateByIDParams) (*repository.User, error) {
	user := repository.User{
		Name:  params.Name,
		Email: params.Email,
		Role:  params.Role,
	}

	if _, err := us.DB.
		NewUpdate().
		Model(&user).
		OmitZero().
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return nil, err
	}
	if err := us.DB.
		NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (us UserService) DeleteByID(ctx context.Context, id int32) error {
	if _, err := us.DB.
		NewDelete().
		Model((*repository.User)(nil)).
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return err
	}

	return nil
}
