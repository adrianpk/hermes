package ssg

import (
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

type LayoutForm struct {
	*am.BaseForm
	Name        string `form:"name" required:"true"`
	Description string `form:"description"`
	Code        string `form:"code"`
}

func NewLayoutForm(r *http.Request) LayoutForm {
	return LayoutForm{
		BaseForm: am.NewBaseForm(r),
	}
}

func LayoutFormFromRequest(r *http.Request) (lf LayoutForm, err error) {
	err = r.ParseForm()
	if err != nil {
		return lf, err
	}
	lf = NewLayoutForm(r)
	err = am.ToForm(r, &lf)
	return lf, err
}

func (form *LayoutForm) Validate() error {
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
