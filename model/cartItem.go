package model

type CartItem struct {
	Base
	CartID    int64 `gorm:"index:idx_user_item"` // 购物车ID，JSON字段名为 cart_id
	ProductID int64 `gorm:"index:idx_user_item"` // 商品ID，JSON字段名为 product_id
	Quantity  int32 `gorm:"not null;default:1"`  // 商品数量，JSON字段名为 quantity
}

// TableName 设置表名
func (CartItem) TableName() string {
	return "tb_cart_item"
}
