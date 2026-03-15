package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	key []byte
}
type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTToken(secret string) *JWTToken {
	return &JWTToken{
		key: []byte(secret),
	}
}

func (tm *JWTToken) GenerateToken(userID string) (string, error) {
	// buat claim opt
	claims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			// Subject:   userID, //via TokenClaim utk Payload
		},
	}
	// symetric/assymetric sign (hs256/ecdsa)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(tm.key)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (tm *JWTToken) VerifyToken(token string) (*TokenClaims, error) {
	t, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(t *jwt.Token) (any, error) {
		return tm.key, nil
	})
	if err != nil {
		return &TokenClaims{}, err
	}
	//prevent alg attack (assert hmac)
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	claims, ok := t.Claims.(*TokenClaims)
	if !ok || !t.Valid {
		return &TokenClaims{}, errors.New("Invalid Token")
	}
	return claims, nil
}
