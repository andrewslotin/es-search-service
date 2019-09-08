package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_WithAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic "+base64.RawURLEncoding.EncodeToString([]byte("user1:password2")))

	var numRequests int
	h := AuthMiddleware(func(w http.ResponseWriter, req AuthenticatedRequest) {
		numRequests++
		assert.Equal(t, "user1", req.Username)
		w.Write([]byte("welcome"))
	})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "welcome", rec.Body.String())
	assert.Equal(t, 1, numRequests)
}

func TestAuthMiddleware_WithoutAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	var numRequests int
	h := AuthMiddleware(func(w http.ResponseWriter, req AuthenticatedRequest) {
		numRequests++
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, 0, numRequests)
}
