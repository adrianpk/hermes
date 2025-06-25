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

func (svc *BaseService) CreateContent(ctx context.Context, role Content) error {
	return svc.repo.CreateContent(ctx, role)
}

