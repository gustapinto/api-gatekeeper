package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/gustapinto/api-gatekeeper/internal/config"
	"github.com/gustapinto/api-gatekeeper/internal/service"
	httputil "github.com/gustapinto/api-gatekeeper/pkg/http_util"
)

type BackendHandler struct {
	Service service.Backend
}

func (b BackendHandler) HandleBackendRouteRequest(w http.ResponseWriter, r *http.Request, backend config.Backend, route config.Route) {
	uid := r.Context().Value("userId")

	userId := ""
	if uidStr, ok := uid.(string); ok {
		userId = uidStr
	}

	if strings.ToUpper(r.Method) != route.Method {
		httputil.WriteMethodNotAllowed(w)
		return
	}

	response, err := b.Service.DoRequestToBackendRoute(userId, backend, route, r.Body)
	if err != nil {
		httputil.WriteInternalServerError(w, err)
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		httputil.WriteInternalServerError(w, err)
		return
	}

	responseContentType := response.Header.Get("Content-Type")
	if responseContentType != "" {
		w.Header().Add("Content-Type", responseContentType)
	}

	w.WriteHeader(response.StatusCode)
	w.Write(responseBody)
}
