package model

// OrderItem 表示订单项表
type OrderItem struct {
	Base
	OrderID      int64   `gorm:"column:order_id"`     // 订单ID
	ProductID    int64   `gorm:"column:product_id"`   // 商品ID
	ProductName  string  `gorm:"column:product_name"` // 商品名称
	ProductPrice float32 `gorm:"type:decimal(10,2)"`  // 商品单价
	Quantity     int     `gorm:"column:quantity"`     // 购买数量
}

// TableName 返回订单项表名
func (oi *OrderItem) TableName() string {
	return "tb_order_item"
}
