package am

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

type HTTPMethods struct {
	GET    string
	POST   string
	PUT    string
	PATCH  string
	DELETE string
	HEAD   string
}

var HTTPMethod = HTTPMethods{
	GET:    "GET",
	POST:   "POST",
	PUT:    "PUT",
	PATCH:  "PATCH",
	DELETE: "DELETE",
	HEAD:   "HEAD",
}

// PathID extracts a UUID from the request's path values based on the provided key.
func PathID(r *http.Request, key string) (uuid.UUID, error) {
	idStr := r.PathValue(key)
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("ID '%s' is missing in the URL", key)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid ID format for '%s': %w", key, err)
	}

	return id, nil
}
