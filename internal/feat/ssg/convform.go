package ssg

// ToContentForm converts a Content business object to a ContentForm.
func ToContentForm(content Content) ContentForm {
	return ContentForm{
		Heading: content.Heading,
		Body:    content.Body,
	}
}

// ToContentFromForm creates a Content business object from a ContentForm.
// Note: You may want to set additional fields (e.g., UserID, BaseModel) elsewhere.
func ToContentFromForm(form ContentForm) Content {
	return Content{
		Heading: form.Heading,
		Body:    form.Body,
	}
}

