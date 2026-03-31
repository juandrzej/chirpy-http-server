package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	mySigningKey := []byte(tokenSecret)
	now := time.Now().UTC()

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	} else if !token.Valid {
		return uuid.Nil, errors.New("Invalid token.")
	} else {
		return uuid.Parse(claims.Subject)
	}
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if len(token) == 0 {
		return "", errors.New("No authorization in headers.")
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return "", errors.New("Token does not start with 'Bearer '.")
	}
	token = strings.TrimPrefix(token, "Bearer ")

	return token, nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}

func GetAPIKey(headers http.Header) (string, error) {
	key := headers.Get("Authorization")
	if len(key) == 0 {
		return "", errors.New("No authorization in headers.")
	}

	if !strings.HasPrefix(key, "ApiKey ") {
		return "", errors.New("Key does not start with ApiKey.")
	}
	key = strings.TrimPrefix(key, "ApiKey ")

	return key, nil
}
