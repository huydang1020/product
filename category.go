package main

import (
	"context"
	"errors"
	"time"

	"github.com/huyshop/header/common"
	pb "github.com/huyshop/header/product"
	"github.com/huyshop/product/utils"
)

func (p *Product) CreateCategory(ctx context.Context, req *pb.Category) (*common.Empty, error) {
	if req.Name == "" {
		return nil, errors.New(utils.E_category_name_empty)
	}
	if req.GetLogo() == "" {
		return nil, errors.New(utils.E_category_logo_empty)
	}
	req.Id = utils.MakeCategoryId()
	req.CreatedAt = time.Now().Unix()
	if err := p.Db.CreateCategory(req); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (p *Product) UpdateCategory(ctx context.Context, req *pb.Category) (*common.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_category_id)
	}
	req.UpdatedAt = time.Now().Unix()
	if err := p.Db.UpdateCategory(req, &pb.Category{Id: req.GetId()}); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (p *Product) ListCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.Categories, error) {
	cates, err := p.Db.ListCategory(req)
	if err != nil {
		return nil, err
	}
	return &pb.Categories{Categories: cates}, nil
}

func (p *Product) GetCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.Category, error) {
	cate, err := p.Db.GetCategory(req.Id)
	if err != nil {
		return nil, err
	}
	return cate, nil
}

func (p *Product) DeleteCategory(ctx context.Context, req *pb.Category) (*common.Empty, error) {
	if err := p.Db.DeleteCategory(req); err != nil {
		return nil, err
	}
	if err := p.cache.Del(ctx, req.Id).Err(); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
