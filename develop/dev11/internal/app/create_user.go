package app

import "fmt"

// CreateUserReq request struct
type CreateUserReq struct {
	Nickname string
}

// CreateUserResp response struct
type CreateUserResp struct {
	UserID int64
}

// CreateUser request
func (a *App) CreateUser(req CreateUserReq) (CreateUserResp, error) {
	id, err := a.repository.CreateUser(req.Nickname)
	if err != nil {
		return CreateUserResp{}, fmt.Errorf("can't create user in repo: %v", err)
	}

	return CreateUserResp{UserID: id}, nil
}
