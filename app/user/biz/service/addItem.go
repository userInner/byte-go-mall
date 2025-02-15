package service

import (
	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/common/logging"
	"byte-go-mall/common/tracing"
	"byte-go-mall/constant/config"
	"byte-go-mall/constant/errno"
	cart "byte-go-mall/kitex_gen/cart"
	"byte-go-mall/model"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CartService struct {
	ctx  context.Context
	repo repo.Repository
}

// NewCartService 创建 CartService
func NewCartService(ctx context.Context, cartRepo repo.Repository) *CartService {
	return &CartService{
		ctx:  ctx,
		repo: cartRepo,
	}
}

// AddItem 添加商品到购物车
func (s *CartService) Run(req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	ctx, span := tracing.Tracer.Start(s.ctx, "CartService.AddItem")
	defer span.End()
	logging.SetSpanWithHostname(span)
	log := logging.LogService(config.AppConfig.App.ServiceName)
	resp = &cart.AddItemResp{}

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
		return nil, errno.NewError(errno.UserNotFound, "user not found")
	}

	// 2. 获取用户的购物车
	shoppingCart, err := s.repo.Cart().GetByUserID(ctx, int64(req.UserId))
	if err != nil {
		// 如果购物车不存在，创建一个新的购物车
		if err == gorm.ErrRecordNotFound {
			shoppingCart = &model.ShoppingCart{
				UserID: int64(req.UserId),
			}
			if err = s.repo.Cart().Create(ctx, shoppingCart); err != nil {
				log.WithFields(
					zap.Int64("user_id", int64(req.UserId)),
					zap.Error(err),
				).Error(ctx, "create shopping cart failed")
				return nil, errno.NewError(errno.DatabaseError, "create shopping cart failed")
			}
		} else {
			log.WithFields(
				zap.Int64("user_id", int64(req.UserId)),
				zap.Error(err),
			).Error(ctx, "get shopping cart failed")
			return nil, errno.NewError(errno.DatabaseError, "get shopping cart failed")
		}
	}

	// 3. 检查商品是否已经在购物车中
	exist, err := s.repo.Cart().ExistsByProductID(ctx, int64(shoppingCart.ID), int64(req.Item.ProductId))
	if err != nil {
		log.WithFields(
			zap.Int64("cart_id", int64(shoppingCart.ID)),
			zap.Int64("product_id", int64(req.Item.ProductId)),
			zap.Error(err),
		).Error(ctx, "check product exist in cart failed")
		return nil, errno.NewError(errno.DatabaseError, "check product exist in cart failed")
	}

	if exist {
		// 4. 如果商品已存在，更新数量
		cartItem, err := s.repo.Cart().GetByProductID(ctx, int64(shoppingCart.ID), int64(req.Item.ProductId))
		if err != nil {
			log.WithFields(
				zap.Int64("cart_id", int64(shoppingCart.ID)),
				zap.Int64("product_id", int64(req.Item.ProductId)),
				zap.Error(err),
			).Error(ctx, "get cart item failed")
			return nil, errno.NewError(errno.DatabaseError, "get cart item failed")
		}

		cartItem.Quantity += int32(req.Item.Quantity)
		if err = s.repo.Cart().UpdateCartItem(ctx, cartItem); err != nil {
			log.WithFields(
				zap.Int64("cart_id", int64(shoppingCart.ID)),
				zap.Int64("product_id", int64(req.Item.ProductId)),
				zap.Error(err),
			).Error(ctx, "update cart item failed")
			return nil, errno.NewError(errno.DatabaseError, "update cart item failed")
		}
	} else {
		// 5. 如果商品不存在，创建新的购物车项
		cartItem := &model.CartItem{
			CartID:    int64(shoppingCart.ID),
			ProductID: int64(req.Item.ProductId),
			Quantity:  int32(req.Item.Quantity),
		}
		if err = s.repo.Cart().AddToCart(ctx, cartItem); err != nil {
			log.WithFields(
				zap.Int64("cart_id", int64(shoppingCart.ID)),
				zap.Int64("product_id", int64(req.Item.ProductId)),
				zap.Error(err),
			).Error(ctx, "create cart item failed")
			return nil, errno.NewError(errno.DatabaseError, "create cart item failed")
		}
	}

	// 6. 返回成功响应
	resp.Status = errno.Success
	log.WithFields(
		zap.Int64("user_id", int64(req.UserId)),
		zap.Int64("cart_id", int64(shoppingCart.ID)),
		zap.Int64("product_id", int64(req.Item.ProductId)),
	).Info(ctx, "item added to cart successfully")
	return resp, nil
}
