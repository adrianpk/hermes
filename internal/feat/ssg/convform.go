package ssg

import "github.com/adrianpk/hermes/internal/am"

// Content related

func ToContentForm(content Content) ContentForm {
	return ContentForm{
		Heading: content.Heading,
		Body:    content.Body,
	}
}

func ToContentFromForm(form ContentForm) Content {
	return Content{
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
	return Section{
		Name:        form.Name,
		Description: form.Description,
		Path:        form.Path,
		LayoutID:    am.ParseUUID(form.LayoutID),
	}
}
