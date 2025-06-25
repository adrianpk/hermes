package ssg

import (
	"database/sql"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/google/uuid"
)

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
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		UserID:  am.ParseUUID(da.UserID),
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
