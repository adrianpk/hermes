package ssg

import (
	"context"

	"github.com/adrianpk/hermes/internal/am"
)

type Service interface {
	CreateContent(ctx context.Context, content Content) error
	GetAllContent(ctx context.Context) ([]Content, error)
	GetContent(ctx context.Context, id string) (Content, error)
	UpdateContent(ctx context.Context, content Content) error
	// GetAllContent(ctx context.Context) ([]Content, error)
	// DeleteContent(ctx context.Context, id uuid.UUID) error
	CreateSection(ctx context.Context, section Section) error
	GetSections(ctx context.Context) ([]Section, error)
	CreateLayout(ctx context.Context, layout Layout) error
	GetAllLayouts(ctx context.Context) ([]Layout, error)
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

func (svc *BaseService) GetAllContent(ctx context.Context) ([]Content, error) {
	return svc.repo.GetAllContent(ctx)
}

func (svc *BaseService) GetContent(ctx context.Context, id string) (Content, error) {
	return svc.repo.GetContent(ctx, id)
}

func (svc *BaseService) UpdateContent(ctx context.Context, content Content) error {
	return svc.repo.UpdateContent(ctx, content)
}

// Section related
func (svc *BaseService) CreateSection(ctx context.Context, section Section) error {
	return svc.repo.CreateSection(ctx, section)
}

func (svc *BaseService) GetSections(ctx context.Context) ([]Section, error) {
	return svc.repo.GetSections(ctx)
}

// Layout related
func (svc *BaseService) CreateLayout(ctx context.Context, layout Layout) error {
	return svc.repo.CreateLayout(ctx, layout)
}

func (svc *BaseService) GetAllLayouts(ctx context.Context) ([]Layout, error) {
	return svc.repo.GetAllLayouts(ctx)
}
