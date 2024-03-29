package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	hijacktest "github.com/getlantern/httptest"
	"github.com/stretchr/testify/mock"

	"github.com/id-tarzanych/lets-go-chat/internal/testserver"
	"github.com/id-tarzanych/lets-go-chat/mocks"
)

func TestLogMiddleware_LogError(t *testing.T) {
	s := &testserver.Server{}

	tests := []struct {
		name       string
		req        *http.Request
		called     bool
		statusCode int
	}{
		{
			name:       "Bad Request",
			req:        httptest.NewRequest(http.MethodGet, "/badRequest", nil),
			called:     true,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Internal Server Error",
			req:        httptest.NewRequest(http.MethodGet, "/internalServerError", nil),
			called:     true,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Valid Request",
			req:        httptest.NewRequest(http.MethodGet, "/test", nil),
			statusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := hijacktest.NewRecorder(nil)

			loggerMock := &mocks.FieldLogger{}
			loggerMock.On("Println", mock.AnythingOfType("string")).Return()

			l := &LogMiddleware{
				logger: loggerMock,
			}

			l.LogError(s).ServeHTTP(w, tt.req)

			if tt.called {
				loggerMock.AssertCalled(t, "Println", fmt.Sprintf("Non-successful response! Status code: %d", tt.statusCode))
			} else {
				loggerMock.AssertNotCalled(t, "Println", fmt.Sprintf("Non-successful response! Status code: %d", tt.statusCode))
			}
		})
	}
}

func TestLogMiddleware_LogPanicRecovery(t *testing.T) {
	s := &testserver.Server{}

	tests := []struct {
		name     string
		req      *http.Request
		panicked bool
	}{
		{
			name:     "Panic",
			req:      httptest.NewRequest(http.MethodGet, "/panic", nil),
			panicked: true,
		},
		{
			name: "Bad Request",
			req:  httptest.NewRequest(http.MethodGet, "/badRequest", nil),
		},
		{
			name: "Internal Server Error",
			req:  httptest.NewRequest(http.MethodGet, "/internalServerError", nil),
		},
		{
			name: "Valid Request",
			req:  httptest.NewRequest(http.MethodGet, "/test", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := hijacktest.NewRecorder(nil)

			loggerMock := &mocks.FieldLogger{}
			loggerMock.On("Errorln", mock.AnythingOfType("string"), mock.Anything).Return()

			l := &LogMiddleware{
				logger: loggerMock,
			}

			l.LogPanicRecovery(s).ServeHTTP(w, tt.req)

			if tt.panicked {
				loggerMock.AssertCalled(t, "Errorln", "Recovered from panic!", mock.Anything)
			} else {
				loggerMock.AssertNotCalled(t, "Errorln", "Recovered from panic!", mock.Anything)
			}
		})
	}
}

func TestLogMiddleware_LogRequest(t *testing.T) {
	s := &testserver.Server{}

	tests := []struct {
		name string
		req  *http.Request
	}{
		{
			name: "Bad Request",
			req:  httptest.NewRequest(http.MethodGet, "/badRequest", nil),
		},
		{
			name: "Internal Server Error",
			req:  httptest.NewRequest(http.MethodGet, "/internalServerError", nil),
		},
		{
			name: "Valid Request",
			req:  httptest.NewRequest(http.MethodGet, "/test", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := hijacktest.NewRecorder(nil)

			loggerMock := &mocks.FieldLogger{}
			loggerMock.On("Println", mock.AnythingOfType("string")).Return()

			l := &LogMiddleware{
				logger: loggerMock,
			}

			l.LogRequest(s).ServeHTTP(w, tt.req)

			dump, _ := httputil.DumpRequest(tt.req, true)
			loggerMock.AssertCalled(t, "Println", fmt.Sprintf("%q", dump))
		})
	}
}
