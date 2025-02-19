package service

import (
	"context"
	"time"

	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/common/logging"
	"byte-go-mall/common/tracing"
	"byte-go-mall/constant/config"
	"byte-go-mall/constant/errno"
	"byte-go-mall/kitex_gen/order"
	"byte-go-mall/model"

	"go.uber.org/zap"
)

type PlaceOrderService struct {
	ctx  context.Context
	repo repo.Repository
}

// NewPlaceOrderService 创建 PlaceOrderService 服务
func NewPlaceOrderService(ctx context.Context, orderRepo repo.Repository) *PlaceOrderService {
	return &PlaceOrderService{
		ctx:  ctx,
		repo: orderRepo,
	}
}

// Run 启动订单创建服务
func (s *PlaceOrderService) Run(req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	// 启动 tracing span
	ctx, span := tracing.Tracer.Start(s.ctx, "PlaceOrderService.Run")
	defer span.End()
	logging.SetSpanWithHostname(span)
	log := logging.LogService(config.AppConfig.App.ServiceName)
	resp = &order.PlaceOrderResp{}

	// 1. 检查用户是否存在
	userExist, err := s.repo.User().ExistsByID(ctx, int64(req.UserId))
	if err != nil {
		log.With(
			zap.Int64("user_id", int64(req.UserId)),
			zap.Error(err),
		).Error("check user exist failed")
		return nil, errno.NewError(errno.DatabaseError, "check user exist failed")
	}
	if !userExist {
		return nil, errno.NewError(errno.UserNotFound, "user does not exist")
	}

	// 2. 创建订单
	orderList := &model.Orders{
		OrderNo:     generateOrderNo(), // 生成订单号
		UserID:      int64(req.UserId),
		TotalAmount: calculateTotalAmount(req.OrderItems), // 计算总金额
		Status:      model.OrderStatusPending,
		ExpireTime:  func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(), // 设置过期时间
		OrderItems:  convertOrderItems(req.OrderItems),                                      // 转换订单项
	}

	// 保存订单到数据库
	err = s.repo.Order().Create(ctx, orderList)
	if err != nil {
		log.With(
			zap.Int64("user_id", int64(req.UserId)),
			zap.Error(err),
		).Error("create order failed")
		return nil, errno.NewError(errno.DatabaseError, "create order failed")
	}

	// 3. 返回订单结果
	resp.Order = &order.OrderResult{
		OrderId: orderList.OrderNo,
	}
	resp.Status = errno.Success

	log.With(
		zap.Int64("user_id", int64(req.UserId)),
		zap.String("order_no", orderList.OrderNo),
	).Info("place order successfully")
	return resp, nil
}

// 生成订单号
func generateOrderNo() string {
	return time.Now().Format("20060102150405") // 示例：20231010123045
}

// 计算订单总金额
func calculateTotalAmount(items []*order.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Cost)
	}
	return total
}

// 转换订单项
func convertOrderItems(items []*order.OrderItem) []model.OrderItem {
	var orderItems []model.OrderItem
	for _, item := range items {
		orderItems = append(orderItems, model.OrderItem{
			ProductID:    int64(item.Item.ProductId),
			ProductName:  item.Item.ProductName,
			ProductPrice: float32(item.Cost),
			Quantity:     int(item.Item.Quantity),
		})
	}
	return orderItems
}
