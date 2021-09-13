package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

const jsonContentType = "application/json"

func FillFromRequest(req *http.Request, destination interface{}) (url.Values, error) {
	if destination != nil && strings.ToLower(req.Header.Get("content-type")) == jsonContentType {
		err := json.NewDecoder(req.Body).Decode(destination)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, errors.Wrap(err, "error parsing http request body")
		}
	}
	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	result := map[string][]string{}

	// first fill variables from query string
	urlQuery := req.URL.Query()
	for name, param := range urlQuery {
		result[name] = param
	}

	// then fill variables from POST data
	postVars := req.Form
	for name, param := range postVars {
		result[name] = param
	}

	// and then fill variables from URL
	routeCtx := chi.RouteContext(req.Context())
	for key, name := range routeCtx.URLParams.Keys {
		result[name] = []string{routeCtx.URLParams.Values[key]}
	}

	// fill variables from headers
	for name, param := range req.Header {
		// Authorize becomes authorize
		result[strings.ToLower(name)] = param
	}

	return url.Values(result), nil
}
