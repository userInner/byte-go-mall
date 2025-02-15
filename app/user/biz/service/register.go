package service

import (
	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/common/logging"
	"byte-go-mall/common/tracing"
	"byte-go-mall/constant/config"
	"byte-go-mall/constant/errno"
	"byte-go-mall/kitex_gen/base"
	user "byte-go-mall/kitex_gen/user"
	"byte-go-mall/model"
	"byte-go-mall/utils"
	"context"
	"errors"

	"go.uber.org/zap"
)

type RegisterService struct {
	ctx  context.Context
	repo repo.Repository
} // NewRegisterService new RegisterService
func NewRegisterService(ctx context.Context, userRepo repo.Repository) *RegisterService {
	return &RegisterService{
		ctx:  ctx,
		repo: userRepo,
	}
}

// Run create note info
func (s *RegisterService) Run(req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// Finish your business logic.
	ctx, span := tracing.Tracer.Start(s.ctx, "Register.Run")
	defer span.End()
	logging.SetSpanWithHostname(span)
	log := logging.LogService(config.AppConfig.App.ServiceName)
	resp = &user.RegisterResp{}

	// 参数验证
	if err = s.validateRegister(req); err != nil {
		resp.Status = errno.WithMessage(&base.Status{Code: 10003}, err.Error())

		return resp, err
	}
	// 判断密码是否严格
	if !utils.IsValidPassword(req.Password) {
		resp.Status = errno.WeakPassword
		return resp, nil
	}

	// 检查用户是否存在
	exist, err := s.repo.User().ExistsByEmail(ctx, req.Email)
	if err != nil {
		resp.Status = errno.DatabaseError
		log.WithFields(
			zap.String("email", req.Email),
			zap.Error(err),
		).Info(ctx, "check email exist failed")
		return resp, err
	}
	if exist {
		resp.Status = errno.UserExist
		return resp, nil
	}

	// 创建用户
	u := &model.User{
		Username: utils.GenerateRandomString(7),
		Email:    req.Email,
	}
	// 密码加密
	if err = u.SetPassword(req.Password); err != nil {
		log.WithFields(
			zap.Error(err),
			zap.String("password", req.Password),
		).Error(ctx, "set password failed")
		resp.Status = errno.InternalError
		return resp, err
	}

	// 保存用户
	if err = s.repo.User().Create(ctx, u); err != nil {
		log.WithFields(
			zap.Error(err),
		).Error(ctx, "create user failed")
		resp.Status = errno.DatabaseError
		return resp, err
	}
	resp.Status = errno.Success
	resp.UserId = int32(u.ID)
	log.WithFields(
		zap.Int64("user_id", int64(u.ID)),
		zap.String("email", u.Email),
	).Info(ctx, "user registered successfully")
	return resp, nil
}

// validateRegister 验证注册参数
func (s *RegisterService) validateRegister(req *user.RegisterReq) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if !utils.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if req.Password != req.ConfirmPassword {
		return errors.New("password not match")
	}
	if len(req.Password) < 6 || len(req.Password) > 32 {
		return errors.New("password length should be between 6 and 32")
	}
	return nil
}
