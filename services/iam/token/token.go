package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"github.com/nickbryan/collectable/libraries/alup/jwt"
)

type Service struct {
	UnimplementedTokenServiceServer
}

func (s Service) CreateToken(ctx context.Context, request *CreateTokenRequest) (*CreateTokenResponse, error) {
	if request.Email != "test@example.org" || request.Password != "testpassword123!" {
		return nil, errors.New("invalid auth credentials")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("unable to generate uuid: %w", err)
	}

	token, err := jwt.NewToken(id)
	if err != nil {
		return nil, fmt.Errorf("unable to create jwt: %w", err)
	}

	return &CreateTokenResponse{Token: token}, nil
}

func NewTokenService() TokenServiceServer {
	return &Service{}
}
