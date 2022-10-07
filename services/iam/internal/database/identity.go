package database

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nickbryan/collectable/services/iam/internal/database/postgresql"
)

type IdentityRepository struct {
	queries *postgresql.Queries
}

func NewIdentityRepository(db *postgresql.Queries) *IdentityRepository {
	return &IdentityRepository{queries: db}
}

func (r *IdentityRepository) Create(ctx context.Context, id uuid.UUID, email, password string) error {
	now := time.Now()

	return r.queries.CreateIdentity(ctx, postgresql.CreateIdentityParams{
		ID:        id,
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	})
}
