package postgresql

import (
	"context"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/models"
)

func (s *Storage) AddOrder(ctx context.Context, orderID models.OrderID, userID models.UserID) error {
	var u models.UserID

	//проверка наличия заказа
	row := s.pool.QueryRow(ctx, "select o.userid from orders o where o.orderid = $1;", orderID)
	err := row.Scan(&u)
	if err == nil {
		if u == userID {
			return errs.ErrOrderAlreadyLoadThisUser
		}
		return errs.ErrOrderLoadAnotherUser
	}

	//вставка заказа
	_, err = s.pool.Exec(ctx, "insert into orders(orderid, userid, type, status, accrual) values($1, $2, 'accrual', 'NEW', 0);", orderID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Withdraw(ctx context.Context, orderID models.OrderID, userID models.UserID, sum float64) error {
	balance, err := s.GetBalance(ctx, userID)
	if err != nil {
		return err
	}

	if balance.Current < sum {
		return errs.ErrNotEnoughFunds
	}

	_, err = s.pool.Exec(ctx, "insert into orders(orderid, userid, type, status, accrual) values($1, $2, 'withdraw', 'NEW', $3);", orderID, userID, sum)
	if err != nil {
		return err
	}

	return nil
}
