package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

func New(db *gorm.DB) *echo.Echo {
	s := &Server{
		db: db,
	}
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/query", s.queryHandler)
	e.GET("/fake", s.fakeHandler)

	return e
}
