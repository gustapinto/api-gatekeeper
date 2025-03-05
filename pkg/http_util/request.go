package httputil

import (
	"net/http"
)

func GetHeadersAsMap(r *http.Request) map[string]string {
	if r == nil {
		return nil
	}

	headers := make(map[string]string)
	for key, values := range r.Header {
		for _, value := range values {
			headers[key] = value
		}
	}

	return headers
}

func GetQueryParamsAsMap(r *http.Request) map[string]string {
	if r == nil {
		return nil
	}

	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		for _, value := range values {
			queryParams[key] = value
		}
	}

	return queryParams
}
