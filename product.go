package main

import (
	"context"
	"errors"
	"log"

	pb "github.com/huyshop/header/product"
	"github.com/huyshop/product/utils"
)

func (p *Product) ListProduct(ctx context.Context, req *pb.ProductRequest) (*pb.Products, error) {
	log.Println("ListProduct", req)
	products, err := p.Db.ListProduct(req)
	if err != nil {
		return nil, err
	}
	return &pb.Products{Products: products}, nil
}

func (p *Product) GetProduct(ctx context.Context, req *pb.ProductRequest) (*pb.Product, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	product, err := p.Db.GetProduct(req.Id)
	if err != nil {
		log.Println("GetProduct error:", err)
		return nil, errors.New(utils.E_not_found_product)
	}
	return product, nil
}

func (p *Product) DeleteProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	if err := p.Db.DeleteProduct(req); err != nil {
		return nil, err
	}
	return req, nil
}
