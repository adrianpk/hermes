package ssg

import (
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

type SectionForm struct {
	*am.BaseForm
	Name        string `form:"name" required:"true"`
	Description string `form:"description"`
	Path        string `form:"path"`
	LayoutID    string `form:"layout_id"`
	Image       string `form:"image"`
	Header      string `form:"header"`
}

func NewSectionForm(r *http.Request) SectionForm {
	return SectionForm{
		BaseForm: am.NewBaseForm(r),
	}
}

func SectionFormFromRequest(r *http.Request) (sf SectionForm, err error) {
	err = r.ParseForm()
	if err != nil {
		return sf, err
	}
	sf = NewSectionForm(r)
	err = am.ToForm(r, &sf)
	return sf, err
}

func (form *SectionForm) Validate() error {
	validate := am.ComposeValidators(
		am.MinLength("name", form.Name, 3),
		am.MaxLength("name", form.Name, 100),
	)
	v, err := validate(*form)
	if err != nil {
		return err
	}
	form.SetValidation(&v)
	return nil
}
