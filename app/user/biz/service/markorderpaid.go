package service

import (
	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/common/logging"
	"byte-go-mall/common/tracing"
	"byte-go-mall/constant/config"
	"byte-go-mall/constant/errno"
	"byte-go-mall/kitex_gen/order"
	"byte-go-mall/model"
	"context"
	"time"

	"go.uber.org/zap"
)

type MarkOrderPaid struct {
	ctx  context.Context
	repo repo.Repository
}

// NewListOrder 创建 ListOrder 服务
func NewMarkOrderPaid(ctx context.Context, orderRepo repo.Repository) *MarkOrderPaid {
	return &MarkOrderPaid{
		ctx:  ctx,
		repo: orderRepo,
	}
}

// MarkOrderPaid 标记订单为已支付
func (s *MarkOrderPaid) Run(req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	ctx, span := tracing.Tracer.Start(s.ctx, "OrderService.MarkOrderPaid")
	defer span.End()
	logging.SetSpanWithHostname(span)
	log := logging.LogService(config.AppConfig.App.ServiceName)
	resp = &order.MarkOrderPaidResp{}

	// 1. 检查订单是否存在
	orderExist, err := s.repo.Order().ExistsByOrderNo(ctx, req.OrderId)
	if err != nil {
		log.WithFields(
			zap.String("order_id", req.OrderId),
			zap.Error(err),
		).Error(ctx, "check order exist failed")
		return nil, errno.NewError(errno.DatabaseError, "check order exist failed")
	}
	if !orderExist {
		return nil, errno.NewError(errno.OrderNotFound, "order does not exist")
	}

	// 2. 更新订单状态为已支付
	err = s.repo.Order().UpdateStatus(ctx, req.OrderId, model.OrderStatusPaid)
	if err != nil {
		log.WithFields(
			zap.String("order_id", req.OrderId),
			zap.Error(err),
		).Error(ctx, "mark order paid failed")
		return nil, errno.NewError(errno.DatabaseError, "mark order paid failed")
	}

	// 3. 更新支付时间
	paymentTime := time.Now()
	err = s.repo.Order().UpdatePaymentTime(ctx, req.OrderId, paymentTime)
	if err != nil {
		log.WithFields(
			zap.String("order_id", req.OrderId),
			zap.Error(err),
		).Error(ctx, "update payment time failed")
		return nil, errno.NewError(errno.DatabaseError, "update payment time failed")
	}

	resp.Status = errno.Success
	log.WithFields(
		zap.String("order_id", req.OrderId),
	).Info(ctx, "mark order paid successfully")
	return resp, nil
}
