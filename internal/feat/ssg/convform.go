package ssg

import (
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

// Content related

func ToContentForm(r *http.Request, content Content) ContentForm {
	return ContentForm{
		BaseForm:  am.NewBaseForm(r),
		ID:        content.ID().String(),
		Heading:   content.Heading,
		Body:      content.Body,
		SectionID: content.SectionID.String(),
	}
}

func ToContentFromForm(form ContentForm) Content {
	return Content{
		BaseModel: am.NewModel(am.WithID(am.ParseUUID(form.ID)), am.WithType(contentType)),
		Heading:   form.Heading,
		Body:      form.Body,
		SectionID: am.ParseUUID(form.SectionID),
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
