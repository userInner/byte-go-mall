package model

import (
	"time"
)

// OrderStatus 定义订单状态常量
type OrderStatus int8

const (
	OrderStatusPending   OrderStatus = 0 // 待支付
	OrderStatusPaid      OrderStatus = 1 // 已支付
	OrderStatusCancelled OrderStatus = 2 // 已取消
	OrderStatusCompleted OrderStatus = 3 // 已完成
)

// Address 表示订单的地址信息
type Address struct {
	StreetAddress string `gorm:"street_address"` // 街道地址
	City          string `gorm:"city"`           // 城市
	State         string `gorm:"state"`          // 州/省
	Country       string `gorm:"country"`        // 国家
	ZipCode       int32  `gorm:"zip_code"`       // 邮政编码
}

// Order 表示订单表
type Orders struct {
	Base
	OrderNo     string      `gorm:"order_no"`           // 订单编号
	UserID      int64       `gorm:"user_id"`            // 用户ID
	TotalAmount float64     `gorm:"type:decimal(10,2)"` // 订单总金额
	Status      OrderStatus `gorm:"status"`             // 订单状态
	PaymentTime *time.Time  `gorm:"payment_time"`       // 支付时间
	CancelTime  *time.Time  `gorm:"cancel_time"`        // 取消时间
	ExpireTime  *time.Time  `gorm:"expire_time"`        // 过期时间
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"` // 订单项列表
	Address     Address     `gorm:"embedded"`           // 订单地址信息
	Email       string      `gorm:"email"`              // 订单的电子邮件地址
	CreatedAt   time.Time   `gorm:"created_at"`         // 订单创建时间
}

// TableName 返回订单表名
func (o *Orders) TableName() string {
	return "tb_order"
}
