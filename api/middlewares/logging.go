package middlewares

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	logger *logrus.Logger
}

func NewLogMiddleware(logger *logrus.Logger) *LogMiddleware {
	return &LogMiddleware{logger}
}

type LoggingResponseWriterInterface interface {
	http.ResponseWriter
	http.Hijacker
}

type LoggingResponseWriter struct {
	ResponseWriter LoggingResponseWriterInterface
	statusCode int
}

func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.ResponseWriter.Header()
}

func (lrw *LoggingResponseWriter) Write(bytes []byte) (int, error) {
	return lrw.ResponseWriter.Write(bytes)
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error)  {
	return lrw.ResponseWriter.Hijack()
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
			logrus.Println(fmt.Sprintf("Non-successful response! Status code: %d", lrw.Status()))
		}
	})
}