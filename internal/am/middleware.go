package am

import (
	"context"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

const defaultCSRFKey = "set-a-csrf-key!"

var (
	csrfMiddleware func(http.Handler) http.Handler
	once           sync.Once
)

// LogHeadersMw is a middleware that logs all request headers.
func LogHeadersMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := NewLogger("request-headers")
		log.Info("Incoming Request Headers:")
		for name, headers := range r.Header {
			for _, h := range headers {
				log.Infof("  %s: %s", name, h)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// MethodOverrideMw is a middleware that checks for a _method form field and overrides the request method.
func MethodOverrideMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if override := r.FormValue("_method"); override != "" {
				r.Method = override
			}
		}
		next.ServeHTTP(w, r)
	})
}

// CSRFMw is a middleware that protects against CSRF attacks.
func CSRFMw(cfg *Config) func(http.Handler) http.Handler {
	if cfg == nil {
		return passThroughMw
	}

	initCSRF(cfg)

	return func(next http.Handler) http.Handler {
		return csrfMiddleware(next)
	}
}

func passThroughMw(next http.Handler) http.Handler {
	return next
}

func initCSRF(cfg *Config) {
	once.Do(func() {
		key := cfg.StrValOrDef(Key.SecCSRFKey, defaultCSRFKey)
		to := cfg.StrValOrDef(Key.SecCSRFRedirect, "/csrf-error")

		csrfMiddleware = csrf.Protect(
			[]byte(key),
			csrf.FieldName(CSRFFieldName),
			csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, to, http.StatusFound)
			})),
		)
	})
}

// ReqIDKey is the context key for the request ID.
const ReqIDKey = "requestID"

// RequestIDMw is a middleware that assigns a unique ID to each request and stores it in the context and as a header.
func RequestIDMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		ctx := context.WithValue(r.Context(), ReqIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ReqID returns the request ID from the context, or an empty string if not set.
func ReqID(r *http.Request) string {
	if v := r.Context().Value(ReqIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
