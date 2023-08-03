package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func wrapLoggingError(err error) {
	log.Printf("%s %v", "ERROR_WRAPPER", err)
}

func wrapHTTPError(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)

	bytes, err := json.Marshal(ErrorResp{Error: error})
	if err != nil {
		wrapLoggingError(fmt.Errorf("can't marshal ErrorResp: %v", err))
		http.Error(w, "json marshal error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		wrapLoggingError(fmt.Errorf("can't write bytes to response writer in wrapHTTPError: %v", err))
		http.Error(w, "writing response body error", http.StatusInternalServerError)
		return
	}
}
