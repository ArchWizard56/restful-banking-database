package main

import (
    "github.com/dgrijalva/jwt-go"
    //"encoding/json"
	"time"
)
type Credentials struct {
    Username string
    Password string
}

type Claims struct {
    Username string `json:"username"`
    TokenValue int `json:"tokenvalue"`
	jwt.StandardClaims
}

var ValidTokenValues map[string]int

func GenToken (username string, tokenvalue int) (string, error) {
    DualDebug("Started JWT Generation")
    expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		TokenValue: tokenvalue,
		StandardClaims : jwt.StandardClaims {
			ExpiresAt: expirationTime.Unix(),
            IssuedAt:  time.Now().Unix(),
	},
	}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    DualDebug("Signing Token")
    tokenString, err := token.SignedString([]byte(Secret.Jwtsecret)) 
    if err != nil {
        DualErr(err)
    }
    DualDebug("Success! Generated Token")
    return tokenString, nil
}
