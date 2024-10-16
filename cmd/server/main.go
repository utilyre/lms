package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/utilyre/lms/internal/repository"
	"golang.org/x/crypto/bcrypt"
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

	if _, err := db.
		NewCreateTable().
		IfNotExists().
		Model((*repository.User)(nil)).
		Exec(context.TODO()); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.GET("/helloworld", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
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

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := repository.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: hash,
			Role:     req.Role,
		}

		_, err = db.NewInsert().Model(&user).Exec(c.Request().Context())
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
