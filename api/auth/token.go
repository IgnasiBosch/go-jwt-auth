package auth

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strings"
	"time"
)

const AccessTokenTTL = time.Minute * 10
const RefreshTokenTTL = time.Hour * 24

func CreateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"token_type": "access",
		"sub":        userID,
		"exp":        time.Now().UTC().Add(AccessTokenTTL).Unix(), // expiration at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func CreateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"token_type": "refresh",
		"sub":        userID,
		"exp":        time.Now().UTC().Add(RefreshTokenTTL).Unix(), // expiration at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func AccessTokenValid(r *http.Request) bool {
	tokenString := ExtractToken(r)
	token, err := ParseToken(tokenString)
	if err != nil {
		return false
	}

	return token.Valid
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexppected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
}

func ExtractToken(r *http.Request) string {
	headerAuth := r.Header.Get("Authorization")
	authValue := strings.Split(headerAuth, " ")
	if len(authValue) != 2 || authValue[0] != "Bearer" {
		return ""
	}

	return authValue[1]
}

func ExtractUserID(r *http.Request) (uuid.UUID, error) {
	tokenString := ExtractToken(r)
	token, err := ParseToken(tokenString)
	if err != nil {
		return uuid.UUID{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok {
		return uuid.Parse(fmt.Sprint(claims["sub"]))
	}

	return uuid.UUID{}, errors.New("invalid claims")
}
