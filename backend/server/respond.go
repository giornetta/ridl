package server

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func respond(w http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		respond(w, http.StatusInternalServerError, e(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func e(err string) map[string]string {
	return map[string]string{
		"error": err,
	}
}
