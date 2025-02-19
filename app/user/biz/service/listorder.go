package service

import (
	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/common/logging"
	"byte-go-mall/common/tracing"
	"byte-go-mall/constant/config"
	"byte-go-mall/constant/errno"
	cart "byte-go-mall/kitex_gen/cart"
	"byte-go-mall/kitex_gen/order"
	"byte-go-mall/model"
	"context"

	"go.uber.org/zap"
)

type ListOrder struct {
	ctx  context.Context
	repo repo.Repository
}

// NewListOrder 创建 ListOrder 服务
func NewListOrder(ctx context.Context, orderRepo repo.Repository) *ListOrder {
	return &ListOrder{
		ctx:  ctx,
		repo: orderRepo,
	}
}

func (s *ListOrder) Run(req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	ctx, span := tracing.Tracer.Start(s.ctx, "OrderService.ListOrder")
	defer span.End()
	logging.SetSpanWithHostname(span)
	log := logging.LogService(config.AppConfig.App.ServiceName)

	resp = &order.ListOrderResp{}

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

	// 2. 查询用户订单列表
	orders, err := s.repo.Order().ListByUserID(ctx, int64(req.UserId))
	if err != nil {
		log.WithFields(
			zap.Int64("user_id", int64(req.UserId)),
			zap.Error(err),
		).Error(ctx, "list orders failed")
		return nil, errno.NewError(errno.DatabaseError, "list orders failed")
	}

	// 3. 构造返回结果
	resp.Orders = make([]*order.Order, 0, len(orders))
	for _, o := range orders {
		resp.Orders = append(resp.Orders, &order.Order{
			OrderId:      o.OrderNo,
			UserId:       uint32(o.UserID),
			UserCurrency: o.OrderNo,
			Address: &order.Address{
				StreetAddress: o.Address.StreetAddress,
				City:          o.Address.City,
				State:         o.Address.State,
				Country:       o.Address.Country,
				ZipCode:       o.Address.ZipCode,
			},
			Email:      o.Email,
			OrderItems: convertToProtoOrderItems(o.OrderItems),
			CreatedAt:  int32(o.CreatedAt.Unix()),
		})
	}

	resp.Status = errno.Success
	log.WithFields(
		zap.Int64("user_id", int64(req.UserId)),
	).Info(ctx, "list orders successfully")

	return resp, nil
}

// 转换订单项为 proto 格式
func convertToProtoOrderItems(items []model.OrderItem) []*order.OrderItem {
	var protoItems []*order.OrderItem
	for _, item := range items {
		protoItems = append(protoItems, &order.OrderItem{
			Item: &cart.CartItem{
				ProductId:   uint32(item.ProductID),
				ProductName: item.ProductName,
				Quantity:    int32(item.Quantity),
			},
			Cost: item.ProductPrice,
		})
	}
	return protoItems
}
