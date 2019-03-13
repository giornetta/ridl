package server

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/giornetta/ridl"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
)

func Router(svc ridl.Service) *chi.Mux {
	h := &ridlHandler{svc}

	// Create an HTTP multiplexer
	mux := chi.NewMux()

	// add the needed middlewares to the multiplexer
	mux.Use(
		allowCORS,
		middleware.Logger,
		middleware.StripSlashes,
		middleware.Recoverer,
	)

	mux.Mount("/ridl", h.routes())

	return mux
}

// New returns an *http.Server that can correctly handle requests
func New(mux http.Handler, port string) *http.Server {
	// tlsConfig contains the best settings to correctly serve over the web in a secure way.
	// TLS certificates should be added in order to use HTTPS.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	s := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 120,
		TLSConfig:    tlsConfig,
		Handler:      mux,
	}

	return s
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				headers := []string{"Content-Type", "Accept", "Authorization"}
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
				methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
