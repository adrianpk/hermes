package ssg

import (
	"github.com/adrianpk/hermes/internal/am"
)

// Content related

func ToContentForm(content Content) ContentForm {
	return ContentForm{
		Heading: content.Heading,
		Body:    content.Body,
	}
}

func ToContentFromForm(form ContentForm) Content {
	return Content{
		BaseModel: am.NewModel(am.WithType(sectionType)),
		Heading:   form.Heading,
		Body:      form.Body,
	}
}

// Section related
func ToSectionForm(section Section) SectionForm {
	return SectionForm{
		Name:        section.Name,
		Description: section.Description,
		Path:        section.Path,
		LayoutID:    section.LayoutID.String(),
		Image:       section.Image,
		Header:      section.Header,
	}
}

func ToSectionFromForm(form SectionForm) Section {
	return Section{
		BaseModel:   am.NewModel(am.WithType(sectionType)),
		Name:        form.Name,
		Description: form.Description,
		Path:        form.Path,
		LayoutID:    am.ParseUUID(form.LayoutID),
		Image:       form.Image,
		Header:      form.Header,
	}
}

// Layout related
func ToLayoutForm(layout Layout) LayoutForm {
	return LayoutForm{
		Name:        layout.Name,
		Description: layout.Description,
		Code:        layout.Code,
	}
}

func ToLayoutFromForm(form LayoutForm) Layout {
	return Layout{
		BaseModel:   am.NewModel(am.WithType("layout")),
		Name:        form.Name,
		Description: form.Description,
		Code:        form.Code,
	}
}
