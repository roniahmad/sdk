package sdk

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/roniahmad/sdk"
	"github.com/roniahmad/sdk/model"
)

type MapClaims map[string]interface{}

func CreateToken(email string, id string, secret string, expiry int,
	keyPhrase string, issuer string) (string, time.Time, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry))

	claims := &model.JwtCustomClaims{
		Email: email,
		ID:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, exp, nil
}

func GetAuthToken(xAuthHeader string) (string, error) {
	const expectedScheme = sdk.Bearer
	parts := strings.SplitN(xAuthHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != expectedScheme {
		return "", sdk.ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func ExtractClaimsFromToken(token string, secret string) (MapClaims, error) {
	t, err := ParseTokenString(token, secret)
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)

	if !(ok && t.Valid) {
		return nil, sdk.ErrInvalidToken
	}

	mc := MapClaims{}
	for key, value := range claims {
		mc[key] = value
	}
	return mc, nil
}

func IsAuthorized(token string, secret string) (bool, error) {
	t, err := ParseTokenString(token, secret)
	if err != nil {
		return false, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !(ok && t.Valid) {
		return false, sdk.ErrInvalidToken
	}

	if expiresAt, ok := claims["exp"]; ok && int64(expiresAt.(float64)) < time.Now().Local().Unix() {
		return false, sdk.ErrExpiredToken
	}

	return true, nil
}

func ParseTokenString(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
}
