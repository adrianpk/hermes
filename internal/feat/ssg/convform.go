package ssg

import (
	"fmt"

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
		BaseModel:   am.NewModel(am.WithType(sectionType)),
		Heading: form.Heading,
		Body:    form.Body,
	}
}

// Section related
func ToSectionForm(section Section) SectionForm {
	return SectionForm{
		Name:        section.Name,
		Description: section.Description,
		Path:        section.Path,
		LayoutID:    section.LayoutID.String(),
	}
}

func ToSectionFromForm(form SectionForm) Section {
	fmt.Printf("Converting SectionForm to Section: %+v\n", form)
	return Section{
		BaseModel:   am.NewModel(am.WithType(sectionType)),
		Name:        form.Name,
		Description: form.Description,
		Path:        form.Path,
		LayoutID:    am.ParseUUID(form.LayoutID),
	}
}
