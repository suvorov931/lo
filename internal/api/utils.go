package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"lo/internal/logger"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, as *logger.AsyncLogger, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := errorResponse{
		Code:    code,
		Message: message,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		as.Error("WriteError: failed to encode response", slog.String("error", err.Error()))
	}
}
