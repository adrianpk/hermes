package auth

import (
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

const (
	resUserName    = "user"
	resUserNameCap = "User"

	resRoleName    = "role"
	resRoleNameCap = "Role"

	resPermissionName    = "permission"
	resPermissionNameCap = "Permission"

	resResourceName    = "resource"
	resResourceNameCap = "Resource"

	resOrgName    = "org"
	resOrgNameCap = "Org"

	resTeamName    = "team"
	resTeamNameCap = "Team"
)

type APIHandler struct {
	*am.APIHandler
	service Service
	crypto  *am.Crypto
}

func NewAPIHandler(name string, service Service, options ...am.Option) *APIHandler {
	h := am.NewAPIHandler(name, options...)
	crypto := am.NewCrypto(h.Cfg().ByteSliceVal(am.Key.SecEncryptionKey))

	return &APIHandler{
		APIHandler: h,
		service:    service,
		crypto:     crypto,
	}
}

func (h *APIHandler) OK(w http.ResponseWriter, message string, data interface{}) {
	wrappedData := h.wrapData(data)
	h.APIHandler.OK(w, message, wrappedData)
}

func (h *APIHandler) Created(w http.ResponseWriter, message string, data interface{}) {
	wrappedData := h.wrapData(data)
	h.APIHandler.Created(w, message, wrappedData)
}

func (h *APIHandler) wrapData(data interface{}) interface{} {
	switch v := data.(type) {
	// Single entities
	case User:
		return map[string]interface{}{"user": v}
	case Role:
		return map[string]interface{}{"role": v}
	case Permission:
		return map[string]interface{}{"permission": v}
	case Resource:
		return map[string]interface{}{"resource": v}
	case Org:
		return map[string]interface{}{"org": v}
	case Team:
		return map[string]interface{}{"team": v}

	// Slices of entities
	case []User:
		return map[string]interface{}{"users": v}
	case []Role:
		return map[string]interface{}{"roles": v}
	case []Permission:
		return map[string]interface{}{"permissions": v}
	case []Resource:
		return map[string]interface{}{"resources": v}
	case []Org:
		return map[string]interface{}{"orgs": v}
	case []Team:
		return map[string]interface{}{"teams": v}

	// Default case for nil, maps, or other types
	default:
		return data
	}
}
