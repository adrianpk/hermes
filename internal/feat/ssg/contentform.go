package ssg

import (
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

type ContentForm struct {
	*am.BaseForm
	ID        string `form:"id"`
	Heading   string `form:"heading" required:"true"`
	Body      string `form:"body"`
	SectionID string `form:"section_id"`
}

func NewContentForm(r *http.Request) ContentForm {
	return ContentForm{
		BaseForm: am.NewBaseForm(r),
		ID:       r.URL.Query().Get("id"),
	}
}

func ContentFormFromRequest(r *http.Request) (cf ContentForm, err error) {
	err = r.ParseForm()
	if err != nil {
		return cf, err
	}

	return ContentForm{
		BaseForm: am.NewBaseForm(r),
		ID:       r.Form.Get("id"),
		Heading:  r.Form.Get("heading"),
		Body:     r.Form.Get("body"),
		SectionID: r.Form.Get("section_id"),
	}, nil
}

// Validate validates a ContentForm using am validators.
func (form *ContentForm) Validate() (err error) {
	validate := am.ComposeValidators(
		am.MinLength("heading", form.Heading, 3),
		am.MaxLength("heading", form.Heading, 100),
		am.MinLength("body", form.Body, 1),
	)

	v, err := validate(*form)
	if err != nil {
		return err
	}

	form.SetValidation(&v)

	return err
}
