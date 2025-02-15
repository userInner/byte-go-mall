package model

type ShoppingCart struct {
	Base
	UserID int64 `gorm:"index:idx_user_item"`
}

// TableName 设置表名
func (ShoppingCart) TableName() string {
	return "tb_shopping_cart"
}
