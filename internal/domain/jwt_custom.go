package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
	SSID string `json:"ssid"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID   int64  `json:"id"`
	SSID string `json:"ssid"`
	jwt.RegisteredClaims
}
