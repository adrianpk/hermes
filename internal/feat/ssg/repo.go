package ssg

import (
	"context"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

type Repo interface {
	am.Repo

	CreateContent(ctx context.Context, content Content) error
	GetContent(ctx context.Context, id uuid.UUID) (Content, error)
	UpdateContent(ctx context.Context, content Content) error
	DeleteContent(ctx context.Context, id uuid.UUID) error
	GetAllContent(ctx context.Context) ([]Content, error)

	CreateSection(ctx context.Context, section Section) error
	GetSection(ctx context.Context, id uuid.UUID) (Section, error)
	GetSections(ctx context.Context) ([]Section, error)
	UpdateSection(ctx context.Context, section Section) error
	DeleteSection(ctx context.Context, id uuid.UUID) error

	CreateLayout(ctx context.Context, layout Layout) error
	GetLayout(ctx context.Context, id uuid.UUID) (Layout, error)
	GetAllLayouts(ctx context.Context) ([]Layout, error)
	UpdateLayout(ctx context.Context, layout Layout) error
	DeleteLayout(ctx context.Context, id uuid.UUID) error
}
