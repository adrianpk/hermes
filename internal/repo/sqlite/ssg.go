package sqlite

import (
	"context"

	"github.com/adrianpk/hermes/internal/feat/ssg"
)

const (
	ssgAuth    = "ssg"
	resContent = "content"
)

func (repo *HermesRepo) CreateContent(ctx context.Context, content ssg.Content) error {
	query, err := repo.Query().Get(ssgAuth, resContent, "Create")
	if err != nil {
		return err
	}
	contentDA := ssg.ToContentDA(content)
	_, err = repo.db.NamedExecContext(ctx, query, contentDA)
	return err
}
