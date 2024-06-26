package auth

import (
	"errors"
	"net/http"
	"strings"
)

// extracts api key from headers of http request
// example:
// Authorization: ApiKey <api_key>
func GetApiKey(headers http.Header) (string, error) {
	authHeaderValue := headers.Get("Authorization")
	if authHeaderValue == "" {
		return "", errors.New("no authorization key provided")
	}

	headerVals := strings.Split(authHeaderValue, " ")
	if len(headerVals) != 2 {
		return "", errors.New("malformed authorization Header")
	}

	if headerVals[0] != "ApiKey" {
		return "", errors.New("malformed authorization Header")
	}

	return headerVals[1], nil

}
