package sqlite

import (
	"context"

	"github.com/adrianpk/hermes/internal/feat/ssg"
)

const (
	ssgAuth    = "ssg"
	resContent = "content"
	resSection = "section"
)

// Content related

func (repo *HermesRepo) CreateContent(ctx context.Context, content ssg.Content) error {
	query, err := repo.Query().Get(ssgAuth, resContent, "Create")
	if err != nil {
		return err
	}

	contentDA := ssg.ToContentDA(content)
	_, err = repo.db.NamedExecContext(ctx, query, contentDA)
	return err
}

// Section related

func (repo *HermesRepo) CreateSection(ctx context.Context, section ssg.Section) error {
	query, err := repo.Query().Get(ssgAuth, resSection, "Create")
	if err != nil {
		return err
	}

	sectionDA := ssg.ToSectionDA(section)
	_, err = repo.db.NamedExecContext(ctx, query, sectionDA)
	return err
}

func (repo *HermesRepo) GetSections(ctx context.Context) ([]ssg.Section, error) {
	query, err := repo.Query().Get(ssgAuth, resSection, "GetAll")
	if err != nil {
		return nil, err
	}
	var das []ssg.SectionDA
	err = repo.db.SelectContext(ctx, &das, query)
	if err != nil {
		return nil, err
	}
	return ssg.ToSections(das), nil
}
