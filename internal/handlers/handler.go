package handlers

import "github.com/dubrovsky1/gophermart/internal/service"

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}
