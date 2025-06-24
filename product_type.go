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

const (
	DESC = "desc"
	ASC  = "asc"
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
		req.State = pb.ProductType_pending.String()
	}
	if req.StoreId == "" {
		return nil, errors.New(utils.E_invalid_store_id)
	}
	req.CreatedAt = time.Now().Unix()
	req.Id = utils.MakeProductTypeId()
	req.Slug = utils.ToSlug(req.GetName())
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
	log.Println("ListProductType req", req)
	listCate, err := p.Db.ListCategory(&pb.CategoryRequest{Id: req.CategoryId, Slug: req.Category})
	if err != nil {
		log.Println("err: ", err)
		return nil, err
	}
	mapCate := make(map[string]*pb.Category)
	for _, c := range listCate {
		mapCate[c.Id] = c
	}
	if req.Category != "" {
		if len(listCate) > 0 {
			req.CategoryId = listCate[0].Id
		} else {
			req.CategoryId = "not_category"
		}
	}
	productTypes, err := p.Db.ListProductType(req)
	if err != nil {
		log.Println("err", err)
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
		pty.Category = mapCate[pty.GetCategoryId()]
		listReview, err := p.Db.ListReview(&pb.ReviewRequest{ProductTypeId: pty.Id})
		if err != nil {
			log.Println("list review err: ", err)
		}
		pty.TotalReviews = int32(len(listReview))
		rate := p.CaculateAvgrating(listReview)
		pty.AverageRating = rate
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
	pty, err := p.Db.GetProductType(req.Id)
	if err != nil {
		log.Println("GetProductType error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	listProduct, err := p.Db.ListProduct(&pb.ProductRequest{ProductTypeId: pty.Id})
	if err != nil {
		log.Println("ListProductType error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	pty.Products = listProduct
	listReview, err := p.Db.ListReview(&pb.ReviewRequest{ProductTypeId: pty.Id})
	if err != nil {
		log.Println("list review err: ", err)
		return nil, errors.New(utils.E_internal_error)
	}
	pty.Reviews = listReview
	pty.TotalReviews = int32(len(listReview))
	rate := p.CaculateAvgrating(listReview)
	pty.AverageRating = rate
	return pty, nil
}

func (p *Product) GetProductTypeBySlug(ctx context.Context, req *pb.ProductTypeRequest) (*pb.ProductType, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	pty, err := p.Db.GetProductTypeBySlug(req.Id)
	if err != nil {
		log.Println("GetProductType error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	listProduct, err := p.Db.ListProduct(&pb.ProductRequest{ProductTypeId: pty.Id})
	if err != nil {
		log.Println("ListProductType error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	pty.Products = listProduct
	listReview, err := p.Db.ListReview(&pb.ReviewRequest{ProductTypeId: pty.Id})
	if err != nil {
		log.Println("list review err: ", err)
		return nil, errors.New(utils.E_internal_error)
	}
	pty.Reviews = listReview
	pty.TotalReviews = int32(len(listReview))
	rate := p.CaculateAvgrating(listReview)
	pty.AverageRating = rate
	return pty, nil
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

func (p Product) CountProductType(ctx context.Context, req *pb.ProductTypeRequest) (*common.Count, error) {
	count, err := p.Db.CountProductType(req)
	if err != nil {
		log.Println("CountProductType error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return &common.Count{Count: count}, nil
}
