package utils

import (
	"log"
	"net/http"
)

const (
	ErrCodeInvalidRequest    = "INVALID_REQUEST"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeInternal          = "INTERNAL_ERROR"
)

const internalErrorMessage = "internal server error"

type ErrorBody struct {
	Code    string `json:"code" example:"INVALID_REQUEST"`
	Message string `json:"message" example:"invalid request"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

func RespondWithError(w http.ResponseWriter, httpStatus int, code, message string) {
	RespondWithJSON(w, httpStatus, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func RespondBadRequest(w http.ResponseWriter, message string) {
	RespondWithError(w, http.StatusBadRequest, ErrCodeInvalidRequest, message)
}

func RespondUnauthorized(w http.ResponseWriter, message string) {
	RespondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, message)
}

func RespondForbidden(w http.ResponseWriter, message string) {
	RespondWithError(w, http.StatusForbidden, ErrCodeForbidden, message)
}

func RespondNotFound(w http.ResponseWriter, message string) {
	RespondWithError(w, http.StatusNotFound, ErrCodeNotFound, message)
}

func RespondInternalServerError(w http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("internal server error: %v", err)
	}
	RespondWithError(w, http.StatusInternalServerError, ErrCodeInternal, internalErrorMessage)
}
