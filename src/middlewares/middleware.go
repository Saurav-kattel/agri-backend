package middlewares

import (
	"log"
	"net/http"
	"time"
)

type writerWrapper struct {
	http.ResponseWriter
	Status int
}

type Middleware func(http.Handler) http.Handler

func (w *writerWrapper) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.Status = statusCode

}

func LoggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		status := &writerWrapper{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(status, r)
		log.Println(status.Status, r.Method, r.URL.Path, time.Since(start))
	})
}

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
