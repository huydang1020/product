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
	if check := p.Db.IsReviewExist(&pb.Review{UserId: review.UserId, OrderId: review.OrderId, ProductId: review.ProductId}); check {
		return nil, errors.New(utils.E_review_already_exists)
	}
	ord, err := p.Db.GetOrder(review.OrderId)
	if err != nil {
		return nil, err
	}
	if ord == nil || ord.UserId != review.UserId {
		return nil, errors.New(utils.E_invalid_order)
	}

	// Kiểm tra sản phẩm có nằm trong đơn hàng không
	found := false
	for _, odt := range ord.ProductOrdered {
		if odt.ProductId == review.ProductId {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New(utils.E_product_not_in_order)
	}

	// Kiểm tra rating hợp lệ
	if review.Rating < 1 || review.Rating > 5 {
		return nil, errors.New(utils.E_invalid_rating)
	}

	// Tạo review
	review.Id = utils.MakeReviewId()
	review.CreatedAt = time.Now().Unix()
	if err := p.Db.CreateReview(review); err != nil {
		return nil, err
	}

	return &common.Empty{}, nil
}

func (p *Product) UpdateReview(ctx context.Context, req *pb.Review) (*common.Empty, error) {
	if req.Id == "" {
		return nil, errors.New(utils.E_not_found_review_id)
	}
	req.UpdatedAt = time.Now().Unix()
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
	log.Println("limit: ", len(rq.ProductIds))
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
	review, err := p.Db.GetReview(&pb.Review{Id: req.Id})
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
	req.UpdatedAt = time.Now().Unix()
	if err := p.Db.UpdateReview(req, &pb.Review{Id: req.Id}); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
