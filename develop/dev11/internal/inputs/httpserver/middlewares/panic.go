package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"
)

// PanicMW middleware constructor
func PanicMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				wrapHTTPError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Printf("%s %s", "PANIC", string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
