package middlewares

import (
	"log"
	"net/http"
)

// LoggingMW middleware constructor
func LoggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", "LOG", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		// TODO: добавить ответ в логи (хэдеры и тело ответа)
	})
}
