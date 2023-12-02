package api

import (
	"logistics_status_tracking/internal/service/pkg/cachestore"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Server struct {
	db         *gorm.DB
	cacheStore *cachestore.RedisStore
}

func New(db *gorm.DB, store *cachestore.RedisStore) *echo.Echo {
	s := &Server{
		db:         db,
		cacheStore: store,
	}
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/query", s.queryHandler)
	e.GET("/fake", s.fakeHandler)

	return e
}
