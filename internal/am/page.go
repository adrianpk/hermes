package am

import (
	"net/http"
	"net/url"
	"path"

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
}

// Feat struct represents a feature with base path, feature name, and action.
type Feat struct {
	Path       string // The base path for the action (i.e. `/feat/auth`)
	Action     string // Action name (command or query, i.e. `edit-user`)
	PathSuffix string // The path suffix for the feature (i.e.: `/auth`)
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
