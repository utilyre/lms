package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/utilyre/lms/internal/repository"
	"github.com/utilyre/lms/internal/service"
)

var listenPort string

func init() {
	flag.StringVar(&listenPort, "port", "8080", "specify port to listen on")
	flag.Parse()
}

func main() {
	log.Printf("Connecting to %s\n", os.Getenv("DB_URL"))
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DB_URL"))))
	db := bun.NewDB(sqldb, pgdialect.New())

	userSVC := service.UserService{DB: db}

	if _, err := db.
		NewCreateTable().
		IfNotExists().
		Model((*repository.User)(nil)).
		Exec(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if _, err := db.
		NewCreateTable().
		IfNotExists().
		Model((*repository.Book)(nil)).
		Exec(context.TODO()); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/helloworld", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	e.DELETE("/books/:id", func(c echo.Context) error {
		type Req struct {
			ID int32 `param:"id"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			return err
		}

		if _, err := db.
			NewDelete().
			Model((*repository.Book)(nil)).
			Where("id = ?", req.ID).
			Exec(c.Request().Context()); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"message": "Book deleted successfully",
		})
	})

	e.PUT("/books/:id", func(c echo.Context) error {
		type Req struct {
			ID                 int32  `param:"id"`
			Title              string `json:"title"`
			Author             string `json:"author"`
			ISBN               string `json:"isbn"`
			AvailabilityStatus string `json:"availability_status"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			return err
		}

		book := repository.Book{
			ID:                 req.ID,
			Title:              req.Title,
			Author:             req.Author,
			ISBN:               req.ISBN,
			AvailabilityStatus: req.AvailabilityStatus,
		}

		// TODO: make these a singular tx
		if _, err := db.
			NewUpdate().
			Model(&book).
			OmitZero().
			WherePK().
			Exec(c.Request().Context()); err != nil {
			return err
		}
		if err := db.
			NewSelect().
			Model(&book).
			WherePK().
			Scan(c.Request().Context()); err != nil {
			return err
		}

		type Resp struct {
			ID                 int32  `json:"id"`
			Title              string `json:"title"`
			Author             string `json:"author"`
			ISBN               string `json:"isbn"`
			AvailabilityStatus string `json:"availability_status"`
		}
		return c.JSON(http.StatusCreated, Resp{
			ID:                 book.ID,
			Title:              book.Title,
			Author:             book.Author,
			ISBN:               book.ISBN,
			AvailabilityStatus: book.AvailabilityStatus,
		})
	})

	e.GET("/books/:id", func(c echo.Context) error {
		type Req struct {
			ID int32 `param:"id"`
		}

		var req Req
		if err := c.Bind(&req); err != nil {
			return err
		}

		var book repository.Book
		if err := db.
			NewSelect().
			Model(&book).
			Where("id = ?", req.ID).
			Scan(c.Request().Context()); err != nil {
			return err
		}

		type Resp struct {
			ID                 int32  `json:"id"`
			Title              string `json:"title"`
			Author             string `json:"author"`
			ISBN               string `json:"isbn"`
			AvailabilityStatus string `json:"availability_status"`
		}
		return c.JSON(http.StatusCreated, Resp{
			ID:                 book.ID,
			Title:              book.Title,
			Author:             book.Author,
			ISBN:               book.ISBN,
			AvailabilityStatus: book.AvailabilityStatus,
		})
	})

	e.POST("/books", func(c echo.Context) error {
		type Req struct {
			Title  string `json:"title"`
			Author string `json:"author"`
			ISBN   string `json:"isbn"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return err
		}

		book := repository.Book{
			Title:              req.Title,
			Author:             req.Author,
			ISBN:               req.ISBN,
			AvailabilityStatus: "available",
		}

		_, err := db.NewInsert().Model(&book).Exec(c.Request().Context())
		if err != nil {
			return err
		}

		type Resp struct {
			ID                 int32  `json:"id"`
			Title              string `json:"title"`
			Author             string `json:"author"`
			ISBN               string `json:"isbn"`
			AvailabilityStatus string `json:"availability_status"`
		}
		return c.JSON(http.StatusCreated, Resp{
			ID:                 book.ID,
			Title:              book.Title,
			Author:             book.Author,
			ISBN:               book.ISBN,
			AvailabilityStatus: book.AvailabilityStatus,
		})
	})

	e.DELETE("/users/:id", func(c echo.Context) error {
		type Req struct {
			ID int32 `param:"id"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := userSVC.DeleteByID(c.Request().Context(), req.ID); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]any{
			"message": "User deleted successfully",
		})
	})

	e.PUT("/users/:id", func(c echo.Context) error {
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

		user, err := userSVC.UpdateByID(c.Request().Context(), req.ID, service.UserUpdateByIDParams{
			Name:  req.Name,
			Email: req.Email,
			Role:  req.Role,
		})
		if err != nil {
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
	})

	e.GET("/users/:id", func(c echo.Context) error {
		type Req struct {
			ID int32 `param:"id"`
		}
		var req Req
		if err := c.Bind(&req); err != nil {
			return err
		}

		user, err := userSVC.GetByID(c.Request().Context(), req.ID)
		if err != nil {
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
	})

	e.POST("/users", func(c echo.Context) error {
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

		user, err := userSVC.Create(c.Request().Context(), service.UserCreateParams{
			Name:     req.Name,
			Email:    req.Email,
			Password: []byte(req.Password),
			Role:     req.Role,
		})
		if err != nil {
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
	})

	e.Logger.Fatal(e.Start(":" + listenPort))
}
