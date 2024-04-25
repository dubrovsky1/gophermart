package service

import (
	"context"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/middleware/logger"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/theplant/luhn"
)

type Storager interface {
	Register(context.Context, models.User) (models.UserID, error)
	Login(context.Context, models.User) (models.UserID, error)
	AddOrder(context.Context, models.OrderID, models.UserID) error
	GetOrderList(context.Context, models.UserID) ([]models.Order, error)
	GetBalance(context.Context, models.UserID) (models.Balance, error)
	Withdraw(context.Context, models.OrderID, models.UserID, int) error
	Withdrawals(context.Context, models.UserID) ([]models.Withdraw, error)
	GetOrderAccrualInfo(context.Context, models.OrderID) (models.OrderAccrual, error)
}

type Service struct {
	storage Storager
}

func New(storage Storager) *Service {
	return &Service{storage: storage}
}

func (s *Service) AddOrder(ctx context.Context, orderID models.OrderID, userID models.UserID) error {
	logger.Sugar.Infow("Service Add Log", "orderID", orderID, "userID", userID)

	if !luhn.Valid(int(orderID)) {
		return errs.ErrOrderNum
	}
	if err := s.storage.AddOrder(ctx, orderID, userID); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetOrderList(ctx context.Context, userID models.UserID) ([]models.Order, error) {
	orders, err := s.storage.GetOrderList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Service) GetBalance(ctx context.Context, userID models.UserID) (models.Balance, error) {
	balance, err := s.storage.GetBalance(ctx, userID)
	if err != nil {
		return models.Balance{}, err
	}
	return balance, nil
}

func (s *Service) Withdraw(ctx context.Context, orderID models.OrderID, userID models.UserID, sum int) error {
	logger.Sugar.Infow("Service Withdrawn Log", "orderID", orderID, "userID", userID, "sum", sum)

	if !luhn.Valid(int(orderID)) {
		return errs.ErrOrderNum
	}
	if err := s.storage.Withdraw(ctx, orderID, userID, sum); err != nil {
		return err
	}
	return nil
}

func (s *Service) Withdrawals(ctx context.Context, userID models.UserID) ([]models.Withdraw, error) {
	withdrawals, err := s.storage.Withdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (s *Service) GetOrderAccrualInfo(ctx context.Context, orderID models.OrderID) (models.OrderAccrual, error) {
	orderAccrual, err := s.storage.GetOrderAccrualInfo(ctx, orderID)
	if err != nil {
		return models.OrderAccrual{}, err
	}
	return orderAccrual, nil
}
