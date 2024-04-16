package storage

import (
	"context"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/dubrovsky1/gophermart/internal/storage/postgresql"
	"io"
	"time"
)

type Storager interface {
	Register(context.Context, models.User) (string, error)
	Login(context.Context, models.User) (string, error)
	io.Closer
}

type Storage struct {
	Storage Storager
}

func New(connectionString string, connTimeout time.Duration, version int64) (*Storager, error) {
	var db Storager
	var err error

	db, err = postgresql.New(connectionString, connTimeout, version)
	if err != nil {
		return nil, err
	}

	return &db, nil
}
