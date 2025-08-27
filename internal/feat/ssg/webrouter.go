package ssg

import (
	"github.com/adrianpk/hermes/internal/am"
)

func NewWebRouter(handler *BFF, mw []am.Middleware, opts ...am.Option) *am.Router {
	core := am.NewWebRouter("web-router", opts...)
	core.SetMiddlewares(mw)

	// Content routes
	core.Get("/new-content", handler.NewContent)
	core.Post("/create-content", handler.CreateContent)
	core.Get("/edit-content", handler.EditContent)
	core.Post("/update-content", handler.UpdateContent)
	core.Get("/list-content", handler.ListContent)
	core.Get("/show-content", handler.ShowContent)
	core.Post("/delete-content", handler.DeleteContent)
	// Section routes
	core.Get("/new-section", handler.NewSection)
	core.Post("/create-section", handler.CreateSection)
	core.Get("/edit-section", handler.EditSection)
	core.Post("/update-section", handler.UpdateSection)
	core.Get("/list-sections", handler.ListSections)
	core.Get("/show-section", handler.ShowSection)
	core.Post("/delete-section", handler.DeleteSection)

	// Layout routes
	core.Get("/new-layout", handler.NewLayout)
	core.Post("/create-layout", handler.CreateLayout)
	core.Get("/edit-layout", handler.EditLayout)
	core.Post("/update-layout", handler.UpdateLayout)
	core.Get("/list-layouts", handler.ListLayouts)
	core.Post("/delete-layout", handler.DeleteLayout)

	return core
}
