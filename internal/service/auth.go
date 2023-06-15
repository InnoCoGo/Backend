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
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetTgUser(user core.User) (int, error) {
	id, err := s.repo.GetTgUser(user.Username, user.TgId)
	if err != nil {
		id, err = s.repo.CreateUser(user)
	}
	return id, err
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	userId, err := s.repo.GetUserId(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

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

// type TgUserWebSite struct {
// 	Id        int     `json:"id" binding:"required"`
// 	FirstName string  `json:"first_name" binding:"required"`
// 	LastName  string  `json:"last_name" binding:"required"`
// 	Username  string  `json:"username" binding:"required"`
// 	PhotoUrl  string  `json:"photo_url" binding:"required"`
// 	AuthDate  float64 `json:"auth_date" binding:"required"`
// 	Hash      string  `json:"hash" binding:"required"`
// }

// type TgUserWebApp struct {
// 	User     userTgWA `json:"user" binding:"required"`
// 	AuthDate float64  `json:"auth_date" binding:"required"`
// 	QueryId  string   `json:"query_id" binding:"required"`
// 	Hash     string   `json:"hash" binding:"required"`
// }

// type userTgWA struct {
// 	Id           int    `json:"id" binding:"required"`
// 	FirstName    string `json:"first_name" binding:"required"`
// 	LastName     string `json:"last_name" binding:"required"`
// 	Username     string `json:"username" binding:"required"`
// 	LanguageCode string `json:"language_code" binding:"required"`
// }

// func (s *AuthService) ParseTgToken(tokenToParse string) (int, error) {

// 	authData, secretKey, user, err := parseInput(c)
// 	if err != nil {
// 		return 0, err
// 	}

// 	checkHash, _ := authData["hash"].(string)

// 	delete(authData, "hash")

// 	var dataCheckArr []string
// 	for key, val := range authData {
// 		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%v", key, val)) // WARNING!
// 	}
// 	sort.Strings(dataCheckArr)

// 	dataCheckString := strings.Join(dataCheckArr, "\n")

// 	hmc := hmac.New(sha256.New, secretKey[:])
// 	hmc.Write([]byte(dataCheckString))
// 	hash := hmc.Sum(nil)

// 	if !hmac.Equal(hash, []byte(checkHash)) {
// 		newErrorResponse(c, http.StatusUnauthorized, "Data is NOT from Telegram")
// 		return
// 	}

// 	authDate, ok := authData["auth_date"].(float64)
// 	if !ok {
// 		newErrorResponse(c, http.StatusUnauthorized, "Invalid authorization data")
// 		return
// 	}
// 	if float64(time.Now().Unix()-int64(authDate)) > 86400 {
// 		newErrorResponse(c, http.StatusUnauthorized, "Data is outdated")
// 		return
// 	}

// 	// 															TODO! (Chech if user in database)
// 	userId, err := h.services.GetTgUser(user)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusUnauthorized, err.Error())
// 		return
// 	}

// 	c.Set(userCtx, userId)
// }

// func parseInput(c *gin.Context) (map[string]interface{}, []byte, core.User, error) {
// 	var tguserWA TgUserWebApp
// 	var tguserWS TgUserWebSite

// 	authData := make(map[string]interface{})

// 	secretKey := make([]byte, 32)

// 	var user core.User

// 	if err := c.BindJSON(&tguserWA); err != nil {
// 		if err = c.BindJSON(&tguserWS); err != nil {
// 			return nil, secretKey, user, err
// 		}

// 		jsonData, _ := json.Marshal(tguserWS)
// 		json.Unmarshal(jsonData, &authData)

// 		h := hmac.New(sha256.New, []byte(webappKeyword))
// 		h.Write([]byte(BotToken))
// 		secretKey := h.Sum(nil)

// 		user = core.User{
// 			FirstName: tguserWS.FirstName,
// 			LastName:  tguserWS.LastName,
// 			Username:  tguserWS.Username,
// 			TgId:      tguserWS.Id,
// 		}

// 		return authData, secretKey, user, nil
// 	}

// 	jsonData, _ := json.Marshal(tguserWA)
// 	json.Unmarshal(jsonData, &authData)

// 	// secretKey = sha256.Sum256([]byte(BotToken))				WARNING!
// 	h := hmac.New(sha256.New, nil)
// 	h.Write([]byte(BotToken))
// 	secretKey = h.Sum(nil)

// 	user = core.User{
// 		FirstName: tguserWA.User.FirstName,
// 		LastName:  tguserWA.User.LastName,
// 		Username:  tguserWA.User.Username,
// 		TgId:      tguserWA.User.Id,
// 	}

// 	return authData, secretKey, user, nil
// }

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
