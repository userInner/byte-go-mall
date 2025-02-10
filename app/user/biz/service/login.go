package service

import (
	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/common/logging"
	"byte-go-mall/common/tracing"
	"byte-go-mall/constant/config"
	"byte-go-mall/constant/errno"
	user "byte-go-mall/kitex_gen/user"
	"byte-go-mall/utils"
	"context"
	"errors"
	"go.uber.org/zap"
)

type LoginService struct {
	ctx  context.Context
	repo repo.Repository
} // NewLoginService new LoginService
func NewLoginService(ctx context.Context, repo repo.Repository) *LoginService {
	return &LoginService{
		ctx:  ctx,
		repo: repo,
	}
}

// Run handles user login
func (s *LoginService) Run(req *user.LoginReq) (resp *user.LoginResp, err error) {
	ctx, span := tracing.Tracer.Start(s.ctx, "Login.Run")
	defer span.End()
	logging.SetSpanWithHostname(span)
	log := logging.LogService(config.AppConfig.App.ServiceName)

	resp = &user.LoginResp{}

	// 参数验证
	if err = s.validateLogin(req); err != nil {
		resp.Status = errno.InvalidParams
		return resp, err
	}

	// 获取用户信息
	u, err := s.repo.User().GetByEmail(ctx, req.Email)
	if err != nil {
		log.WithFields(
			zap.String("email", req.Email),
			zap.Error(err),
		).Error(ctx, "get user by email failed")
		resp.Status = errno.DatabaseError
		return resp, err
	}

	// 用户不存在
	if u == nil {
		resp.Status = errno.UserNotFound
		return resp, nil
	}

	// 验证密码
	if !u.CheckPassword(req.Password) {
		log.WithFields(
			zap.String("email", req.Email),
		).Info(ctx, "password incorrect")
		resp.Status = errno.PasswordError
		return resp, nil
	}

	// 设置响应
	resp.Status = errno.Success
	resp.UserId = int32(u.ID)

	log.WithFields(
		zap.Int64("user_id", int64(u.ID)),
		zap.String("email", u.Email),
	).Info(ctx, "user login successfully")

	return resp, nil
}

// validateLogin 验证登录参数
func (s *LoginService) validateLogin(req *user.LoginReq) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if !utils.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
