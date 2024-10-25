package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utilyre/lms/internal/service"
)

type UserHandler struct {
	UserSVC service.UserService
}

func (uh UserHandler) Delete(c echo.Context) error {
	type Req struct {
		ID int32 `param:"id"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := uh.UserSVC.DeleteByID(c.Request().Context(), req.ID); err != nil {
		var validationErr service.ValidationError
		if errors.As(err, &validationErr) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]any{
				"type":    "validation",
				"message": validationErr.Error(),
			})
		}

		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "User deleted successfully",
	})
}

func (uh UserHandler) Update(c echo.Context) error {
	type Req struct {
		ID    int32  `param:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := uh.UserSVC.UpdateByID(c.Request().Context(), req.ID, service.UserUpdateByIDParams{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	})
	if err != nil {
		var validationErr service.ValidationError
		if errors.As(err, &validationErr) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]any{
				"type":    "validation",
				"message": validationErr.Error(),
			})
		}

		return err
	}

	type Resp struct {
		ID    int32  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	return c.JSON(http.StatusCreated, Resp{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	})
}

func (uh UserHandler) Get(c echo.Context) error {
	type Req struct {
		ID int32 `param:"id"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := uh.UserSVC.GetByID(c.Request().Context(), req.ID)
	if err != nil {
		var validationErr service.ValidationError
		if errors.As(err, &validationErr) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]any{
				"type":    "validation",
				"message": validationErr.Error(),
			})
		}
		if errors.Is(err, service.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"type":    "resource",
				"message": "user not found",
			})
		}

		return err
	}

	type Resp struct {
		ID    int32  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	return c.JSON(http.StatusCreated, Resp{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	})
}

func (uh UserHandler) Create(c echo.Context) error {
	type Req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := uh.UserSVC.Create(c.Request().Context(), service.UserCreateParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: []byte(req.Password),
		Role:     req.Role,
	})
	if err != nil {
		var validationErr service.ValidationError
		if errors.As(err, &validationErr) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]any{
				"type":    "validation",
				"message": validationErr.Error(),
			})
		}
		if errors.Is(err, service.ErrUserDup) {
			return c.JSON(http.StatusConflict, map[string]any{
				"type":    "logic",
				"message": "user already exists",
			})
		}

		return err
	}

	type Resp struct {
		ID    int32  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	return c.JSON(http.StatusCreated, Resp{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	})
}
