package repo

import "gorm.io/gorm"

// Repository 仓储接口集合
type Repository interface {
	User() *UserRepository
	// Add other repositories...
	Cart() *CartRepository
	Order() *OrderRepository
}

// repository 仓储实现
type repository struct {
	db   *gorm.DB
	user *UserRepository
	// Add other repositories...
	cart  *CartRepository
	order *OrderRepository
}

// NewRepository 创建仓储实例
func NewRepository(db *gorm.DB) Repository {
	repo := &repository{
		db: db,
	}

	// Initialize sub-repositories
	repo.user = NewUserRepository(db)
	repo.cart = NewCartRepository(db)
	repo.order = NewOrderRepository(db)
	return repo
}

func (r *repository) User() *UserRepository {
	return r.user
}

func (r *repository) Cart() *CartRepository {
	return r.cart
}

func (r *repository) Order() *OrderRepository {
	return r.order
}
