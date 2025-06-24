package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	pb "github.com/huyshop/header/product"
	"github.com/huyshop/product/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Product struct {
	Db    IDatabase
	cache *redis.Client
}

type IDatabase interface {
	CreateCategory(category *pb.Category) error
	UpdateCategory(updator, selector *pb.Category) error
	ListCategory(rq *pb.CategoryRequest) ([]*pb.Category, error)
	GetCategory(id string) (*pb.Category, error)
	DeleteCategory(category *pb.Category) error
	CountCategory(rq *pb.CategoryRequest) (int64, error)

	CreateProduct(product *pb.Product) error
	UpdateProduct(updator, selector *pb.Product) error
	DeleteProduct(product *pb.Product) error
	ListProduct(rq *pb.ProductRequest) ([]*pb.Product, error)
	GetProduct(id string) (*pb.Product, error)
	CountProduct(rq *pb.ProductRequest) (int64, error)

	CreateProductType(productType *pb.ProductType) error
	UpdateProductType(updator, selector *pb.ProductType) error
	DeleteProductType(id string) error
	ListProductType(rq *pb.ProductTypeRequest) ([]*pb.ProductType, error)
	GetProductType(id string) (*pb.ProductType, error)
	CountProductType(rq *pb.ProductTypeRequest) (int64, error)
	GetProductTypeBySlug(key string) (*pb.ProductType, error)

	TransCreateProductType(pt *pb.ProductType) error
	TransUpdateStateProductType(pt *pb.ProductType) error
	TransUpdateProductType(in *pb.ProductType) error
	TransDeleteProductType(ptid string) error

	CreateBanner(banner *pb.Banner) error
	UpdateBanner(updator, selector *pb.Banner) error
	DeleteBanner(banner *pb.Banner) error
	ListBanner(rq *pb.BannerRequest) ([]*pb.Banner, error)
	GetBanner(id string) (*pb.Banner, error)
	CountBanner(rq *pb.BannerRequest) (int64, error)

	CreateOrder(order *pb.Order) error
	UpdateOrder(updator, selector *pb.Order) error
	DeleteOrder(order *pb.Order) error
	ListOrder(rq *pb.OrderRequest) ([]*pb.Order, error)
	GetOrder(id string) (*pb.Order, error)
	CountOrder(rq *pb.OrderRequest) (int64, error)
	TransCreateOrder(order *pb.Order) error
	TransUpdateOrder(order *pb.Order) error

	CreateOrderDetail(orderDetail *pb.OrderDetail) error
	UpdateOrderDetail(updator, selector *pb.OrderDetail) error
	DeleteOrderDetail(orderDetail *pb.OrderDetail) error
	ListOrderDetail(rq *pb.OrderDetailRequest) ([]*pb.OrderDetail, error)
	GetOrderDetail(req *pb.OrderDetail) (*pb.OrderDetail, error)
	CountOrderDetail(rq *pb.OrderDetailRequest) (int64, error)

	CreateOrderShip(orderShip *pb.OrderShip) error
	UpdateOrderShip(updator, selector *pb.OrderShip) error
	DeleteOrderShip(orderShip *pb.OrderShip) error
	ListOrderShip(rq *pb.OrderShipRequest) ([]*pb.OrderShip, error)
	GetOrderShip(re *pb.OrderShip) (*pb.OrderShip, error)
	CountOrderShip(rq *pb.OrderShipRequest) (int64, error)

	CreateReview(review *pb.Review) error
	UpdateReview(updator, selector *pb.Review) error
	DeleteReview(review *pb.Review) error
	ListReview(rq *pb.ReviewRequest) ([]*pb.Review, error)
	GetReview(review *pb.Review) (*pb.Review, error)
	CountReview(rq *pb.ReviewRequest) (int64, error)
	IsReviewExist(review *pb.Review) bool
}

func NewRedisCache(addr, pw string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
	tick := time.NewTicker(10 * time.Minute)
	ctx := context.Background()
	go func(client *redis.Client) {
		for {
			select {
			case <-tick.C:
				if err := client.Ping(ctx).Err(); err != nil {
					panic(err)
				}
			}
		}
	}(client)
	return client
}

func NewProduct(cf *Configs) (*Product, error) {
	dbase := &db.DB{}
	if err := dbase.ConnectDb(cf.DBPath, cf.DBName); err != nil {
		return nil, err
	}
	log.Println("Connect db successful")
	redisDb, _ := strconv.Atoi(config.RedisDb)
	rd := NewRedisCache(config.RedisAddr, config.RedisPassword, redisDb)
	log.Println("Connect redis successful")
	return &Product{
		Db:    dbase,
		cache: rd,
	}, nil
}

func startGRPCServe(port string, p *Product) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: 15 * time.Second,
		}),
	}
	serve := grpc.NewServer(opts...)
	pb.RegisterProductServiceServer(serve, p)
	reflection.Register(serve)
	return serve.Serve(listen)
}
