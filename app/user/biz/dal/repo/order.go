package repo

import (
	"byte-go-mall/model"
	"context"
	"time"

	"gorm.io/gorm"
)

// OrdersRepository 订单仓储
type OrderRepository struct {
	db *gorm.DB
}

// NewOrdersRepository 创建订单仓储
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

// Create 创建订单
func (r *OrderRepository) Create(ctx context.Context, orders *model.Orders) error {
	return r.db.WithContext(ctx).Create(orders).Error
}

// GetByID 根据订单ID查询订单
func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*model.Orders, error) {
	var orders model.Orders
	err := r.db.WithContext(ctx).Preload("OrderItems").Where("id = ?", id).First(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

// GetByOrdersNo 根据订单编号查询订单
func (r *OrderRepository) GetByOrdersNo(ctx context.Context, ordersNo string) (*model.Orders, error) {
	var orders model.Orders
	err := r.db.WithContext(ctx).Preload("OrderItems").Where("orders_no = ?", ordersNo).First(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

// ListByUserID 根据用户ID查询订单列表
func (r *OrderRepository) ListByUserID(ctx context.Context, userID int64) ([]model.Orders, error) {
	var orderss []model.Orders
	err := r.db.WithContext(ctx).Preload("OrderItems").Where("user_id = ?", userID).Find(&orderss).Error
	if err != nil {
		return nil, err
	}
	return orderss, nil
}

// UpdateStatus 更新订单状态
func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&model.Orders{}).Where("id = ?", id).Update("status", status).Error
}

// UpdatePaymentTime 更新订单支付时间
func (r *OrderRepository) UpdatePaymentTime(ctx context.Context, id string, paymentTime time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Orders{}).Where("id = ?", id).Update("payment_time", paymentTime).Error
}

// UpdateCancelTime 更新订单取消时间
func (r *OrderRepository) UpdateCancelTime(ctx context.Context, id string, cancelTime time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Orders{}).Where("id = ?", id).Update("cancel_time", cancelTime).Error
}

// Delete 删除订单
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Orders{}).Error
}

// DeleteOrdersItemsByOrdersID 根据订单ID删除订单项
func (r *OrderRepository) DeleteOrdersItemsByOrdersID(ctx context.Context, ordersID string) error {
	return r.db.WithContext(ctx).Where("orders_id = ?", ordersID).Delete(&model.OrderItem{}).Error
}

// 1. 检查订单是否存在
func (r *OrderRepository) ExistsByOrderNo(ctx context.Context, orderNo string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Orders{}).Where("order_no = ?", orderNo).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
