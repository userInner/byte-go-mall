package main

import (
	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/app/user/biz/service"
	user "byte-go-mall/kitex_gen/user"
	"context"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	repo repo.Repository // 用户CRUD
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	resp, err = service.NewRegisterService(ctx, s.repo).Run(req)

	return resp, err
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	resp, err = service.NewLoginService(ctx, s.repo).Run(req)

	return resp, err
}
