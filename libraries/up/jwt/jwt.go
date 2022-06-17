package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var token = []byte("secret_token")

type Claims struct {
	jwt.RegisteredClaims

	AuthID uuid.UUID
}

func NewSignedString(authID uuid.UUID) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		AuthID: authID,
	}).SignedString(token)
}
