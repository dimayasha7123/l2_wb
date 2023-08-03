package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HandlerWithResults interface
type HandlerWithResults interface {
	ServeHTTP(http.ResponseWriter, *http.Request) (any, *ServeHTTPError)
}

// ResultsWrapper struct
type ResultsWrapper struct {
	handler HandlerWithResults
}

// NewResultsWrapper constructor
func NewResultsWrapper(handler HandlerWithResults) *ResultsWrapper {
	return &ResultsWrapper{handler: handler}
}

func (wr *ResultsWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result, serveErr := wr.handler.ServeHTTP(w, r)
	if serveErr != nil {
		wrapLoggingError(fmt.Errorf("handler serve http error: %v", serveErr))
		wrapHTTPError(w, serveErr.Message, serveErr.Code)
		return
	}

	// accurate!
	if result == nil {
		result = "ok"
	}
	resultResp := ResultResp{Result: result}

	bytes, err := json.Marshal(resultResp)
	if err != nil {
		wrapLoggingError(fmt.Errorf("can't marshal ResultResp: %v", err))
		wrapHTTPError(w, "json marshal error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		wrapLoggingError(fmt.Errorf("can't write bytes to response writer: %v", err))
		wrapHTTPError(w, "writing response body error", http.StatusInternalServerError)
		return
	}
}
