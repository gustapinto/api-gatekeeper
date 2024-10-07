package util

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteMethodNotAllowed(w http.ResponseWriter) {
	errorJson, e := json.Marshal(ErrorResponse{
		Message: "Method not allowed",
	})
	if e != nil {
		errorJson = []byte("{}")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(errorJson)
}

func WriteInternalServerError(w http.ResponseWriter, err error) {
	errorJson, e := json.Marshal(ErrorResponse{
		Message: err.Error(),
	})
	if e != nil {
		errorJson = []byte("{}")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(errorJson)
}

func WriteBadRequest(w http.ResponseWriter, err error) {
	errorJson, e := json.Marshal(ErrorResponse{
		Message: err.Error(),
	})
	if e != nil {
		errorJson = []byte("{}")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(errorJson)
}

func WriteUnauthorized(w http.ResponseWriter) {
	errorJson, e := json.Marshal(ErrorResponse{
		Message: "Unauthorized",
	})
	if e != nil {
		errorJson = []byte("{}")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(errorJson)
}

func WriteForbidden(w http.ResponseWriter) {
	errorJson, e := json.Marshal(ErrorResponse{
		Message: "You do not have the scopes to access this resource",
	})
	if e != nil {
		errorJson = []byte("{}")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write(errorJson)
}
