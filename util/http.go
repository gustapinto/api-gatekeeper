package util

import (
	"encoding/json"
	"fmt"
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
		Message: fmt.Sprintf("Failed to build request, got error %s", err.Error()),
	})
	if e != nil {
		errorJson = []byte("{}")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(errorJson)
}
