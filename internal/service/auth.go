package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"github.com/dubrovsky1/gophermart/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	salt      = "as3rasd421dsf"
	tokenExp  = time.Hour * 3
	SecretKey = "supersecretkey"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID models.UserID `json:"userid"`
}

func (s *Service) Register(ctx context.Context, user models.User) (models.UserID, error) {
	user.Password = generateHashedPassword(user.Password)

	userID, err := s.storage.Register(ctx, user)
	if err != nil {
		return "", err
	}

	//сразу аутентифицируемся после регистрации, передавая в ответе токен с созданным userid
	token, err := buildJWTString(userID)
	if err != nil {
		return "", err
	}

	return models.UserID(token), nil
}

func (s *Service) Login(ctx context.Context, user models.User) (models.UserID, error) {
	user.Password = generateHashedPassword(user.Password)

	userID, err := s.storage.Login(ctx, user)
	if err != nil {
		return "", err
	}

	token, err := buildJWTString(userID)
	if err != nil {
		return "", err
	}

	return models.UserID(token), nil
}

func generateHashedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPassword := hash.Sum([]byte(salt))
	return "{SHA256}" + base64.StdEncoding.EncodeToString(hashedPassword)
}

func buildJWTString(userID models.UserID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
