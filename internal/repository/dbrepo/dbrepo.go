package dbrepo

import (
	"database/sql"

	"github.com/prashant9154/Booking_System/internal/config"
	"github.com/prashant9154/Booking_System/internal/repository"
)

type postgressDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgressDBRepo{
		App: a,
		DB:  conn,
	}
}
