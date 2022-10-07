-- name: CreateIdentity :exec
INSERT INTO identities (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);