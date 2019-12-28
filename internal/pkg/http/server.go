package http

import (
	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/strangnet/go-webauthn-login/internal/pkg/user"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

type Server interface {
	Open() error
	Close()
	Handler() http.Handler
}

type server struct {
	listener     net.Listener
	addr         string
	us           user.Service
	webAuthn     *webauthn.WebAuthn
	sessionStore *session.Store
	logger       log.Logger
	encoder      *encoder
}

func NewServer(addr string, logger log.Logger, us user.Service, webAuthn *webauthn.WebAuthn, sessionStore *session.Store) Server {
	return &server{
		addr:         addr,
		logger:       logger,
		us:           us,
		webAuthn:     webAuthn,
		sessionStore: sessionStore,
		encoder:      NewEncoder(logger),
	}
}

func (s *server) Open() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.listener = listener
	server := http.Server{
		Handler: s.Handler(),
	}

	return server.Serve(s.listener)
}

func (s *server) Close() {
	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		workDir, _ := os.Getwd()
		filesDir := filepath.Join(workDir, "web")
		fs := http.StripPrefix("/", http.FileServer(http.Dir(filesDir)))
		fs.ServeHTTP(w, req)
	}))

	r.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Route("/register", newRegisterHandler(s.encoder, s.us, s.webAuthn, s.sessionStore).Routes)
			r.Route("/login", newLoginHandler(s.encoder, s.us, s.webAuthn, s.sessionStore).Routes)
		})
	})

	return r
}
