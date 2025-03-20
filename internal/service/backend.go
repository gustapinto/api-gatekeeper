package service

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gustapinto/api-gatekeeper/internal/config"
)

type Backend struct{}

func NewBackend() Backend {
	return Backend{}
}

func (Backend) mergeHeaders(headerMaps ...map[string]string) map[string]string {
	headers := make(map[string]string)

	for _, headerMap := range headerMaps {
		for key, value := range headerMap {
			headers[key] = value
		}
	}

	return headers
}

func (b Backend) DoRequestToBackendRoute(
	userId string,
	requestId string,
	backend config.Backend,
	route config.Route,
	body io.ReadCloser,
	requestHeaders map[string]string,
	queryParams map[string]string,
) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(route.TimeoutSeconds) * time.Second,
	}
	defer client.CloseIdleConnections()

	backendPath, err := url.JoinPath(backend.Host, route.BackendPath)
	if err != nil {
		return nil, err
	}

	backendUrl, err := url.Parse(backendPath)
	if err != nil {
		return nil, err
	}

	backendUrlQuery := backendUrl.Query()
	for key, value := range queryParams {
		backendUrlQuery.Add(key, value)
	}

	backendUrl.RawQuery = backendUrlQuery.Encode()

	request, err := http.NewRequest(route.Method, backendUrl.String(), body)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	additionalHeaders := make(map[string]string)
	if backend.PassHeaders || route.PassHeaders {
		additionalHeaders = requestHeaders
	}

	headers := b.mergeHeaders(backend.Headers, route.Headers, additionalHeaders)
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	request.Header.Add("X-Api-Gatekeeper-User", userId)
	request.Header.Add("X-Api-Gatekeeper-Request", userId)

	return client.Do(request)
}
