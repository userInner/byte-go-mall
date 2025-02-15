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

type GetCart struct {
	ctx  context.Context
	repo repo.Repository
}

// NewGetCart 创建 GetCart 服务
func NewGetCart(ctx context.Context, cartRepo repo.Repository) *GetCart {
	return &GetCart{
		ctx:  ctx,
		repo: cartRepo,
	}
}

func (s *GetCart) Run(req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	ctx, span := tracing.Tracer.Start(s.ctx, "CartService.GetCart")
	defer span.End()
	logging.SetSpanWithHostname(span)
	log := logging.LogService(config.AppConfig.App.ServiceName)
	resp = &cart.GetCartResp{}

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
			// 如果购物车不存在，返回一个空的购物车
			resp.Cart = &cart.Cart{
				UserId: req.UserId,
				Items:  []*cart.CartItem{},
			}
			log.WithFields(
				zap.Int64("user_id", int64(req.UserId)),
			).Info(ctx, "cart is empty")
			return resp, nil
		}
		log.WithFields(
			zap.Int64("user_id", int64(req.UserId)),
			zap.Error(err),
		).Error(ctx, "get shopping cart failed")
		return nil, errno.NewError(errno.DatabaseError, "get shopping cart failed")
	}

	// 3. 获取购物车项
	cartItems, err := s.repo.Cart().GetCartItemsByCartID(ctx, int64(shoppingCart.ID))
	if err != nil {
		log.WithFields(
			zap.Int64("cart_id", int64(shoppingCart.ID)),
			zap.Error(err),
		).Error(ctx, "get cart items failed")
		return nil, errno.NewError(errno.DatabaseError, "get cart items failed")
	}

	// 4. 构造返回结果
	resp.Cart = &cart.Cart{
		UserId: req.UserId,
		Items:  make([]*cart.CartItem, 0, len(cartItems)),
	}
	for _, item := range cartItems {
		resp.Cart.Items = append(resp.Cart.Items, &cart.CartItem{
			ProductId: uint32(item.ProductID),
			Quantity:  int32(item.Quantity),
		})
	}

	log.WithFields(
		zap.Int64("user_id", int64(req.UserId)),
		zap.Int64("cart_id", int64(shoppingCart.ID)),
	).Info(ctx, "get cart successfully")
	return resp, nil
}
