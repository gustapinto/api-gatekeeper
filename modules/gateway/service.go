package gateway

import (
	"io"
	"net/http"
	"time"
)

func doRequestToBackendRoute(userId string, service Backend, route Route, body io.ReadCloser) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(route.TimeoutSeconds) * time.Second,
	}
	defer client.CloseIdleConnections()

	request, err := http.NewRequest(route.Method, service.Host+"/"+route.BackendPath, body)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	if service.Headers != nil {
		for key, value := range route.Headers {
			request.Header.Add(key, value)
		}
	}

	if route.Headers != nil {
		for key, value := range route.Headers {
			request.Header.Add(key, value)
		}
	}

	request.Header.Add("X-Api-Gatekeeper-User", userId)

	return client.Do(request)
}
