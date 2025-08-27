package ssg

import (
	"context"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

type Service interface {
	CreateContent(ctx context.Context, content Content) error
	GetAllContent(ctx context.Context) ([]Content, error)
	GetContent(ctx context.Context, id uuid.UUID) (Content, error)
	UpdateContent(ctx context.Context, content Content) error
	DeleteContent(ctx context.Context, id uuid.UUID) error

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

func (svc *BaseService) GetContent(ctx context.Context, id uuid.UUID) (Content, error) {
	return svc.repo.GetContent(ctx, id)
}

func (svc *BaseService) UpdateContent(ctx context.Context, content Content) error {
	return svc.repo.UpdateContent(ctx, content)
}

func (svc *BaseService) DeleteContent(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteContent(ctx, id)
}

// Section related
func (svc *BaseService) CreateSection(ctx context.Context, section Section) error {
	return svc.repo.CreateSection(ctx, section)
}

func (svc *BaseService) GetSection(ctx context.Context, id uuid.UUID) (Section, error) {
	return svc.repo.GetSection(ctx, id)
}

func (svc *BaseService) GetSections(ctx context.Context) ([]Section, error) {
	return svc.repo.GetSections(ctx)
}

func (svc *BaseService) UpdateSection(ctx context.Context, section Section) error {
	return svc.repo.UpdateSection(ctx, section)
}

func (svc *BaseService) DeleteSection(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteSection(ctx, id)
}

// Layout related
func (svc *BaseService) CreateLayout(ctx context.Context, layout Layout) error {
	return svc.repo.CreateLayout(ctx, layout)
}

func (svc *BaseService) GetLayout(ctx context.Context, id uuid.UUID) (Layout, error) {
	return svc.repo.GetLayout(ctx, id)
}

func (svc *BaseService) GetAllLayouts(ctx context.Context) ([]Layout, error) {
	return svc.repo.GetAllLayouts(ctx)
}

func (svc *BaseService) UpdateLayout(ctx context.Context, layout Layout) error {
	return svc.repo.UpdateLayout(ctx, layout)
}

func (svc *BaseService) DeleteLayout(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteLayout(ctx, id)
}
