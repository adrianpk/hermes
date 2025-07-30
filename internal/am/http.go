package am

import "net/http"

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
