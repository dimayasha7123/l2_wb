package middlewares

import "fmt"

// ServeHTTPError struct
type ServeHTTPError struct {
	InternalError error
	Message       string
	Code          int
}

func (e *ServeHTTPError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.InternalError)
}

// ResultResp response struct
type ResultResp struct {
	Result any `json:"result"`
}

// ErrorResp response struct
type ErrorResp struct {
	Error string `json:"error"`
}
