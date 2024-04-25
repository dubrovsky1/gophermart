package postgresql

import (
	"context"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/jackc/pgx/v5"
	"strconv"
	"time"
)

func (s *Storage) GetOrderList(ctx context.Context, userID models.UserID) ([]models.Order, error) {
	rows, err := s.pool.Query(ctx, "select (o.orderid, o.status, o.accrual, o.uploaded_at) from orders o where o.userid = $1 order by o.uploaded_at;", userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.ErrNotExists
		}
		return nil, err
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowTo[models.GetOrderListResult])
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errs.ErrNotExists
	}

	//logger.Sugar.Infow("GetOrderList Log.", "result", result)

	var orders []models.Order

	for _, item := range result {
		orderID, _ := strconv.Atoi(item.OrderID)

		order := models.Order{
			OrderID: models.OrderID(orderID),
			Status:  item.Status,
			Accrual: item.Accrual,
			Upload:  item.Upload.Format(time.RFC3339),
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *Storage) GetBalance(ctx context.Context, userID models.UserID) (models.Balance, error) {
	var balance models.Balance

	row := s.pool.QueryRow(ctx, `
										with balance as 
										(
											select o.type, 
											       o.accrual 
											from orders o
											where o.userid = $1
											and o.status = 'PROCESSED'
											and o.accrual is not null

											union all

											select 'accrual' as type, --на случай, когда нет подходящих заказов, чтобы вернуть нули вместо пустого ответа
											        0 as accrual 
										)
										select distinct sum(case when b.type = 'accrual' then b.accrual else b.accrual * -1 end) over() as accrual,
											            sum(case when b.type = 'withdraw' then b.accrual else 0 end) over() as withdraw
										from balance b;`, userID)

	err := row.Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return models.Balance{}, err
	}
	return balance, nil
}

func (s *Storage) Withdrawals(ctx context.Context, userID models.UserID) ([]models.Withdraw, error) {
	rows, err := s.pool.Query(ctx, "select (o.orderid, o.accrual, o.uploaded_at) from orders o where o.userid = $1 and o.status = 'PROCESSED' and o.type = 'withdraw' order by o.uploaded_at;", userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.ErrNotExists
		}
		return nil, err
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowTo[models.WithdrawListResult])
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errs.ErrNotExists
	}

	var withdrawals []models.Withdraw

	for _, item := range result {
		orderID, _ := strconv.Atoi(item.OrderID)

		withdraw := models.Withdraw{
			OrderID: models.OrderID(orderID),
			Sum:     item.Sum,
			Upload:  item.Upload.Format(time.RFC3339),
		}

		withdrawals = append(withdrawals, withdraw)
	}

	return withdrawals, nil
}

func (s *Storage) GetOrderAccrualInfo(ctx context.Context, orderID models.OrderID) (models.OrderAccrual, error) {
	var orderAccrual models.OrderAccrual
	row := s.pool.QueryRow(ctx, "select o.orderid, o.status, o.accrual from orders o where o.orderid = $1;", orderID)

	err := row.Scan(&orderAccrual.OrderID, &orderAccrual.Status, &orderAccrual.Accrual)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.OrderAccrual{}, errs.ErrNotExists
		}
		return models.OrderAccrual{}, err
	}

	return orderAccrual, nil
}
