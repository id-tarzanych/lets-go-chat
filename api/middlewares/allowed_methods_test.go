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
		{
			name:     "GET method",
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			wantCode: http.StatusOK,
		},
		{
			name:     "POST method",
			req:      httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}")),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "PUT method",
			req:      httptest.NewRequest(http.MethodPut, "/test", strings.NewReader("{}")),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "DELETE method",
			req:      httptest.NewRequest(http.MethodDelete, "/test", nil),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "OPTIONS method",
			req:      httptest.NewRequest(http.MethodOptions, "/test", nil),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "HEAD method",
			req:      httptest.NewRequest(http.MethodHead, "/test", nil),
			wantCode: http.StatusBadRequest,
		},
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
		{
			name:     "GET method",
			req:      httptest.NewRequest(http.MethodGet, "/test", nil),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "POST method",
			req:      httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}")),
			wantCode: http.StatusOK,
		},
		{
			name:     "PUT method",
			req:      httptest.NewRequest(http.MethodPut, "/test", strings.NewReader("{}")),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "DELETE method",
			req:      httptest.NewRequest(http.MethodDelete, "/test", nil),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "OPTIONS method",
			req:      httptest.NewRequest(http.MethodOptions, "/test", nil),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "HEAD method",
			req:      httptest.NewRequest(http.MethodHead, "/test", nil),
			wantCode: http.StatusBadRequest,
		},
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
