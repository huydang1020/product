package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/huyshop/header/common"
	pb "github.com/huyshop/header/product"
	"github.com/huyshop/product/utils"
)

func (p *Product) CreateProductType(ctx context.Context, req *pb.ProductType) (*common.Empty, error) {
	if req.GetName() == "" {
		return nil, errors.New(utils.E_invalid_name)
	}
	if req.CategoryId == "" {
		return nil, errors.New(utils.E_invalid_category_id)
	}
	if req.GetDescription() == "" {
		return nil, errors.New(utils.E_invalid_name)
	}
	if req.GetState() == "" {
		req.State = pb.ProductType_approving.String()
	}
	if req.StoreId == "" {
		return nil, errors.New(utils.E_invalid_store_id)
	}
	req.CreatedAt = time.Now().Unix()
	req.Id = utils.MakeProductTypeId()
	for _, pro := range req.GetProducts() {
		if pro.GetName() == "" {
			return nil, errors.New(utils.E_invalid_name)
		}
		if pro.GetId() == "" {
			pro.Id = utils.MakeProductId()
			pro.CreatedAt = req.GetCreatedAt()
			pro.ProductTypeId = req.GetId()
		}
		if pro.GetState() == "" {
			pro.State = req.GetState()
		}
	}
	if err := p.Db.TransCreateProductType(req); err != nil {
		log.Println("CreateProductType error:", err)
		return nil, errors.New(utils.E_can_not_insert)
	}
	return &common.Empty{}, nil
}

func (p *Product) UpdateProductType(ctx context.Context, req *pb.ProductType) (*common.Empty, error) {
	log.Println("UpdateProductType", req)
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	req.UpdatedAt = time.Now().Unix()
	for _, pro := range req.GetProducts() {
		pro.UpdatedAt = req.GetUpdatedAt()
		if req.GetState() != pb.ProductType_active.String() {
			pro.State = req.GetState()
		}
	}
	if err := p.Db.TransUpdateProductType(req); err != nil {
		log.Println("UpdateProductType error:", err)
		return nil, errors.New(utils.E_can_not_update_product_type)
	}
	return &common.Empty{}, nil
}

func (p *Product) UpdateStateProductType(ctx context.Context, req *pb.ProductType) (*common.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	if req.GetState() == "" {
		return nil, errors.New(utils.E_invalid_state)
	}
	if err := p.Db.TransUpdateStateProductType(req); err != nil {
		log.Println("UpdateStateProductType error:", err)
		return nil, errors.New(utils.E_can_not_update_product_type)
	}
	return &common.Empty{}, nil
}
func (p *Product) ListProductType(ctx context.Context, req *pb.ProductTypeRequest) (*pb.ProductTypes, error) {
	log.Println("ListProductType", req)
	productTypes, err := p.Db.ListProductType(req)
	if err != nil {
		return nil, err
	}
	if len(productTypes) == 0 {
		return &pb.ProductTypes{}, nil
	}
	for _, pty := range productTypes {
		listPr, err := p.Db.ListProduct(&pb.ProductRequest{ProductTypeId: pty.Id})
		if err != nil {
			log.Println("ListProductType error:", err)
			return nil, errors.New(utils.E_internal_error)
		}
		if len(listPr) > 0 {
			pty.Products = listPr
		}
	}
	count, err := p.Db.CountProductType(req)
	if err != nil {
		log.Println("CountProductType error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return &pb.ProductTypes{ProductTypes: productTypes, Total: int32(count)}, nil
}

func (p *Product) GetProductType(ctx context.Context, req *pb.ProductTypeRequest) (*pb.ProductType, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	productType, err := p.Db.GetProductType(req.Id)
	if err != nil {
		log.Println("GetProductType error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return productType, nil
}

func (p *Product) DeleteProductType(ctx context.Context, req *pb.ProductType) (*common.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	if err := p.Db.TransDeleteProductType(req.GetId()); err != nil {
		return nil, errors.New(utils.E_can_not_delete_product_type)
	}
	return &common.Empty{}, nil
}
