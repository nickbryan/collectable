// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: identity.sql

package postgresql

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createIdentity = `-- name: CreateIdentity :exec
INSERT INTO identities (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)
`

type CreateIdentityParams struct {
	ID        uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateIdentity(ctx context.Context, arg CreateIdentityParams) error {
	_, err := q.db.Exec(ctx, createIdentity,
		arg.ID,
		arg.Email,
		arg.Password,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}
