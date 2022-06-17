package token

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nickbryan/collectable/libraries/up/jwt"
	token "github.com/nickbryan/collectable/proto/iam/token/service/v1"
)

type Service struct {
	token.UnimplementedTokenServiceServer
}

func (s Service) CreateToken(ctx context.Context, request *token.CreateTokenRequest) (*token.CreateTokenResponse, error) {
	if request.Email != "test@example.org" || request.Password != "testpassword123!" {
		return nil, status.Error(codes.NotFound, "unable to create token from auth details")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("unable to generate uuid: %w", err)
	}

	tkn, err := jwt.NewSignedString(id)
	if err != nil {
		return nil, fmt.Errorf("unable to create jwt: %w", err)
	}

	return &token.CreateTokenResponse{Token: tkn}, nil
}

func NewTokenService() token.TokenServiceServer {
	return &Service{}
}
