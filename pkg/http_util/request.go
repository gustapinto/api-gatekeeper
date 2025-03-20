package httputil

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
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

func ParseBasicAuthorizationToken(token string) (string, string, error) {
	if token == "" {
		return "", "", errors.New("badparams: missing Authorization token")
	}

	if strings.Contains(token, "Basic") {
		token = strings.TrimSpace(strings.ReplaceAll(token, "Basic", ""))
	}

	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", "", err
	}

	data := strings.Split(string(decodedToken), ":")
	if len(data) < 2 {
		return "", "", err
	}
	login := data[0]
	password := data[1]

	return login, password, nil
}
