package main

import (
    "github.com/dgrijalva/jwt-go"
    //"encoding/json"
	"time"
    "errors"
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

func VerifyToken (Token string) (bool, error) {
    claims := &Claims{}
    DualDebug("Verifying Token")
    tkn, err := jwt.ParseWithClaims(Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret.Jwtsecret), nil
	})
    if err != nil {
        return false, err
    }
    if !tkn.Valid {
        DualDebug("Invalid Token")
		return false, nil
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
            DualDebug("Invalid Token")
			return false, nil
		}
        DualDebug("Bad request found")
		return false, errors.New("Bad Request")
	}
    var DBTokenValue int
    DBTokenValue, err = GetToken(Database, claims.Username)
    if err != nil {
        return false, nil
    }
    if claims.TokenValue == DBTokenValue {
        DualDebug("Success! Token Verified")
        return true, nil
    }
    DualDebug("Invalid Token Value Found")
    return false, errors.New("Invalid Token Value")
}
func GetTokenClaims (Token string) (*Claims, error) {
    claims := &Claims{}
    DualDebug("Getting Token Claims")
    _, err := jwt.ParseWithClaims(Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret.Jwtsecret), nil
	})
    if err != nil {
        DualWarning(err.Error())
        return claims, err 
    }
    return claims, nil
}
