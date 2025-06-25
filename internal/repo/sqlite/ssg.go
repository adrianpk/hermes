package sqlite

import (
	"context"

	"github.com/adrianpk/hermes/internal/feat/ssg"
)

func (repo *HermesRepo) CreateContent(ctx context.Context, content ssg.Content) error {
	query := `
		INSERT INTO content (
			id, user_id, heading, body, status, short_id, created_by, updated_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	contentDA := ssg.ToContentDA(content) 
	exec := repo.getExec(ctx)
	_, err := exec.ExecContext(ctx, query,
		contentDA.ID,
		contentDA.UserID,
		contentDA.Heading,
		contentDA.Body,
		contentDA.Status,
		contentDA.ShortID,
		contentDA.CreatedBy,
		contentDA.UpdatedBy,
		contentDA.CreatedAt,
		contentDA.UpdatedAt,
	)
	return err
}
