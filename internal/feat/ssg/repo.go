package ssg

import (
	"context"

	"github.com/adrianpk/hermes/internal/am"
)

type Repo interface {
	am.Repo

	CreateContent(ctx context.Context, role Content) error
	// GetContent(ctx context.Context, roleID uuid.UUID, preload ...bool) (Content, error)
	// UpdateContent(ctx context.Context, role Content) error
	// GetAllContent(ctx context.Context) ([]Content, error)
	// DeleteContent(ctx context.Context, roleID uuid.UUID) error
}

