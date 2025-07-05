package ssg

import (
	"context"

	"github.com/adrianpk/hermes/internal/am"
)

type Service interface {
	CreateContent(ctx context.Context, content Content) error
	// GetContent(ctx context.Context, id uuid.UUID) (Content, error)
	// UpdateContent(ctx context.Context, content Content) error
	// GetAllContent(ctx context.Context) ([]Content, error)
	// DeleteContent(ctx context.Context, id uuid.UUID) error
	CreateSection(ctx context.Context, section Section) error
}

var (
	key = am.Key
)

type BaseService struct {
	*am.Service
	repo Repo
}

func NewService(repo Repo) *BaseService {
	return &BaseService{
		Service: am.NewService("ssg-service"),
		repo:    repo,
	}
}

// Content related

func (svc *BaseService) CreateContent(ctx context.Context, content Content) error {
	return svc.repo.CreateContent(ctx, content)
}

// Section related
func (svc *BaseService) CreateSection(ctx context.Context, section Section) error {
	return svc.repo.CreateSection(ctx, section)
}
