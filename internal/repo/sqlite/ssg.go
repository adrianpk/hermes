package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/adrianpk/hermes/internal/feat/ssg"
	"github.com/google/uuid"
)

var (
	featSSG    = "ssg"
	resLayout  = "layout"
	resContent = "content"
	resSection = "section"
)

// Content related

func (repo *HermesRepo) CreateContent(ctx context.Context, content ssg.Content) error {
	query, err := repo.Query().Get(featSSG, resContent, "Create")
	if err != nil {
		return err
	}

	// Note: m_type is not persisted, so it's not included here.
	// SectionID is now included.
	_, err = repo.db.ExecContext(ctx, query,
		content.GetID(),
		content.GetShortID(), // short_id
		content.UserID,
		content.SectionID,
		content.Heading,
		content.Body,
		content.Status,
		content.GetCreatedBy(),
		content.GetUpdatedBy(),
		content.GetCreatedAt(),
		content.GetUpdatedAt(),
	)
	return err
}

func (repo *HermesRepo) GetAllContent(ctx context.Context) ([]ssg.Content, error) {
	query, err := repo.Query().Get(featSSG, resContent, "GetAll")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []ssg.Content
	for rows.Next() {
		var (
			id        uuid.UUID
			userID    uuid.UUID
			sectionID uuid.UUID
			heading   string
			body      string
			status    string
			shortID   string
			createdBy uuid.UUID
			updatedBy uuid.UUID
			createdAt time.Time
			updatedAt time.Time
		)

		err := rows.Scan(
			&id, &userID, &sectionID, &heading, &body, &status, &shortID,
			&createdBy, &updatedBy, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		content := ssg.NewContent(heading, body)
		content.SetID(id)
		content.UserID = userID
		content.SectionID = sectionID
		content.Status = status
		content.SetShortID(shortID)
		content.SetCreatedBy(createdBy)
		content.SetUpdatedBy(updatedBy)
		content.SetCreatedAt(createdAt)
		content.SetUpdatedAt(updatedAt)
		content.SetType(resContent) // Set mType manually

		contents = append(contents, content)
	}

	return contents, nil
}

func (repo *HermesRepo) GetContent(ctx context.Context, id uuid.UUID) (ssg.Content, error) {
	query, err := repo.Query().Get(featSSG, resContent, "Get")
	if err != nil {
		return ssg.Content{}, err
	}

	row := repo.db.QueryRowxContext(ctx, query, id)

	var (
		contentID uuid.UUID
		userID    uuid.UUID
		sectionID uuid.UUID
		heading   string
		body      string
		status    string
		shortID   string
		createdBy uuid.UUID
		updatedBy uuid.UUID
		createdAt time.Time
		updatedAt time.Time
	)

	err = row.Scan(
		&contentID, &userID, &sectionID, &heading, &body, &status, &shortID,
		&createdBy, &updatedBy, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.Content{}, errors.New("content not found")
		}
		return ssg.Content{}, err
	}

	content := ssg.NewContent(heading, body)
	content.SetID(contentID)
	content.UserID = userID
	content.SectionID = sectionID
	content.Status = status
	content.SetShortID(shortID)
	content.SetCreatedBy(createdBy)
	content.SetUpdatedBy(updatedBy)
	content.SetCreatedAt(createdAt)
	content.SetUpdatedAt(updatedAt)
	content.SetType(resContent) // Set mType manually

	return content, nil
}

func (repo *HermesRepo) UpdateContent(ctx context.Context, content ssg.Content) error {
	query, err := repo.Query().Get(featSSG, resContent, "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query,
		content.Heading,
		content.Body,
		content.Status,
		content.GetUpdatedBy(),
		content.GetUpdatedAt(),
		content.GetID(),
	)
	return err
}

func (repo *HermesRepo) DeleteContent(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featSSG, resContent, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

// Section related

func (repo *HermesRepo) CreateSection(ctx context.Context, section ssg.Section) error {
	query, err := repo.Query().Get(featSSG, resSection, "Create")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query,
		section.GetID(),
		section.Name,
		section.Description,
		section.Path,
		section.LayoutID,
		section.Image,
		section.Header,
		section.GetShortID(),
		section.GetCreatedBy(),
		section.GetUpdatedBy(),
		section.GetCreatedAt(),
		section.GetUpdatedAt(),
	)
	return err
}

func (repo *HermesRepo) GetSections(ctx context.Context) ([]ssg.Section, error) {
	query, err := repo.Query().Get(featSSG, resSection, "GetAll")
	if err != nil {
		return nil, err
	}
	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []ssg.Section
	for rows.Next() {
		var s ssg.Section
		err := rows.Scan(
			&s.ID, &s.ShortID, &s.Name, &s.Description, &s.Path, &s.LayoutID, &s.Image, &s.Header,
			&s.CreatedBy, &s.UpdatedBy, &s.CreatedAt, &s.UpdatedAt, &s.LayoutName,
		)
		if err != nil {
			return nil, err
		}
		sections = append(sections, s)
	}
	return sections, nil
}

func (repo *HermesRepo) GetSection(ctx context.Context, id uuid.UUID) (ssg.Section, error) {
	query, err := repo.Query().Get(featSSG, resSection, "Get")
	if err != nil {
		return ssg.Section{}, err
	}

	row := repo.db.QueryRowxContext(ctx, query, id)

	var (
		sectionID   uuid.UUID
		name        string
		description string
		path        string
		layoutID    uuid.UUID
		image       string
		header      string
		shortID     string
		createdBy   uuid.UUID
		updatedBy   uuid.UUID
		createdAt   time.Time
		updatedAt   time.Time
	)

	err = row.Scan(
		&sectionID, &name, &description, &path, &layoutID, &image, &header, &shortID,
		&createdBy, &updatedBy, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ssg.Section{}, errors.New("section not found")
		}
		return ssg.Section{}, err
	}

	section := ssg.NewSection(name, description, path, layoutID)
	section.SetID(sectionID)
	section.Image = image
	section.Header = header
	section.SetShortID(shortID)
	section.SetCreatedBy(createdBy)
	section.SetUpdatedBy(updatedBy)
	section.SetCreatedAt(createdAt)
	section.SetUpdatedAt(updatedAt)

	return section, nil
}

func (repo *HermesRepo) UpdateSection(ctx context.Context, section ssg.Section) error {
	query, err := repo.Query().Get(featSSG, resSection, "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query,
		section.Name,
		section.Description,
		section.Path,
		section.LayoutID,
		section.Image,
		section.Header,
		section.GetShortID(),
		section.GetUpdatedBy(),
		section.GetUpdatedAt(),
		section.GetID(),
	)
	return err
}

func (repo *HermesRepo) DeleteSection(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featSSG, resSection, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

// Layout related

func (repo *HermesRepo) CreateLayout(ctx context.Context, layout ssg.Layout) error {
	query, err := repo.Query().Get(featSSG, "layout", "Create")
	if err != nil {
		return err
	}
	_, err = repo.db.ExecContext(ctx, query,
		layout.GetID(),
		layout.GetShortID(),
		layout.Name,
		layout.Description,
		layout.Code,
		layout.GetCreatedBy(),
		layout.GetUpdatedBy(),
		layout.GetCreatedAt(),
		layout.GetUpdatedAt(),
	)
	return err
}

func (repo *HermesRepo) GetAllLayouts(ctx context.Context) ([]ssg.Layout, error) {
	query, err := repo.Query().Get(featSSG, resLayout, "GetAll")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var layouts []ssg.Layout
	for rows.Next() {
		var (
			id          uuid.UUID
			shortID     string
			name        string
			description string
			code        string
			createdBy   uuid.UUID
			updatedBy   uuid.UUID
			createdAt   time.Time
			updatedAt   time.Time
		)

		err := rows.Scan(
			&id, &shortID, &name, &description, &code,
			&createdBy, &updatedBy, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		layout := ssg.Newlayout(name, description, code)
		layout.SetID(id)
		layout.SetShortID(shortID)
		layout.SetCreatedBy(createdBy)
		layout.SetUpdatedBy(updatedBy)
		layout.SetCreatedAt(createdAt)
		layout.SetUpdatedAt(updatedAt)

		layouts = append(layouts, layout)
	}

	return layouts, nil
}

func (repo *HermesRepo) GetLayout(ctx context.Context, id uuid.UUID) (ssg.Layout, error) {
	query, err := repo.Query().Get(featSSG, "layout", "Get")
	if err != nil {
		return ssg.Layout{}, err
	}

	var layout ssg.Layout
	err = repo.db.GetContext(ctx, &layout, query, id)
	if err != nil {
		return ssg.Layout{}, err
	}

	return layout, nil
}

func (repo *HermesRepo) UpdateLayout(ctx context.Context, layout ssg.Layout) error {
	query, err := repo.Query().Get(featSSG, resLayout, "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query,
		layout.Name,
		layout.Description,
		layout.Code,
		layout.GetUpdatedBy(),
		layout.GetUpdatedAt(),
		layout.GetID(),
	)
	return err
}

func (repo *HermesRepo) DeleteLayout(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featSSG, "layout", "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}
