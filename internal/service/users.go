package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/uptrace/bun"
	"github.com/utilyre/lms/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrRequired     = errors.New("required")
	ErrTooShort     = errors.New("too short")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidID    = errors.New("invalid id")
	ErrUserNotFound = errors.New("user not found")
)

type ValidationError struct {
	Field string
	Err   error
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %v", ve.Field, ve.Err)
}

func (ve ValidationError) Unwrap() error {
	return ve.Err
}

type UserService struct {
	DB bun.IDB
}

type UserCreateParams struct {
	Name     string
	Email    string
	Password []byte
	Role     string
}

var reEmail = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)

func (us UserService) Create(ctx context.Context, params UserCreateParams) (*repository.User, error) {
	if len(params.Name) == 0 {
		return nil, ValidationError{
			Field: "name",
			Err:   ErrRequired,
		}
	}
	if len(params.Email) == 0 {
		return nil, ValidationError{
			Field: "email",
			Err:   ErrRequired,
		}
	}
	if !reEmail.Match([]byte(params.Email)) {
		return nil, ValidationError{
			Field: "email",
			Err:   ErrInvalidEmail,
		}
	}
	if len(params.Password) < 3 {
		return nil, ValidationError{
			Field: "password",
			Err:   ErrTooShort,
		}
	}

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
	if id < 1 {
		return nil, ValidationError{
			Field: "id",
			Err:   ErrInvalidID,
		}
	}

	var user repository.User
	if err := us.DB.
		NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

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
	if id < 1 {
		return nil, ValidationError{
			Field: "id",
			Err:   ErrInvalidID,
		}
	}
	if len(params.Name) == 0 {
		return nil, ValidationError{
			Field: "name",
			Err:   ErrRequired,
		}
	}
	if len(params.Email) == 0 {
		return nil, ValidationError{
			Field: "email",
			Err:   ErrRequired,
		}
	}
	if !reEmail.Match([]byte(params.Email)) {
		return nil, ValidationError{
			Field: "email",
			Err:   ErrInvalidEmail,
		}
	}

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
	if id < 1 {
		return ValidationError{
			Field: "id",
			Err:   ErrInvalidID,
		}
	}

	if _, err := us.DB.
		NewDelete().
		Model((*repository.User)(nil)).
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return err
	}

	return nil
}
