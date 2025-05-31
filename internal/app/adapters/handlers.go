package adapters

import (
	"ai-agent-business/internal/app/adapters/handler"
	"ai-agent-business/internal/app/domain"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	book   *handler.BookHandler
	health *handler.HealthHandler
}

func NewHandlers(services *domain.Services) *Handlers {
	return &Handlers{
		book:   handler.NewBookHandler(services),
		health: handler.NewHealthHandler(services),
	}
}

func (h *Handlers) Configure(server *echo.Echo) {
	h.book.Configure(server)
	h.health.Configure(server)
}
