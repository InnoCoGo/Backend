package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/internal/repository"
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
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
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

func (s *AuthService) GenerateToken(uctx core.UserCtx) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:   uctx.UserId,
		Username: uctx.Username,
	})

	return token.SignedString([]byte(signInKey))
}

func (s *AuthService) ParseToken(accessToken string) (core.UserCtx, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}
		return []byte(signInKey), nil
	})
	if err != nil {
		return core.UserCtx{}, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return core.UserCtx{}, err
	}

	return core.UserCtx{
		UserId:   claims.UserId,
		Username: claims.Username,
	}, nil
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) VerifyTgAuthData(authData map[string]interface{}, keyword string) (bool, error) {
	checkHash, _ := authData["hash"].(string)
	authData["auth_date"] = int(authData["auth_date"].(float64))

	authIdWS, ok := authData["id"].(float64)
	if ok {
		authData["id"] = int(authIdWS)
	}

	delete(authData, "hash")

	var dataCheckArr []string
	for key, val := range authData {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%v", key, val)) // WARNING!
	}
	sort.Strings(dataCheckArr)
	dataCheckString := strings.Join(dataCheckArr, "\n")

	var h hash.Hash
	if keyword == "" {
		h = sha256.New()
	} else {
		h = hmac.New(sha256.New, []byte(keyword))
	}
	h.Write([]byte(os.Getenv("BOT_TOKEN")))
	secretKey := h.Sum(nil)

	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	hash := hex.EncodeToString(h.Sum(nil))

	log.Printf("\nsecretkey: %s\nkeyword: %s\nhash: %s\nchechHash: %s", hex.EncodeToString(secretKey), keyword, hash, checkHash)

	if hash != checkHash {
		return false, fmt.Errorf("the hashes don't match")
	}

	authDate, _ := authData["auth_date"].(int)
	if time.Now().Unix()-int64(authDate) > 86400 {
		return false, fmt.Errorf("expired auth date")
	}

	return true, nil
}
