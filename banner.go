package main

import (
	"context"
	"errors"
	"time"

	"github.com/huyshop/header/common"
	pb "github.com/huyshop/header/product"
	"github.com/huyshop/product/utils"
)

func (p *Product) CreateBanner(ctx context.Context, req *pb.Banner) (*common.Empty, error) {
	if req.Name == "" {
		return nil, errors.New(utils.E_banner_name_empty)
	}
	req.Id = utils.MakeBannerId()
	req.CreatedAt = time.Now().Unix()
	req.State = pb.Banner_active.String()
	if err := p.Db.CreateBanner(req); err != nil {
		return nil, errors.New(utils.E_can_not_insert)
	}
	return &common.Empty{}, nil
}

func (p *Product) UpdateBanner(ctx context.Context, req *pb.Banner) (*common.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_banner_id)
	}
	req.UpdatedAt = time.Now().Unix()
	if err := p.Db.UpdateBanner(req, &pb.Banner{Id: req.GetId()}); err != nil {
		return nil, errors.New(utils.E_can_not_update)
	}
	return &common.Empty{}, nil
}

func (p *Product) ListBanner(ctx context.Context, req *pb.BannerRequest) (*pb.Banners, error) {
	bn, err := p.Db.ListBanner(req)
	if err != nil {
		return nil, errors.New(utils.E_internal_error)
	}
	count, err := p.Db.CountBanner(req)
	if err != nil {
		return nil, errors.New(utils.E_internal_error)
	}
	return &pb.Banners{Banners: bn, Total: int32(count)}, nil
}

func (p *Product) GetBanner(ctx context.Context, req *pb.BannerRequest) (*pb.Banner, error) {
	cate, err := p.Db.GetBanner(req.Id)
	if err != nil {
		return nil, errors.New(utils.E_internal_error)
	}
	return cate, nil
}

func (p *Product) DeleteBanner(ctx context.Context, req *pb.Banner) (*common.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_banner_id)
	}
	if err := p.Db.DeleteBanner(req); err != nil {
		return nil, errors.New(utils.E_internal_error)
	}
	return &common.Empty{}, nil
}
