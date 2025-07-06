package ssg

import (
	"database/sql"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

// Content related

func ToContentDA(content Content) ContentDA {
	return ContentDA{
		ID:        content.ID(),
		UserID:    sql.NullString{String: content.UserID.String(), Valid: content.UserID != uuid.Nil},
		Heading:   sql.NullString{String: content.Heading, Valid: content.Heading != ""},
		Body:      sql.NullString{String: content.Body, Valid: content.Body != ""},
		Status:    sql.NullString{String: content.Status, Valid: content.Status != ""},
		ShortID:   sql.NullString{String: content.ShortID(), Valid: content.ShortID() != ""},
		CreatedBy: sql.NullString{String: content.CreatedBy().String(), Valid: content.CreatedBy() != uuid.Nil},
		UpdatedBy: sql.NullString{String: content.UpdatedBy().String(), Valid: content.UpdatedBy() != uuid.Nil},
		CreatedAt: sql.NullTime{Time: content.CreatedAt(), Valid: !content.CreatedAt().IsZero()},
		UpdatedAt: sql.NullTime{Time: content.UpdatedAt(), Valid: !content.UpdatedAt().IsZero()},
	}
}

func ToContent(da ContentDA) Content {
	return Content{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID.String),
			am.WithType(contentType),
			am.WithCreatedBy(am.ParseUUIDNull(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUIDNull(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		UserID:  am.ParseUUIDNull(da.UserID),
		Heading: da.Heading.String,
		Body:    da.Body.String,
		Status:  da.Status.String,
	}
}

func ToContents(das []ContentDA) []Content {
	contents := make([]Content, len(das))
	for i, da := range das {
		contents[i] = ToContent(da)
	}
	return contents
}

// Section related

func ToSectionDA(section Section) SectionDA {
	return SectionDA{
		ID:          section.ID(),
		Name:        sql.NullString{String: section.Name, Valid: section.Name != ""},
		Description: sql.NullString{String: section.Description, Valid: section.Description != ""},
		Path:        sql.NullString{String: section.Path, Valid: section.Path != ""},
		LayoutID:    sql.NullString{String: section.LayoutID.String(), Valid: section.LayoutID != uuid.Nil},
		ShortID:     sql.NullString{String: section.ShortID(), Valid: section.ShortID() != ""},
		CreatedBy:   sql.NullString{String: section.CreatedBy().String(), Valid: section.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: section.UpdatedBy().String(), Valid: section.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: section.CreatedAt(), Valid: !section.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: section.UpdatedAt(), Valid: !section.UpdatedAt().IsZero()},
		Image:       sql.NullString{String: section.Image, Valid: section.Image != ""},
		Header:      sql.NullString{String: section.Header, Valid: section.Header != ""},
	}
}

func ToSection(da SectionDA) Section {
	return Section{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID.String),
			am.WithType(sectionType),
			am.WithCreatedBy(am.ParseUUIDNull(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUIDNull(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:   da.Name.String,
		Image:  da.Image.String,
		Header: da.Header.String,
	}
}

func ToSections(das []SectionDA) []Section {
	sections := make([]Section, len(das))
	for i, da := range das {
		sections[i] = ToSection(da)
	}
	return sections
}
