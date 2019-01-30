package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/giornetta/ridl/ridl"
)

type ridlHandler struct {
	s ridl.Service
}

func (h *ridlHandler) routes() *chi.Mux {
	mux := chi.NewMux()

	mux.Post("/encrypt", h.encrypt)
	mux.Post("/decrypt", h.decrypt)

	return mux
}

// encrypt wraps ridl.Service.Encrypt, decoding the request body and sending it to the service.
func (h *ridlHandler) encrypt(w http.ResponseWriter, r *http.Request) {
	var req ridl.EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond(w, http.StatusBadRequest, e(err.Error()))
		return
	}

	res, err := h.s.Encrypt(&req)
	if err != nil {
		respond(w, http.StatusInternalServerError, e(err.Error()))
		return
	}

	respond(w, http.StatusOK, res)
}

// decrypt wraps ridl.Service.Decrypt, decoding the request body and sending it to the service.
func (h *ridlHandler) decrypt(w http.ResponseWriter, r *http.Request) {
	var req ridl.DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond(w, http.StatusBadRequest, e(err.Error()))
		return
	}

	res, err := h.s.Decrypt(&req)
	if err != nil {
		respond(w, http.StatusBadRequest, e(err.Error()))
		return
	}

	respond(w, http.StatusOK, res)
}
