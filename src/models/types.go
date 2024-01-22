package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Table struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type Man struct {
	Email         string `json:"email"`
	PreferredName string `json:"preferredName"`
	Password      string `json:"password"`
	TableID       int    `json:"tableID"`
}

type Event struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	StartTime string `json:"startTime"`
	TableID   int    `json:"tableID"`
}
