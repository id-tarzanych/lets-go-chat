package middlewares

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	logger logrus.FieldLogger
}

func NewLogMiddleware(logger logrus.FieldLogger) *LogMiddleware {
	return &LogMiddleware{logger}
}

type LoggingResponseWriterInterface interface {
	http.ResponseWriter
	http.Hijacker
}

type LoggingResponseWriter struct {
	LoggingResponseWriterInterface

	statusCode int
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.LoggingResponseWriterInterface.WriteHeader(code)
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	lrw, _ := w.(LoggingResponseWriterInterface)

	return &LoggingResponseWriter{lrw, http.StatusOK}
}

func (lrw *LoggingResponseWriter) Status() int {
	return lrw.statusCode
}

func (l *LogMiddleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		l.logger.Println(fmt.Sprintf("%q", dump))

		next.ServeHTTP(w, r)
	})
}

func (l *LogMiddleware) LogPanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				l.logger.Errorln("Recovered from panic!", err)

				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (l *LogMiddleware) LogError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		errorCategory := math.Floor(float64(lrw.Status()) / 100)

		if errorCategory == 4 || errorCategory == 5 {
			l.logger.Println(fmt.Sprintf("Non-successful response! Status code: %d", lrw.Status()))
		}
	})
}
