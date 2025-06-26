package am

import (
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/csrf"
)

// Page struct represents a web page with data, flash messages, form, menu, and feature information.
// Page is a generic struct that can carry a single entity (Entity) and a collection of entities (Entities).
type Page[T any] struct {
	Form     Form
	Entity   T           // Single entity instance of any type
	Entities []T         // Collection of entities of the same type
	Data     interface{} // For any additional or wrapper data
	Flash    Flash
	Menu     *Menu
	Feat     Feat
}

// Feat struct represents a feature with base path, feature name, and action.
type Feat struct {
	Path       string // The base path for the action (i.e. `/feat/auth`)
	Action     string // Action name (command or query, i.e. `edit-user`)
	PathSuffix string // The path suffix for the feature (i.e.: `/auth`)
}

// NewPage creates a new Page with the given data.
func NewPage[T any](r *http.Request, data interface{}) *Page[T] {
	return &Page[T]{
		Form:  NewBaseForm(r),
		Data:  data,
		Flash: NewFlash(),
		Menu:  NewMenu("/"),
	}
}

// SetForm sets the entire Form for the page.
func (p *Page[T]) SetForm(form Form) {
	p.Form = form
}

// SetData sets the Data property for the page.
func (p *Page[T]) SetData(data interface{}) {
	p.Data = data
}

// SetFlash sets the flash message for the page.
func (p *Page[T]) SetFlash(flash Flash) {
	p.Flash = flash
}

// SetFeat sets the Feat struct for the page.
func (p *Page[T]) SetFeat(feat Feat) {
	p.Feat = feat
}

// SetMenuItems sets the menu items for the page.
func (p *Page[T]) SetMenuItems(items []MenuItem) {
	p.Menu.Items = items
}

// GenCSRFToken generates a csrf token and sets it in the form.
func (p *Page[T]) GenCSRFToken(r *http.Request) {
	p.Form.SetCSRF(csrf.Token(r))
}

// Path generates the href path for a menu item based on the feature and menu item data.
func (p *Page[T]) Path(feat Feat, item MenuItem) string {
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
func (p *Page[T]) NewMenu(path string) *Menu {
	menu := NewMenu(path)
	p.Menu = menu
	return menu
}
