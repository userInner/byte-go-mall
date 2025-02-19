package main

import (
	"byte-go-mall/app/user/biz/dal/repo"
	"byte-go-mall/app/user/biz/service"
	"byte-go-mall/kitex_gen/cart"
	"byte-go-mall/kitex_gen/order"
	user "byte-go-mall/kitex_gen/user"
	"context"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	repo repo.Repository // 用户CRUD
}

type CartServiceImpl struct {
	repo repo.Repository // 用户CRUD
}

type OrderServiceImpl struct {
	repo repo.Repository
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	resp, err = service.NewRegisterService(ctx, s.repo).Run(req)

	return resp, err
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	resp, err = service.NewLoginService(ctx, s.repo).Run(req)

	return resp, err
}

func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	resp, err = service.NewGetCart(ctx, s.repo).Run(req)

	return resp, err
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	resp, err = service.NewEmptyCart(ctx, s.repo).Run(req)

	return resp, err
}

func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	resp, err = service.NewCartService(ctx, s.repo).Run(req)

	return resp, err
}

func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	resp, err = service.NewListOrder(ctx, s.repo).Run(req)

	return resp, err
}

func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	resp, err = service.NewPlaceOrderService(ctx, s.repo).Run(req)

	return resp, err
}

func (s *OrderServiceImpl) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	resp, err = service.NewMarkOrderPaid(ctx, s.repo).Run(req)

	return resp, err
}
