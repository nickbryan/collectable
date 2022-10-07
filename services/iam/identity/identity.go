package identity

import (
	"context"

	"github.com/google/uuid"

	"github.com/nickbryan/collectable/proto/iam/identity/service/v1"
)

type Repository interface {
	Create(ctx context.Context, id uuid.UUID, email, password string) error
}

type Service struct {
	identity.UnimplementedIdentityServiceServer

	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) CreateIdentity(ctx context.Context, request *identity.CreateIdentityRequest) (*identity.CreateIdentityResponse, error) {

	// TODO: add migrations and db env vars/connection - migrations will need to be a single dir with multiple migrations in instead of one file

	// TODO: need to do validation somewhere - check that passwords match etc.

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, nil // TODO: handle error here
	}

	if err := s.repo.Create(ctx, id, request.Email, request.Password); err != nil {
		return nil, nil // TODO: handle error here
	}

	return &identity.CreateIdentityResponse{Id: id.String()}, nil
}
