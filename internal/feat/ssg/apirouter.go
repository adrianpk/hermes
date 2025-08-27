package ssg

import (
	"github.com/adrianpk/hermes/internal/am"
)

func NewAPIRouter(handler *APIHandler, mw []am.Middleware, opts ...am.Option) *am.Router {
	core := am.NewAPIRouter("api-router", opts...)
	core.SetMiddlewares(mw)

	// Layout API routes
	core.Get("/layouts", handler.GetAllLayouts)
	core.Get("/layouts/{id}", handler.GetLayout)
	core.Post("/layouts", handler.CreateLayout)
	core.Put("/layouts/{id}", handler.UpdateLayout)
	core.Delete("/layouts/{id}", handler.DeleteLayout)

	// Section API routes
	core.Get("/sections", handler.GetAllSections)
	core.Get("/sections/{id}", handler.GetSection)
	core.Post("/sections", handler.CreateSection)
	core.Put("/sections/{id}", handler.UpdateSection)
	core.Delete("/sections/{id}", handler.DeleteSection)

	// Content API routes
	core.Get("/contents", handler.GetAllContent)
	core.Get("/contents/{id}", handler.GetContent)
	core.Post("/contents", handler.CreateContent)
	core.Put("/contents/{id}", handler.UpdateContent)
	core.Delete("/contents/{id}", handler.DeleteContent)

	return core
}
