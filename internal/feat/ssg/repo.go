package ssg

import (
	"context"

	"github.com/adrianpk/hermes/internal/am"
)

type Repo interface {
	am.Repo

	CreateContent(ctx context.Context, content Content) error
	GetContent(ctx context.Context, id string) (Content, error)
	UpdateContent(ctx context.Context, content Content) error
	GetAllContent(ctx context.Context) ([]Content, error)
	// DeleteContent(ctx context.Context, contentID uuid.UUID) error
	CreateSection(ctx context.Context, section Section) error
	GetSections(ctx context.Context) ([]Section, error)
	CreateLayout(ctx context.Context, layout Layout) error
	GetAllLayouts(ctx context.Context) ([]Layout, error)
}
