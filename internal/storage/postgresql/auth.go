package postgresql

import (
	"context"
	"errors"
	"github.com/dubrovsky1/gophermart/internal/errs"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) Register(ctx context.Context, user models.User) (string, error) {
	var userID string
	row := s.pool.QueryRow(ctx, "insert into users(userid, login, password) values ($1, $2, $3) returning userid;", user.UserID, user.Login, user.Password)

	err := row.Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return "", errs.ErrAlreadyExists
		}
		return "", err
	}
	return userID, nil
}

func (s *Storage) Login(ctx context.Context, user models.User) (string, error) {
	var userID string
	row := s.pool.QueryRow(ctx, "select u.userid from users u where u.login = $1 and u.password = $2;", user.Login, user.Password)

	err := row.Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", errs.ErrNotExists
		}
		return "", err
	}

	return userID, nil
}
