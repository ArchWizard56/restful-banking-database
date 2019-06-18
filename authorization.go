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
	jwt.StandardClaims
}

func GenToken (username string, tokenvalue string) (string, error) {
    DualDebug("Started JWT Generation")
    expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims : jwt.StandardClaims {
			ExpiresAt: expirationTime.Unix(),
            IssuedAt:  time.Now().Unix(),
	},
	}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    DualDebug("Signing Token")
    secret := append([]byte(Secret.Jwtsecret),[]byte(tokenvalue)...)
    tokenString, err := token.SignedString(secret) 
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
        claim := token.Claims.(*Claims)
        var DBTokenValue string
        var err error
        DBTokenValue, err = GetToken(Database, claim.Username)
        if err != nil {
            return []byte("bad"), err
        }
        secret := append([]byte(Secret.Jwtsecret),[]byte(DBTokenValue)...)
        DualDebug("generated new secret")
		return secret, nil
	})
    DualDebug("Parsed Token")
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
    if tkn.Valid {
        DualDebug("Success! Valid Token")
        return true, nil
    }
    return false, errors.New("Invalid Token Value")
}
func GetTokenClaims (Token string) (*Claims, error) {
    claims := &Claims{}
    DualDebug("Getting Token Claims")
    _, err := jwt.ParseWithClaims(Token, claims, func(token *jwt.Token) (interface{}, error) {
        claim := token.Claims.(*Claims)
        var DBTokenValue string
        var err error
        DBTokenValue, err = GetToken(Database, claim.Username)
        if err != nil {
            return []byte("bad"), err
        }
        secret := append([]byte(Secret.Jwtsecret),[]byte(DBTokenValue)...)
        DualDebug("generated new secret")
		return secret, nil
	})
    if err != nil {
        DualDebug(err.Error())
        return claims, err 
    }
    return claims, nil
}
