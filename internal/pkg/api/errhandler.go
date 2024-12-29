package api

import (
	"errors"
	"log/slog"
	"net/http"

	model "github.com/ShmaykhelDuo/battler/internal/model/api"
)

type APIErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	ID      int    `json:"id"`
	Kind    string `json:"kind"`
	Message string `json:"message,omitempty"`
}

func HandleError(w http.ResponseWriter, err error) {
	var apiError model.Error
	if !errors.As(err, &apiError) {
		slog.Error("unhandled error", "err", err)

		apiError = model.Error{
			Kind: model.KindInternal,
		}
	}

	code := statusCode(apiError)
	res := APIError{
		ID:      apiError.Kind.ID,
		Kind:    apiError.Kind.Description,
		Message: apiError.Message,
	}
	WriteJSONResponse(w, code, APIErrorResponse{Error: res})
}

func statusCode(err model.Error) int {
	switch err.Kind {
	case model.KindInvalidRequest, model.KindInvalidArgument:
		return http.StatusBadRequest
	case model.KindNotFound:
		return http.StatusNotFound
	case model.KindAlreadyExists:
		return http.StatusConflict
	case model.KindInvalidCredentials, model.KindUnauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
