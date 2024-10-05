package gateway

import (
	"io"
	"net/http"
	"strings"

	"github.com/gustapinto/api-gatekeeper/util"
)

func handleBackendRouteRequest(userId string, service Backend, route Route, w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) != route.Method {
		util.WriteMethodNotAllowed(w)
		return
	}

	response, err := doRequestToBackendRoute(userId, service, route, r.Body)
	if err != nil {
		util.WriteInternalServerError(w, err)
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		util.WriteInternalServerError(w, err)
		return
	}

	responseContentType := response.Header.Get("Content-Type")
	if responseContentType != "" {
		w.Header().Add("Content-Type", responseContentType)
	}

	w.WriteHeader(response.StatusCode)
	w.Write(responseBody)
}
