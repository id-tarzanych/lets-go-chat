package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/id-tarzanych/lets-go-chat/internal/testserver"
)

func TestGetOnly(t *testing.T) {
	s := &testserver.Server{}

	tests := []struct {
		name     string
		req      *http.Request
		wantCode int
	}{
		{"GET method", httptest.NewRequest(http.MethodGet, "/test", nil), http.StatusOK},
		{"POST method", httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}")), http.StatusBadRequest},
		{"PUT method", httptest.NewRequest(http.MethodPut, "/test", strings.NewReader("{}")), http.StatusBadRequest},
		{"DELETE method", httptest.NewRequest(http.MethodDelete, "/test", nil), http.StatusBadRequest},
		{"OPTIONS method", httptest.NewRequest(http.MethodOptions, "/test", nil), http.StatusBadRequest},
		{"HEAD method", httptest.NewRequest(http.MethodHead, "/test", nil), http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			GetOnly(s).ServeHTTP(w, tt.req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("GetOnly() = %v, want %v", w.Result().StatusCode, tt.wantCode)
			}
		})
	}
}

func TestPostOnly(t *testing.T) {
	s := &testserver.Server{}

	tests := []struct {
		name     string
		req      *http.Request
		wantCode int
	}{
		{"GET method", httptest.NewRequest(http.MethodGet, "/test", nil), http.StatusBadRequest},
		{"POST method", httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}")), http.StatusOK},
		{"PUT method", httptest.NewRequest(http.MethodPut, "/test", strings.NewReader("{}")), http.StatusBadRequest},
		{"DELETE method", httptest.NewRequest(http.MethodDelete, "/test", nil), http.StatusBadRequest},
		{"OPTIONS method", httptest.NewRequest(http.MethodOptions, "/test", nil), http.StatusBadRequest},
		{"HEAD method", httptest.NewRequest(http.MethodHead, "/test", nil), http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			PostOnly(s).ServeHTTP(w, tt.req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("PostOnly() = %v, want %v", w.Result().StatusCode, tt.wantCode)
			}
		})
	}
}
