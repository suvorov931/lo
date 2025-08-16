package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"lo/internal/logger"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, logger logger.Logger, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := ErrorResponse{
		Code:    code,
		Message: message,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error("WriteError: failed to encode response", slog.String("error", err.Error()))
	}
}
