package main

import (
	"context"
	"errors"
	"log"

	"github.com/huyshop/header/common"
	pb "github.com/huyshop/header/product"
	"github.com/huyshop/product/utils"
)

func (p *Product) CreateOrder(ctx context.Context, req *pb.Order) (*common.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	if err := p.Db.CreateOrder(req); err != nil {
		log.Println("CreateOrder error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return &common.Empty{}, nil
}

func (p *Product) UpdateOrder(ctx context.Context, req *pb.Order) (*common.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	if err := p.Db.UpdateOrder(req, &pb.Order{Id: req.GetId()}); err != nil {
		log.Println("UpdateOrder error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return &common.Empty{}, nil
}

func (p *Product) DeleteOrder(ctx context.Context, req *pb.Order) (*common.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	if err := p.Db.DeleteOrder(req); err != nil {
		log.Println("DeleteOrder error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return &common.Empty{}, nil
}

func (p *Product) ListOrder(ctx context.Context, req *pb.OrderRequest) (*pb.Orders, error) {

	if req.GetLimit() < 1 {
		req.Limit = 10
	}
	list, err := p.Db.ListOrder(req)
	if err != nil {
		log.Println("ListOrder error:", err)
		return nil, errors.New(utils.E_internal_error)
	}

	count, err := p.Db.CountOrder(req)
	if err != nil {
		log.Println("CountOrder error:", err)
		return nil, errors.New(utils.E_internal_error)
	}

	return &pb.Orders{Orders: list, Total: int32(count)}, nil
}

func (p *Product) GetOrder(ctx context.Context, req *pb.OrderRequest) (*pb.Order, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	order, err := p.Db.GetOrder(req.GetId())
	if err != nil {
		log.Println("GetOrder error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return order, nil
}

func (p *Product) AddToCart(ctx context.Context, req *pb.Cart) (*pb.Cart, error) {
	if req.GetUserId() == "" {
		return nil, errors.New(utils.E_not_found_user_id)
	}
	if len(req.Item) < 1 {
		return nil, errors.New(utils.E_not_found_item_cart)
	}
	for _, item := range req.Item {
		if item.GetProductId() == "" {
			return nil, errors.New(utils.E_not_found_product)
		}
		
		if item.GetQuantity() < 1 {
			return nil, errors.New(utils.E_inventory_quantity_not_enough)
		}
		
	}
	return nil, nil
}

func (p *Product) UpdateCart(ctx context.Context, req *pb.Cart) (*pb.Cart, error) {
	return nil, nil
}
func (p *Product) DeleteCart(ctx context.Context, req *pb.Cart) (*pb.Cart, error) {
	return nil, nil
}
func (p *Product) ListCart(ctx context.Context, req *pb.Cart) (*pb.Cart, error) {
	return nil, nil
}

func (p *Product) GetCart(ctx context.Context, req *pb.Cart) (*pb.Cart, error) {
	return nil, nil
}
func (p *Product) DeleteCartItem(ctx context.Context, req *pb.Cart) (*pb.Cart, error) {
	return nil, nil
}
