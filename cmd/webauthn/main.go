package main

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/strangnet/go-webauthn-login/internal/pkg/http"
	"github.com/strangnet/go-webauthn-login/internal/pkg/inmem"
	"github.com/strangnet/go-webauthn-login/internal/pkg/user"
	"os"
	"os/signal"
	"syscall"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/webauthn"
)

const (
	defaultApiAddress    = ":80"
	defaultRPDisplayName = "Foo Bar inc"
	defaultRPID          = "localhost"
	defaultRPOrigin      = "http://localhost"
	defaultRPIcon        = "https://duo.com/assets/img/duoLogo-web.png"
)

func main() {

	var (
		apiAddress = envString("API_ADDRESS", defaultApiAddress)

		RPDisplayName = envString("RP_DISPLAY_NAME", defaultRPDisplayName)
		RPID          = envString("RP_ID", defaultRPID)
		RPOrigin      = envString("RP_ORIGIN", defaultRPOrigin)
		RPIcon        = envString("RP_ICON", defaultRPIcon)

		webAuthn     *webauthn.WebAuthn
		sessionStore *session.Store
	)

	flag.Parse()

	logger := *log.StandardLogger()
	logger.SetFormatter(&log.JSONFormatter{})

	errorChannel := make(chan error)

	var err error
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: RPDisplayName,
		RPID:          RPID,
		RPOrigin:      RPOrigin,
		RPIcon:        RPIcon,
		// AttestationPreference:  "",
		// AuthenticatorSelection: protocol.AuthenticatorSelection{},
		// Timeout:                0,
		// Debug:                  false,
	})
	if err != nil {
		log.Fatal("failed to create Webauthn from config: ", err)
	}

	sessionStore, err = session.NewStore()
	if err != nil {
		log.Fatal("failed to create session store", err)
	}

	// Repositories
	var (
		users = inmem.NewUserRepository()
	)

	// Services
	var (
		us = user.NewService(users)
	)

	go func() {
		log.WithField("addr", apiAddress).Info("http.server.listen")

		server := http.NewServer(
			apiAddress,
			logger,
			us,
			webAuthn,
			sessionStore,
		)

		errorChannel <- server.Open()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errorChannel <- fmt.Errorf("got signal: %s", <-c)
	}()

	if err := <-errorChannel; err != nil {
		log.Error(errors.New("got error: " + err.Error()))
	}

	log.Info("terminated")
}

func Hello() string {
	return "Hello, World!"
}

func envString(key, defaultValue string) string {
	value, ok := syscall.Getenv(key)
	if !ok {
		return defaultValue
	}
	return value
}
