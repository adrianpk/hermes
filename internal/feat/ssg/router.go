package ssg

import (
	"github.com/adrianpk/hermes/internal/am"
)

func NewWebRouter(handler *WebHandler, mw []am.Middleware, opts ...am.Option) *am.Router {
	core := am.NewWebRouter("web-router", opts...)
	core.SetMiddlewares(mw)

	// Content routes
	core.Get("/new-content", handler.NewContent)
	core.Post("/create-content", handler.CreateContent)
	// core.Get("/show-content", handler.ShowContent)
	// core.Get("/edit-content", handler.EditContent)
	// core.Post("/update-content", handler.UpdateContent)
	core.Get("/list-content", handler.ListContent)
	// core.Post("/delete-content", handler.DeleteContent)
	// Section routes
	core.Get("/new-section", handler.NewSection)
	core.Post("/create-section", handler.CreateSection)

	// Layout routes
	core.Get("/new-layout", handler.NewLayout)
	core.Post("/create-layout", handler.CreateLayout)

	return core
}
