package example

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

type Service struct {
	// Add dependencies here (repositories, clients, etc.)
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateGreeting(ctx context.Context, name string, uppercase bool) (string, error) {
	slog.InfoContext(ctx, "generating greeting",
		"name", name,
		"uppercase", uppercase,
	)

	message := fmt.Sprintf("Hello, %s! Welcome to Termplate Go.", name)

	if uppercase {
		message = strings.ToUpper(message)
	}

	return message, nil
}
