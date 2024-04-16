package service

import (
	"context"
	"github.com/dubrovsky1/gophermart/internal/models"
)

type Storager interface {
	Register(context.Context, models.User) (string, error)
	Login(context.Context, models.User) (string, error)
}

type Service struct {
	storage Storager
}

func New(storage Storager) *Service {
	return &Service{storage: storage}
}
