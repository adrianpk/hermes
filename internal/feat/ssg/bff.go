package ssg

import (
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/adrianpk/hermes/internal/feat/auth"
	"github.com/google/uuid"
)

const (
	// WIP: This will be obtained from configuration.
	defaultAPIBaseURL = "http://localhost:8081/api/v1/ssg"
)

type BFF struct {
	*am.WebHandler
	apiClient *am.APIClient
}

func NewBFF(tm *am.TemplateManager, flash *am.FlashManager, opts ...am.Option) *BFF {
	handler := am.NewWebHandler(tm, flash, opts...)
	apiClient := am.NewAPIClient("bff-api-client", defaultAPIBaseURL, opts...)
	return &BFF{
		WebHandler: handler,
		apiClient:  apiClient,
	}
}

func (h *BFF) sampleUserInSession(r *http.Request) auth.User {
	user := auth.NewUser("fakeuser", "Fake User")
	user.SetID(uuid.MustParse("00000000-0000-0000-0000-000000000001"), true)
	user.IsActive = true
	return user
}
