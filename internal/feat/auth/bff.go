package auth

import (
	"github.com/adrianpk/hermes/internal/am"
)

const (
	// WIP: This will be obtained from configuration.
	defaultAPIBaseURL = "http://localhost:8081/api/v1/auth"
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
