package view

import (
	"RestSample/dgrijalva/jwt-go"
)

type Users struct {
	UserId   int    `json:"userid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	IsActive bool   `json:"status"`
}

//for jwt

type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}
