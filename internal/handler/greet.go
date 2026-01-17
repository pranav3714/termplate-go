package handler

import (
	"context"
	"fmt"

	"github.com/blacksilver/ever-so-powerful/internal/model"
	"github.com/blacksilver/ever-so-powerful/internal/service/example"
)

type GreetInput struct {
	Name      string
	Uppercase bool
}

type GreetOutput struct {
	Message string
}

// GreetHandler handles greeting operations
type GreetHandler struct {
	service *example.Service
}

// NewGreetHandler creates a new greet handler
func NewGreetHandler() *GreetHandler {
	return &GreetHandler{
		service: example.NewService(),
	}
}

// Greet generates a greeting message
func (h *GreetHandler) Greet(ctx context.Context, in GreetInput) (*GreetOutput, error) {
	if in.Name == "" {
		return nil, model.NewValidationError("name", "name is required")
	}

	message, err := h.service.GenerateGreeting(ctx, in.Name, in.Uppercase)
	if err != nil {
		return nil, fmt.Errorf("generating greeting: %w", err)
	}

	return &GreetOutput{
		Message: message,
	}, nil
}
