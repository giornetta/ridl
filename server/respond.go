package server

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// respond serializes the given data as JSON and sends it as an http Response
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

// e takes a string and returns a map with key="error"
// This should be called before sending an error data to respond, in order to correctly send an error message
func e(err string) map[string]string {
	return map[string]string{
		"error": err,
	}
}
