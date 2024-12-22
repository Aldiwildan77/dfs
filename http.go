package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody []byte
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error().Err(err).Msg("Error reading request body")
				WriteJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Invalid request body"})
				return
			}
			requestBody = bodyBytes

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Str("referer", r.Referer()).
			Interface("headers", r.Header).
			Interface("cookies", r.Cookies()).
			RawJSON("body", requestBody).
			Msg("Incoming http request")

		next.ServeHTTP(w, r)
	})
}
