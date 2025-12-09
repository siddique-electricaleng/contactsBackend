package contacts

import (
	repo "contacts/internal/adapters/postgresql/sqlc"
)

type Service struct {
	q *repo.Queries
}
