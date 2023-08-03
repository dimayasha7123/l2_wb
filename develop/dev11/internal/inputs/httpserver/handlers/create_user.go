package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"l2_wb/develop/dev11/internal/app"
	"l2_wb/develop/dev11/internal/inputs/httpserver/middlewares"
)

// CreateUserReq request struct
type CreateUserReq struct {
	Nickname string `json:"nickname"`
}

// Validate method
func (r CreateUserReq) Validate() error {
	fields := make([]string, 0, 1)
	if r.Nickname == "" {
		fields = append(fields, "nickname")
	}
	if len(fields) == 0 {
		return nil
	}
	return NewMissingFieldsErr(fields)
}

// CreateUserResp response struct
type CreateUserResp struct {
	UserID int64 `json:"user_id"`
}

// CreateUserHandler handler struct
type CreateUserHandler struct {
	service *app.App
}

// NewCreateUserHandler constructor
func NewCreateUserHandler(service *app.App) *CreateUserHandler {
	return &CreateUserHandler{service: service}
}

func (h *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (any, *middlewares.ServeHTTPError) {
	valErr := validateMethod(http.MethodPost, r.Method)
	if valErr != nil {
		return nil, valErr
	}

	reader := r.Body
	if reader == nil {
		message := "no request body"
		return nil, &middlewares.ServeHTTPError{
			InternalError: errors.New(message),
			Message:       message,
			Code:          http.StatusBadRequest,
		}
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "can't read request body",
			Code:          http.StatusInternalServerError,
		}
	}

	var req CreateUserReq
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "can't unmarshal request body",
			Code:          http.StatusInternalServerError,
		}
	}

	err = req.Validate()
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       err.Error(),
			Code:          http.StatusBadRequest,
		}
	}

	appReq := app.CreateUserReq{
		Nickname: req.Nickname,
	}

	appResp, err := h.service.CreateUser(appReq)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "service error",
			Code:          http.StatusServiceUnavailable,
		}
	}

	return CreateUserResp{UserID: appResp.UserID}, nil
}
