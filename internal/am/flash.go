package am

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/securecookie"
)

type FlashCtxKey string

const (
	FlashKey        FlashCtxKey = "aqmflash"
	FlashCookieName string      = "aqmflash"
)

type notificationTypes struct {
	Success string
	Info    string
	Warn    string
	Error   string
	Debug   string
}

var NotificationType = notificationTypes{
	Success: "success",
	Info:    "info",
	Warn:    "warning",
	Error:   "danger",
	Debug:   "debug",
}

type Notification struct {
	Type string
	Msg  string
}

type Flash struct {
	Notifications []Notification
}

func NewFlash() Flash {
	return Flash{
		Notifications: []Notification{},
	}
}

func (f *Flash) Add(typ, msg string) {
	f.Notifications = append(f.Notifications, Notification{
		Type: typ,
		Msg:  msg,
	})
}

func (f *Flash) Clear() {
	f.Notifications = []Notification{}
}

func (f *Flash) HasMessages() bool {
	return len(f.Notifications) > 0
}

type FlashManager struct {
	Core
	encoder *securecookie.SecureCookie
	flashes sync.Map
}

func NewFlashManager(opts ...Option) *FlashManager {
	core := NewCore("fm-manager", opts...)
	return &FlashManager{
		Core: core,
	}
}

func (fm *FlashManager) Setup(ctx context.Context) error {
	err := fm.Core.Setup(ctx)
	if err != nil {
		return err
	}

	cfg := fm.Cfg()

	hashKey := cfg.ByteSliceVal(Key.SecHashKey)
	blockKey := cfg.ByteSliceVal(Key.SecBlockKey)

	if len(hashKey) == 0 || len(blockKey) == 0 {
		return errors.New("missing hashKey or blockKey in configuration")
	}

	fm.encoder = securecookie.New(hashKey, blockKey)
	return nil
}

func (fm *FlashManager) SetFlashInCookie(w http.ResponseWriter, flash Flash) {
	if fm.encoder == nil {
		return
	}

	encoded, err := fm.encoder.Encode(FlashCookieName, flash)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     FlashCookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})
}

func (fm *FlashManager) GetFlashFromCookie(r *http.Request) Flash {
	if fm.encoder == nil {
		return NewFlash()
	}
	cookie, err := r.Cookie(FlashCookieName)
	if err != nil {
		return NewFlash()
	}
	var flash Flash
	err = fm.encoder.Decode(FlashCookieName, cookie.Value, &flash)
	if err != nil {
		return NewFlash()
	}
	return flash
}

func (fm *FlashManager) ClearFlashCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     FlashCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
	})
}

func (fm *FlashManager) SaveFlash(r *http.Request, flash Flash) {
	key := ReqID(r)
	fm.flashes.Store(key, flash)
}

func (fm *FlashManager) GetFlash(r *http.Request) Flash {
	key := ReqID(r)
	if v, ok := fm.flashes.Load(key); ok {
		return v.(Flash)
	}
	return NewFlash()
}

func (fm *FlashManager) ClearFlash(w http.ResponseWriter, r *http.Request) {
	key := ReqID(r)
	fm.flashes.Delete(key)
	fm.ClearFlashCookie(w)
}

func (fm *FlashManager) AddFlash(r *http.Request, typ, msg string) {
	flash := fm.GetFlash(r)
	flash.Add(typ, msg)
	fm.SaveFlash(r, flash)
}

func (fm *FlashManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flash := fm.GetFlashFromCookie(r)
		fm.SaveFlash(r, flash)
		defer fm.ClearFlash(w, r)
		next.ServeHTTP(w, r)
	})
}

type flashResponseWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func (w *flashResponseWriter) WriteHeader(statusCode int) {
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *flashResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

type responseWriter struct {
	http.ResponseWriter
	flash         Flash
	flashModified bool
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.flashModified = true
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) AddFlash(typ, msg string) {
	rw.flash.Add(typ, msg)
	rw.flashModified = true
}

func (rw *responseWriter) GetFlash() Flash {
	return rw.flash
}

func (fm *FlashManager) Middlewares() []Middleware {
	return []Middleware{fm.Middleware}
}
