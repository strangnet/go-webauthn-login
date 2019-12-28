package http

import (
	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/go-chi/chi"
	"github.com/strangnet/go-webauthn-login/internal/pkg/domain"
	"github.com/strangnet/go-webauthn-login/internal/pkg/user"
	"net/http"
	"strings"
)

type registerHandler struct {
	encoder *encoder
	users user.Service
	webAuthn *webauthn.WebAuthn
	sessionStore *session.Store
}

func newRegisterHandler(encoder *encoder, us user.Service, webAuthn *webauthn.WebAuthn, sessionStore *session.Store) *registerHandler {
	return &registerHandler{
		encoder:      encoder,
		users:        us,
		webAuthn:     webAuthn,
		sessionStore: sessionStore,
	}
}

func (h *registerHandler) Routes(r chi.Router) {
	r.Get("/list", h.listRegistrations)
	r.Get("/begin/{username}", h.beginRegistration)
	r.Post("/finish/{username}", h.finishRegistration)
}

func (h *registerHandler) beginRegistration(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	username := chi.URLParam(request, "username")

	// Check if username exists
	u, err := h.users.FindByUsername(username)
	if err != nil {
		displayName := strings.Split(username, "@")[0]
		u = domain.NewUser(username, displayName)
		h.users.Create(u)
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = u.CredentialExcludeList()
	}

	// generate PublicKeyCredentialCreationOptions and session data
	options, sessionData, err := h.webAuthn.BeginRegistration(u, registerOptions)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
		return
	}

	// store session data as json
	err = h.sessionStore.SaveWebauthnSession("registration", sessionData, request, writer)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
		return
	}

	h.encoder.StatusResponse(writer, options, http.StatusOK)
}

func (h *registerHandler) finishRegistration(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	username := chi.URLParam(request, "username")

	// Check if username exists
	u, err := h.users.FindByUsername(username)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
		return
	}

	// load the session data
	sessionData, err := h.sessionStore.GetWebauthnSession("registration", request)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
		return
	}

	credential, err := h.webAuthn.FinishRegistration(u, sessionData, request)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
		return
	}

	u.AddCredential(*credential)

	err = h.users.Create(u)
	if err != nil {
		h.encoder.Error(ctx, writer, err)
	}

	h.encoder.StatusResponse(writer, "Registration Success", http.StatusOK)
}

func (h *registerHandler) listRegistrations(writer http.ResponseWriter, request *http.Request) {
	// ctx := request.Context()


	// Check if username exists
	users := h.users.ListAllUsers()

	h.encoder.StatusResponse(writer, users, http.StatusOK)
}
