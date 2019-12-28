package http

import (
	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/go-chi/chi"
	"github.com/strangnet/go-webauthn-login/internal/pkg/user"
	"net/http"
)

type loginHandler struct {
	encoder      *encoder
	users        user.Service
	webAuthn     *webauthn.WebAuthn
	sessionStore *session.Store
}

func newLoginHandler(encoder *encoder, us user.Service, webAuthn *webauthn.WebAuthn, sessionStore *session.Store) *loginHandler {
	return &loginHandler{
		encoder:      encoder,
		users:        us,
		webAuthn:     webAuthn,
		sessionStore: sessionStore,
	}
}

func (h *loginHandler) Routes(r chi.Router) {
	r.Get("/begin/{username}", h.beginLogin)
	r.Post("/finish/{username}", h.finishLogin)
}

func (h *loginHandler) beginLogin(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	username := chi.URLParam(request, "username")

	// Check if username exists
	u, err := h.users.FindByUsername(username)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
		return
	}

	// begin webauthn login
	options, sessionData, err := h.webAuthn.BeginLogin(u)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
		return
	}

	// store session data
	err = h.sessionStore.SaveWebauthnSession("authentication", sessionData, request, writer)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
	}

	h.encoder.StatusResponse(writer, options, http.StatusOK)
}

func (h *loginHandler) finishLogin(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	username := chi.URLParam(request, "username")

	// Check if username exists
	u, err := h.users.FindByUsername(username)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
		return
	}

	// load the session data
	sessionData, err := h.sessionStore.GetWebauthnSession("authentication", request)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
		return
	}

	// Light weight credential check, i.e. none
	_, err = h.webAuthn.FinishLogin(u, sessionData, request)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
		return
	}

	h.encoder.StatusResponse(writer, "Login success", http.StatusOK)
}
