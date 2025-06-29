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

func (p *Product) CreateReviews(ctx context.Context, review *pb.Reviews) (*common.Empty, error) {
	if review.UserId == "" {
		return nil, errors.New(utils.E_not_found_user_id)
	}
	if review.ProductId == "" {
		return nil, errors.New(utils.E_not_found_product_id)
	}
	if review.OrderId == "" {
		return nil, errors.New(utils.E_not_found_order_id)
	}
	if check := p.Db.IsReviewsExist(&pb.Reviews{UserId: review.UserId, OrderId: review.OrderId, ProductId: review.ProductId}); check {
		return nil, errors.New(utils.E_reviews_already_exists)
	}
	ord, err := p.Db.GetOrder(review.OrderId)
	if err != nil {
		return nil, err
	}
	if ord == nil || ord.UserId != review.UserId {
		return nil, errors.New(utils.E_invalid_order)
	}
	if ord.State != pb.Order_completed.String() {
		return nil, errors.New(utils.E_invalid_state_order)
	}
	pro, err := p.Db.GetProduct(review.ProductId)
	if err != nil {
		log.Println("err ", err)
		return nil, errors.New(utils.E_not_found_product)
	}
	review.ProductTypeId = pro.ProductTypeId

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
	if err := p.Db.CreateReviews(review); err != nil {
		log.Println("err: ", err)
		return nil, errors.New(utils.E_internal_error)
	}
	// update thông tin product
	pty := &pb.ProductType{}
	listReview, err := p.Db.ListReviews(&pb.ReviewsRequest{ProductTypeId: pro.ProductTypeId})
	if err != nil {
		log.Println("list review err: ", err)
		return nil, errors.New(utils.E_internal_error)
	}
	pty.TotalReviews = int32(len(listReview))
	rate := p.CaculateAvgrating(listReview)
	pty.AverageRating = rate
	if err = p.Db.UpdateProductType(pty, &pb.ProductType{Id: pro.ProductTypeId}); err != nil {
		log.Println("update review by productType err: ", err)
	}
	return &common.Empty{}, nil
}

func (p *Product) UpdateReviews(ctx context.Context, req *pb.Reviews) (*common.Empty, error) {
	if req.Id == "" {
		return nil, errors.New(utils.E_not_found_reviews_id)
	}
	req.UpdatedAt = time.Now().Unix()
	if req.SellerReply != "" {
		req.SellerReplyAt = time.Now().Unix()
	}
	if err := p.Db.UpdateReviews(req, &pb.Reviews{Id: req.Id}); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (p *Product) DeleteReviews(ctx context.Context, req *pb.Reviews) (*common.Empty, error) {
	if req.Id == "" {
		return nil, errors.New(utils.E_not_found_reviews_id)
	}
	if err := p.Db.DeleteReviews(req); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (p *Product) CaculateAvgrating(reviews []*pb.Reviews) float32 {
	var totalRating int32
	for _, ra := range reviews {
		totalRating += ra.Rating
	}

	var avgRating float32
	if len(reviews) > 0 {
		avgRating = float32(totalRating) / float32(len(reviews))
	}

	return avgRating

}

func (p *Product) ListReviews(ctx context.Context, rq *pb.ReviewsRequest) (*pb.Reviewss, error) {
	log.Println("ListReview", rq)
	reviews, err := p.Db.ListReviews(rq)
	if err != nil {
		return nil, err
	}
	for _, rv := range reviews {
		prod, err := p.Db.GetProduct(rv.GetProductId())
		if err != nil {
			continue
		}
		rv.Product = prod
	}
	count, err := p.Db.CountReviews(rq)
	if err != nil {
		return nil, err
	}
	return &pb.Reviewss{Reviews: reviews, Total: int32(count)}, nil
}

func (p *Product) GetReviews(ctx context.Context, req *pb.ReviewsRequest) (*pb.Reviews, error) {
	if req.Id == "" {
		return nil, errors.New(utils.E_not_found_reviews_id)
	}
	review, err := p.Db.GetReviews(&pb.Reviews{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (p *Product) CountReviews(ctx context.Context, rq *pb.ReviewsRequest) (*common.Count, error) {
	count, err := p.Db.CountReviews(rq)
	if err != nil {
		return nil, err
	}
	return &common.Count{Count: count}, nil
}

func (p *Product) ReplyReviews(ctx context.Context, req *pb.Reviews) (*common.Empty, error) {
	if req.Id == "" {
		return nil, errors.New(utils.E_not_found_reviews_id)
	}
	if req.SellerReply == "" {
		return nil, errors.New(utils.E_not_found_reviews_reply)
	}
	req.SellerReplyAt = time.Now().Unix()
	req.UpdatedAt = time.Now().Unix()
	if err := p.Db.UpdateReviews(req, &pb.Reviews{Id: req.Id}); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
