package domain

import (
	dbAdapters "ai-agent-business/internal/app/adapters/database"
	"ai-agent-business/internal/app/adapters/database/postgres"
	"ai-agent-business/internal/app/domain/book"
	"ai-agent-business/internal/app/domain/health"
	"ai-agent-business/internal/infra/database"
)

type Services struct {
	Book   *book.Service
	Health *health.Service
}

func NewServices(dbs *database.Databases) *Services {
	return &Services{
		Book:   book.NewService(postgres.NewBookRepository(dbs)),
		Health: health.NewService(dbAdapters.NewHealthRepository(dbs)),
	}
}
