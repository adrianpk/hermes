package ssg

import (
	"github.com/adrianpk/hermes/internal/am"
)

// Content related

func ToContentDA(content Content) ContentDA {
	return ContentDA{
		ID:        content.ID(),
		UserID:    content.UserID,
		Heading:   content.Heading,
		Body:      content.Body,
		Status:    content.Status,
		ShortID:   content.ShortID(),
		CreatedBy: am.UUIDPtr(content.CreatedBy()),
		UpdatedBy: am.UUIDPtr(content.UpdatedBy()),
		CreatedAt: am.TimePtr(content.CreatedAt()),
		UpdatedAt: am.TimePtr(content.UpdatedAt()),
	}
}

func ToContent(da ContentDA) Content {
	return Content{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID),
			am.WithType(contentType),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		UserID:  da.UserID,
		Heading: da.Heading,
		Body:    da.Body,
		Status:  da.Status,
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
		Name:        section.Name,
		Description: section.Description,
		Path:        section.Path,
		LayoutID:    section.LayoutID.String(),
		ShortID:     section.ShortID(),
		CreatedBy:   am.UUIDPtr(section.CreatedBy()),
		UpdatedBy:   am.UUIDPtr(section.UpdatedBy()),
		CreatedAt:   am.TimePtr(section.CreatedAt()),
		UpdatedAt:   am.TimePtr(section.UpdatedAt()),
		Image:       section.Image,
		Header:      section.Header,
	}
}

func ToSection(da SectionDA) Section {
	return Section{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID),
			am.WithType(sectionType),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		Name:        da.Name,
		Description: da.Description,
		Path:        da.Path,
		LayoutID:    am.ParseUUID(da.LayoutID),
		Image:       da.Image,
		Header:      da.Header,
	}
}

func ToSections(das []SectionDA) []Section {
	sections := make([]Section, len(das))
	for i, da := range das {
		sections[i] = ToSection(da)
	}
	return sections
}

// Layout related

func ToLayoutDA(layout Layout) LayoutDA {
	return LayoutDA{
		ID:          layout.ID(),
		Name:        layout.Name,
		Description: layout.Description,
		Code:        layout.Code,
		ShortID:     layout.ShortID(),
		CreatedBy:   am.UUIDPtr(layout.CreatedBy()),
		UpdatedBy:   am.UUIDPtr(layout.UpdatedBy()),
		CreatedAt:   am.TimePtr(layout.CreatedAt()),
		UpdatedAt:   am.TimePtr(layout.UpdatedAt()),
	}
}

func ToLayout(da LayoutDA) Layout {
	return Layout{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID),
			am.WithType("layout"),
			am.WithCreatedBy(am.UUIDVal(da.CreatedBy)),
			am.WithUpdatedBy(am.UUIDVal(da.UpdatedBy)),
			am.WithCreatedAt(am.TimeVal(da.CreatedAt)),
			am.WithUpdatedAt(am.TimeVal(da.UpdatedAt)),
		),
		Name:        da.Name,
		Description: da.Description,
		Code:        da.Code,
	}
}

func ToLayouts(das []LayoutDA) []Layout {
	layouts := make([]Layout, len(das))
	for i, da := range das {
		layouts[i] = ToLayout(da)
	}
	return layouts
}

