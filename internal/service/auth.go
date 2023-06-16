package service

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
	"github.com/itoqsky/InnoCoTravel_backend/internal/repository"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

const (
	salt      = "nfjdpsabnuirnefnjdsfds"
	signInKey = "fn9ht3s4h8f2finqjwnadfeu93nqfew"
	tokenTTL  = 24 * time.Hour
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"userId"`
}

func (s *AuthService) CreateUser(user core.User) (int, error) {
	if user.PasswordOrHash != "" {
		user.PasswordOrHash = generatePasswordHash(user.PasswordOrHash)
	}
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUserId(user core.User) (int, error) {
	if user.PasswordOrHash != "" {
		user.PasswordOrHash = generatePasswordHash(user.PasswordOrHash)
	}
	return s.repo.GetUserId(user)
}

func (s *AuthService) GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: userId,
	})

	return token.SignedString([]byte(signInKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}
		return []byte(signInKey), nil
	})
	if err != nil {
		return -1, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return -1, err
	}
	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
