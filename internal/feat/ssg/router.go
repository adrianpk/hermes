package ssg

import (
	"github.com/adrianpk/hermes/internal/am"
)

func NewWebRouter(handler *WebHandler, opts ...am.Option) *am.Router {
	core := am.NewRouter("web-router", opts...)

	// Content routes
	core.Get("/new-content", handler.NewContent)
	core.Post("/create-content", handler.CreateContent)
	// core.Get("/show-content", handler.ShowContent)
	// core.Get("/edit-content", handler.EditContent)
	// core.Post("/update-content", handler.UpdateContent)
	// core.Get("/list-contents", handler.ListContents)
	// core.Post("/delete-content", handler.DeleteContent)

	return core
}
