package service

import (
	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/common/logging"
	"byte-go-mall/common/tracing"
	"byte-go-mall/constant/config"
	"byte-go-mall/constant/errno"
	cart "byte-go-mall/kitex_gen/cart"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EmptyCart struct {
	ctx  context.Context
	repo repo.Repository
}

// NewEmptyCart 创建 EmptyCart 服务
func NewEmptyCart(ctx context.Context, cartRepo repo.Repository) *EmptyCart {
	return &EmptyCart{
		ctx:  ctx,
		repo: cartRepo,
	}
}

// Run 清空用户的购物车
func (s *EmptyCart) Run(req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	ctx, span := tracing.Tracer.Start(s.ctx, "CartService.EmptyCart")
	defer span.End()
	logging.SetSpanWithHostname(span)
	log := logging.LogService(config.AppConfig.App.ServiceName)
	resp = &cart.EmptyCartResp{}

	// 1. 检查用户是否存在
	userExist, err := s.repo.User().ExistsByID(ctx, int64(req.UserId))
	if err != nil {
		log.WithFields(
			zap.Int64("user_id", int64(req.UserId)),
			zap.Error(err),
		).Error(ctx, "check user exist failed")
		return nil, errno.NewError(errno.DatabaseError, "check user exist failed")
	}
	if !userExist {
		return nil, errno.NewError(errno.UserNotFound, "user does not exist")
	}

	// 2. 获取用户的购物车
	shoppingCart, err := s.repo.Cart().GetByUserID(ctx, int64(req.UserId))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果购物车不存在，直接返回成功（购物车已经是空的）
			log.WithFields(
				zap.Int64("user_id", int64(req.UserId)),
			).Info(ctx, "cart is already empty")
			resp.Status = errno.Success
			return resp, nil
		}
		log.WithFields(
			zap.Int64("user_id", int64(req.UserId)),
			zap.Error(err),
		).Error(ctx, "get shopping cart failed")
		return nil, errno.NewError(errno.DatabaseError, "get shopping cart failed")
	}

	// 3. 清空购物车项
	if err = s.repo.Cart().DeleteCartItemsByCartID(ctx, int64(shoppingCart.ID)); err != nil {
		log.WithFields(
			zap.Int64("cart_id", int64(shoppingCart.ID)),
			zap.Error(err),
		).Error(ctx, "empty cart failed")
		return nil, errno.NewError(errno.DatabaseError, "empty cart failed")
	}

	log.WithFields(
		zap.Int64("user_id", int64(req.UserId)),
		zap.Int64("cart_id", int64(shoppingCart.ID)),
	).Info(ctx, "cart emptied successfully")
	resp.Status = errno.Success
	return resp, nil
}
