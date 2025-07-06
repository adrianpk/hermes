package am

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

// Page struct represents a web page with data, fm messages, form, menu, and feature information.
type Page struct {
	Form     Form
	Entity   any   // Single entity instance of any type
	Entities []any // Collection of entities of any type
	Data     any   // For any additional or wrapper data
	Flash    Flash
	Menu     *Menu
	Feat     Feat
	Selects  map[string][]SelectOption // For select/dropdown options and similar lists
}

// Feat struct represents a feature with base path, feature name, and action.
type Feat struct {
	Path       string // The base path for the action (i.e. `/feat/auth`)
	Action     string // Action name (command or query, i.e. `edit-user`)
	PathSuffix string // The path suffix for the feature (i.e.: `/auth`)
}

// SelectOption struct represents an option for select/dropdown lists.
type SelectOption struct {
	Value string
	Label string
}

// NewPage creates a new Page with the given data.
func NewPage(r *http.Request, data interface{}) *Page {
	flash, _ := r.Context().Value(FlashKey).(Flash)
	return &Page{
		Form:  NewBaseForm(r),
		Data:  data,
		Flash: flash,
		Menu:  NewMenu("/"),
	}
}

// SetForm sets the entire Form for the page.
func (p *Page) SetForm(form Form) {
	p.Form = form
}

// SetData sets the Data property for the page.
func (p *Page) SetData(data interface{}) {
	p.Data = data
}

// SetFlash sets the fm message for the page.
func (p *Page) SetFlash(flash Flash) {
	p.Flash = flash
}

// SetFeat sets the Feat struct for the page.
func (p *Page) SetFeat(feat Feat) {
	p.Feat = feat
}

// SetMenuItems sets the menu items for the page.
func (p *Page) SetMenuItems(items []MenuItem) {
	p.Menu.Items = items
}

// GenCSRFToken generates a csrf token and sets it in the form.
func (p *Page) GenCSRFToken(r *http.Request) {
	p.Form.SetCSRF(csrf.Token(r))
}

// Path generates the href path for a menu item based on the feature and menu item data.
func (p *Page) Path(feat Feat, item MenuItem) string {
	basePath := path.Join(feat.Path, feat.Action)

	if len(item.QueryParams) == 0 {
		return basePath
	}

	query := url.Values{}
	for key, value := range item.QueryParams {
		query.Add(key, value)
	}

	return basePath + "?" + query.Encode()
}

// NewMenu returns a new menu associated with this page, configured with the
// given path
func (p *Page) NewMenu(path string) *Menu {
	menu := NewMenu(path)
	p.Menu = menu
	return menu
}

// SetSelects sets multiple select option lists for selects/dropdowns.
func (p *Page) SetSelects(selects map[string][]SelectOption) {
	p.Selects = selects
}

// AddSelect adds or updates a single select option list for a select/dropdown.
func (p *Page) AddSelect(key string, values []SelectOption) {
	if p.Selects == nil {
		p.Selects = make(map[string][]SelectOption)
	}
	p.Selects[key] = values
}

// GetSelects retrieves a select option list by key from the Selects map.
func (p *Page) GetSelects(key string) []SelectOption {
	if p.Selects == nil {
		return nil
	}
	return p.Selects[key]
}

// ToAnySlice converts a slice of any type to a slice of any (for use in Page.Selects).
func ToAnySlice[T any](in []T) []any {
	out := make([]any, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

// ToSelectOptions converts a slice of models (with ID() and Name fields) to []SelectOption.
func ToSelectOptions[T interface {
	ID() any
	GetName() string
}](in []T) []SelectOption {
	out := make([]SelectOption, len(in))
	for i, v := range in {
		out[i] = SelectOption{
			Value: fmt.Sprintf("%v", v.ID()),
			Label: v.GetName(),
		}
	}
	return out
}

// ToSelectOptionsWith extracts select options from a slice using custom value and label functions.
func ToSelectOptionsWith[T any](in []T, valueFn func(T) string, labelFn func(T) string) []SelectOption {
	out := make([]SelectOption, len(in))
	for i, v := range in {
		out[i] = SelectOption{
			Value: valueFn(v),
			Label: labelFn(v),
		}
	}
	return out
}

// ToSelectOptionsWithID extracts select options from a slice using a custom label function. Value is always ID() as string.
func ToSelectOptionsWithID[T interface{ ID() uuid.UUID }](in []T, labelFn func(T) string) []SelectOption {
	out := make([]SelectOption, len(in))
	for i, v := range in {
		out[i] = SelectOption{
			Value: v.ID().String(),
			Label: labelFn(v),
		}
	}
	return out
}
