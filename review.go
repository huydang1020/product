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
	DEFAULT_LIMIT = 10
)

func (p *Product) CreateReview(ctx context.Context, review *pb.Review) (*common.Empty, error) {
	if review.UserId == "" {
		return nil, errors.New(utils.E_not_found_user_id)
	}
	if review.ProductId == "" {
		return nil, errors.New(utils.E_not_found_product_id)
	}
	if review.OrderId == "" {
		return nil, errors.New(utils.E_not_found_order_id)
	}
	odt, err := p.Db.ListOrderDetail(&pb.OrderDetailRequest{ProductId: review.ProductId, OrderId: review.OrderId})
	if err != nil {
		return nil, err
	}
	if len(odt) == 0 {
		return nil, errors.New(utils.E_not_found_order_detail)
	}
	if review.Rating < 1 || review.Rating > 5 {
		return nil, errors.New(utils.E_invalid_rating)
	}
	if err := p.Db.CreateReview(review); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (p *Product) UpdateReview(ctx context.Context, req *pb.Review) (*common.Empty, error) {
	if req.Id == "" {
		return nil, errors.New(utils.E_not_found_review_id)
	}
	if err := p.Db.UpdateReview(req, &pb.Review{Id: req.Id}); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (p *Product) DeleteReview(ctx context.Context, req *pb.Review) (*common.Empty, error) {
	if req.Id == "" {
		return nil, errors.New(utils.E_not_found_review_id)
	}
	if err := p.Db.DeleteReview(req); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (p *Product) ListReview(ctx context.Context, rq *pb.ReviewRequest) (*pb.Reviews, error) {
	log.Println("ListReview", rq)
	if rq.Limit > 100 {
		rq.Limit = DEFAULT_LIMIT
	}
	reviews, err := p.Db.ListReview(rq)
	if err != nil {
		return nil, err
	}
	return &pb.Reviews{Reviews: reviews, Total: int32(len(reviews))}, nil
}

func (p *Product) GetReview(ctx context.Context, req *pb.ReviewRequest) (*pb.Review, error) {
	if req.Id == "" {
		return nil, errors.New(utils.E_not_found_review_id)
	}
	review, err := p.Db.GetReview(req.Id)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (p *Product) CountReview(ctx context.Context, rq *pb.ReviewRequest) (*common.Count, error) {
	count, err := p.Db.CountReview(rq)
	if err != nil {
		return nil, err
	}
	return &common.Count{Count: count}, nil
}

func (p *Product) ReplyReview(ctx context.Context, req *pb.Review) (*common.Empty, error) {
	if req.Id == "" {
		return nil, errors.New(utils.E_not_found_review_id)
	}
	if req.SellerReply == "" {
		return nil, errors.New(utils.E_not_found_review_reply)
	}
	req.SellerReplyAt = time.Now().Unix()
	if err := p.Db.UpdateReview(req, &pb.Review{Id: req.Id}); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
