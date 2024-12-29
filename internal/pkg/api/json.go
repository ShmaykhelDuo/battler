package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, code int, v any) {
	bytes, err := json.Marshal(v)
	if err != nil {
		slog.Error("failed to marshal value", "err", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bytes)
}
