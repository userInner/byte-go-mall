package repo

import (
	"byte-go-mall/model"
	"context"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository 创建Cart仓储
func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

// Create 创建购物车
func (r *CartRepository) Create(ctx context.Context, shoppingCart *model.ShoppingCart) error {
	return r.db.WithContext(ctx).Create(shoppingCart).Error
}

// GetByID 根据ID获取购物车
func (r *CartRepository) GetByID(ctx context.Context, id int64) (*model.ShoppingCart, error) {
	var shoppingCart model.ShoppingCart
	err := r.db.WithContext(ctx).First(&shoppingCart, id).Error
	if err != nil {
		return nil, err
	}
	return &shoppingCart, nil
}

// GetByUserID 根据用户ID获取购物车
func (r *CartRepository) GetByUserID(ctx context.Context, userID int64) (*model.ShoppingCart, error) {
	var shoppingCart model.ShoppingCart
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&shoppingCart).Error
	if err != nil {
		return nil, err
	}
	return &shoppingCart, nil
}

// AddToCart 向购物车中添加商品
func (r *CartRepository) AddToCart(ctx context.Context, cartItem *model.CartItem) error {
	return r.db.WithContext(ctx).Create(cartItem).Error
}

// UpdateCartItem 更新购物车项
func (r *CartRepository) UpdateCartItem(ctx context.Context, cartItem *model.CartItem) error {
	return r.db.WithContext(ctx).Save(cartItem).Error
}

// ExistsByProductID 检查某个商品是否已经在购物车中
func (r *CartRepository) ExistsByProductID(ctx context.Context, cartID int64, productID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CartItem{}).
		Where("cart_id = ? AND product_id = ?", cartID, productID).Count(&count).Error
	return count > 0, err
}

// GetByProductID 根据商品ID获取购物车项
func (r *CartRepository) GetByProductID(ctx context.Context, cartID, productID int64) (*model.CartItem, error) {
	var cartItem model.CartItem
	err := r.db.WithContext(ctx).Where("cart_id = ? AND product_id = ?", cartID, productID).First(&cartItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 未找到记录，返回 nil
		}
		return nil, err // 其他错误，返回错误信息
	}
	return &cartItem, nil
}

// DeleteCartItemsByCartID 根据购物车ID删除所有购物车项
func (r *CartRepository) DeleteCartItemsByCartID(ctx context.Context, cartID int64) error {
	return r.db.WithContext(ctx).Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error
}

// GetCartItemsByCartID 根据购物车ID获取所有购物车项
func (r *CartRepository) GetCartItemsByCartID(ctx context.Context, cartID int64) ([]model.CartItem, error) {
	var cartItems []model.CartItem
	err := r.db.WithContext(ctx).Where("cart_id = ?", cartID).Find(&cartItems).Error
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

// BatchUpdateCartItems 批量更新购物车项
func (r *CartRepository) BatchUpdateCartItems(ctx context.Context, cartItems []model.CartItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range cartItems {
			if err := tx.Save(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
