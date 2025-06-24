package auth

import (
	"net/http"

	"github.com/adrianpk/hermes/internal/am"
)

func GetFlash(r *http.Request) am.Flash {
	flash, ok := r.Context().Value(am.FlashKey).(*am.Flash)
	if ok {
		return *flash
	}

	return am.Flash{
		Notifications: []am.Notification{},
	}
}
