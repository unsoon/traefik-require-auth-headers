package traefik_require_auth_headers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/unsoon/traefik-require-auth-headers/helpers"
)

type Config struct {
	RequiredHeaders []string          `json:"requiredHeaders,omitempty"`
	ErrorResponse   ErrorResponseSpec `json:"errorResponse,omitempty"`
}

type ErrorResponseSpec struct {
	Headers     map[string]string `json:"headers,omitempty"`
	StatusCode  int               `json:"statusCode,omitempty"`
	ContentType string            `json:"contentType,omitempty"`
	Body        *interface{}      `json:"body,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		RequiredHeaders: []string{},
		ErrorResponse: ErrorResponseSpec{
			Headers:     map[string]string{},
			StatusCode:  http.StatusUnauthorized,
			ContentType: "text/plain",
		},
	}
}

type RequireAuthHeaders struct {
	next            http.Handler
	requiredHeaders []string
	errorResponse   ErrorResponseSpec
	name            string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.RequiredHeaders) == 0 {
		return nil, fmt.Errorf("requiredHeaders must not be empty")
	}

	return &RequireAuthHeaders{
		next:            next,
		requiredHeaders: config.RequiredHeaders,
		errorResponse:   config.ErrorResponse,
		name:            name,
	}, nil
}

func (h *RequireAuthHeaders) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")

	if authHeader == "" {
		h.writeErrorResponse(rw)
		return
	}

	prefix := strings.Split(authHeader, " ")[0]

	if prefix == "" {
		h.writeErrorResponse(rw)
		return
	}

	hasValidPrefix := false

	for _, header := range h.requiredHeaders {
		if prefix == header {
			hasValidPrefix = true
			break
		}
	}

	if !hasValidPrefix {
		h.writeErrorResponse(rw)
		return
	}

	h.next.ServeHTTP(rw, req)
}

func (h *RequireAuthHeaders) writeErrorResponse(rw http.ResponseWriter) {
	for key, value := range h.errorResponse.Headers {
		rw.Header().Set(key, value)
	}

	rw.Header().Set("Content-Type", h.errorResponse.ContentType)
	rw.WriteHeader(h.errorResponse.StatusCode)

	if h.errorResponse.Body != nil {
		body, err := helpers.ConvertToType(*h.errorResponse.Body, h.errorResponse.ContentType)

		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Write([]byte(body))
	}
}
