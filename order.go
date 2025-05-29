package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/huyshop/header/common"
	pb "github.com/huyshop/header/product"
	"github.com/huyshop/product/utils"
)

const (
	REDIS_KEY_CART  = "cart_"
	REDIS_KEY_ORDER = "order_"
	ACTIVE          = "active"
	INACTIVE        = "inactive"
	COD             = "cod"
)

func (p *Product) CreateOrder(ctx context.Context, req *pb.Order) (*pb.Order, error) {
	if req.GetUserId() == "" {
		return nil, errors.New(utils.E_not_found_user_id)
	}
	log.Println("req:", req)
	if len(req.ProductOrdered) < 1 {
		log.Println("not found product order")
		return nil, errors.New(utils.E_not_found_product)
	}
	if req.GetReceiverName() == "" {
		return nil, errors.New(utils.E_invalid_receiver_name)
	}
	if req.GetReceiverPhone() == "" {
		return nil, errors.New(utils.E_invalid_receiver_phone)
	}
	if req.GetReceiverAddress() == "" {
		return nil, errors.New(utils.E_invalid_receiver_address)
	}

	var totalMoney float64
	for _, ord := range req.ProductOrdered {
		log.Println("ord:", ord)
		prod, err := p.Db.GetProduct(ord.ProductId)
		if err != nil {
			log.Println("get prod err:", err)
			return nil, errors.New(utils.E_not_found_product)
		}
		if ord.Quantity < 1 {
			log.Println("invalid amount product")
			return nil, errors.New(utils.E_invalid_amount_product)
		}
		if ord.Quantity > prod.GetQuantity() {
			log.Println("inventory quantity not enough")
			return nil, errors.New(utils.E_inventory_quantity_not_enough)
		}
		price := float64(prod.SellPrice) * float64(ord.Quantity)
		totalMoney += price
	}
	req.Id = utils.MakeOrderId()
	randNumber := rand.Intn(99999999999999-10000000000000) + 10000000000000
	req.OrderCode = fmt.Sprint(randNumber)
	req.TimeOrder = time.Now().Unix()
	req.State = pb.Order_pending.String()
	req.TotalMoney = totalMoney
	req.ShippingFee = 30000
	history := map[string]int64{}
	history[pb.Order_pending.String()] = req.TimeOrder
	byteHistory, err := json.Marshal(history)
	if err != nil {
		log.Println("marshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	req.History = string(byteHistory)
	// if req.CodeId != "" {
	// 	code, err := p.Db.GetCode(&pb.Code{Id: req.CodeId, UserId: req.UserId})
	// 	if err != nil {
	// 		log.Println("get code err:", err)
	// 		ctx.JSON(200, &Response{Code: -1, Message: err.Error()})
	// 		return
	// 	}
	// 	if code.State != GOT {
	// 		ctx.JSON(200, &Response{Code: -1, Message: utils.INVALID_CODE_STATE})
	// 		return
	// 	}
	// 	vou, err := r.db.GetVoucher(&pb.Voucher{Id: code.VoucherId})
	// 	if err != nil {
	// 		log.Println("get voucher err:", err)
	// 		ctx.JSON(200, &Response{Code: -1, Message: err.Error()})
	// 		return
	// 	}
	// 	discount := float64(vou.Discount) / float64(100)
	// 	req.TotalMoney = req.TotalMoney - (req.TotalMoney * discount)
	// }
	req.TotalMoney = req.TotalMoney + float64(req.ShippingFee)
	if req.MethodPayment == COD {
		if err := p.Db.TransCreateOrder(req); err != nil {
			log.Println("insert order err:", err)
			return nil, errors.New(utils.E_internal_error)
		}
		return &pb.Order{}, nil
	} else if req.MethodPayment == "online" {
		vnpUrl := os.Getenv("VNP_URL")
		vnpSecret := os.Getenv("VNP_HASH_SECRET")
		vnpTmnCode := os.Getenv("VNP_TMNCODE")
		createdDate, err := utils.ConvertUnixToDateTime("20060102150405", req.TimeOrder)
		if err != nil {
			log.Println("convert time err:", err)
			return nil, errors.New(utils.E_internal_error)
		}
		vnpParams := url.Values{}
		vnpParams.Set("vnp_Version", "2.1.0")
		vnpParams.Set("vnp_Command", "pay")
		vnpParams.Set("vnp_TmnCode", vnpTmnCode)
		vnpParams.Set("vnp_Locale", "vn")
		vnpParams.Set("vnp_CurrCode", "VND")
		vnpParams.Set("vnp_TxnRef", req.OrderCode)
		vnpParams.Set("vnp_OrderInfo", "Thanh toán cho giao dịch: "+req.OrderCode)
		vnpParams.Set("vnp_OrderType", "billpayment")
		vnpParams.Set("vnp_Amount", strconv.FormatInt(int64(req.TotalMoney)*100, 10))
		vnpParams.Set("vnp_ReturnUrl", req.VnpayReturnUrl)
		vnpParams.Set("vnp_IpAddr", req.IpAddress)
		vnpParams.Set("vnp_CreateDate", createdDate)
		vnpParams.Set("vnp_BankCode", "VNBANK")

		sortedParams := utils.SortParams(vnpParams)
		signData := sortedParams.Encode()
		hmacSecret := hmac.New(sha512.New, []byte(vnpSecret))
		hmacSecret.Write([]byte(signData))
		signature := fmt.Sprintf("%x", hmacSecret.Sum(nil))
		vnpParams.Set("vnp_SecureHash", signature)
		vnpRedirectURL := vnpUrl + "?" + vnpParams.Encode()
		byteOrder, err := json.Marshal(req)
		if err != nil {
			log.Println("marshal order err:", err)
			return nil, errors.New(utils.E_internal_error)
		}
		exprOderRedis, _ := strconv.Atoi(os.Getenv("TIME_LIVE_ORDER_REDIS"))
		keyRedis := REDIS_KEY_ORDER + req.OrderCode
		if err := p.cache.Set(ctx, keyRedis, string(byteOrder), time.Duration(exprOderRedis)*time.Second).Err(); err != nil {
			log.Println("set data redis error:", err)
			return nil, errors.New(utils.E_internal_error)
		}
		return &pb.Order{VnpRedirectUrl: vnpRedirectURL}, nil
	}
	return &pb.Order{}, errors.New(utils.E_invalid_method_payment)
}

func (p *Product) CreateOrderVNpay(ctx context.Context, req *pb.Order) (*common.Empty, error) {
	if req.GetUserId() == "" {
		log.Println("not found user id")
		return nil, errors.New(utils.E_not_found_user_id)
	}
	if req.GetOrderCode() == "" {
		log.Println("not found order code")
		return nil, errors.New(utils.E_not_found_order_code)
	}
	keyRedis := REDIS_KEY_ORDER + req.OrderCode
	result, err := p.cache.Get(ctx, keyRedis).Result()
	if err != nil {
		log.Println("get data redis error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	order := &pb.Order{}
	if err := json.Unmarshal([]byte(result), order); err != nil {
		log.Println("unmarshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	order.State = pb.Order_completed.String()
	if err := p.Db.TransCreateOrder(order); err != nil {
		log.Println("trans insert order err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	if err := p.cache.Del(ctx, keyRedis); err != nil {
		log.Println("del key redis err:", err)
		// return nil, errors.New(utils.E_internal_error)
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
	if len(list) < 1 {
		log.Println("ListOrder empty")
		return &pb.Orders{}, nil
	}
	for _, order := range list {
		if order.GetProductOrdered() != nil {
			for _, item := range order.GetProductOrdered() {
				prod, err := p.Db.GetProduct(item.GetProductId())
				if err != nil {
					log.Println("GetProduct error:", err)
					return nil, errors.New(utils.E_internal_error)
				}
				item.Product = prod
			}
		}
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
	if order.GetProductOrdered() != nil {
		for _, item := range order.GetProductOrdered() {
			prod, err := p.Db.GetProduct(item.GetProductId())
			if err != nil {
				log.Println("GetProduct error:", err)
				return nil, errors.New(utils.E_internal_error)
			}
			item.Product = prod
		}
	}
	return order, nil
}

func (p *Product) CancelOrder(ctx context.Context, req *pb.Order) (*common.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	order, err := p.Db.GetOrder(req.GetId())
	if err != nil {
		log.Println("GetOrder error:", err)
		return nil, errors.New(utils.E_not_found_order)
	}
	if order.State != pb.Order_completed.String() {
		return nil, errors.New(utils.E_invalid_state)
	}
	order.State = pb.Order_canceled.String()
	history := map[string]int64{}
	history[pb.Order_canceled.String()] = req.TimeOrder
	byteHistory, err := json.Marshal(history)
	if err != nil {
		log.Println("marshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	order.History = string(byteHistory)
	if err := p.Db.UpdateOrder(order, &pb.Order{Id: req.Id}); err != nil {
		log.Println("UpdateOrder error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return &common.Empty{}, nil
}

func (p *Product) UpdateStateOrder(ctx context.Context, req *pb.Order) (*common.Empty, error) {
	if req.GetId() == "" {
		return nil, errors.New(utils.E_not_found_id)
	}
	order, err := p.Db.GetOrder(req.GetId())
	if err != nil {
		log.Println("GetOrder error:", err)
		return nil, errors.New(utils.E_not_found_order)
	}
	if order.State != pb.Order_completed.String() || order.State != pb.Order_canceled.String() {
		return nil, errors.New(utils.E_invalid_state)
	}
	history := map[string]int64{}
	history[req.GetState()] = req.TimeOrder
	byteHistory, err := json.Marshal(history)
	if err != nil {
		log.Println("marshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	order.History = string(byteHistory)
	if err := p.Db.UpdateOrder(order, &pb.Order{Id: req.Id}); err != nil {
		log.Println("UpdateOrder error:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return &common.Empty{}, nil
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
		prod, err := p.Db.GetProduct(item.GetProductId())
		if err != nil {
			log.Println("GetProduct error:", err)
			return nil, errors.New(utils.E_internal_error)
		}
		if item.GetQuantity() > prod.GetQuantity() {
			return nil, errors.New(utils.E_inventory_quantity_not_enough)
		}
	}

	keyCartRedis := REDIS_KEY_CART + fmt.Sprint(req.GetUserId())
	cartExpire, _ := strconv.Atoi(config.RedisCartExpire)
	durationCartRedis := time.Duration(cartExpire) * time.Minute
	result, err := p.cache.Get(ctx, keyCartRedis).Result()
	if err != redis.Nil && err != nil {
		log.Println("redis err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	itemCart := make(map[string]int32, 0)
	if err == redis.Nil {
		log.Println("add new cart")
		for _, item := range req.Item {
			itemCart[item.ProductId] = item.Quantity
		}
		byteItem, err := json.Marshal(itemCart)
		if err != nil {
			log.Println("marshal err:", err)
			return nil, errors.New(utils.E_internal_error)
		}
		if err := p.cache.Set(ctx, keyCartRedis, string(byteItem), durationCartRedis).Err(); err != nil {
			log.Println("set cart redis err:", err)
			return nil, errors.New(utils.E_internal_error)
		}
		resp := []*pb.ProductOrdered{}
		for prodId, quantity := range itemCart {
			cartItem := &pb.ProductOrdered{
				ProductId: prodId,
				Quantity:  quantity,
			}
			prod, err := p.Db.GetProduct(prodId)
			if err != nil {
				log.Println("get prod err:", err)
				continue
			}
			cartItem.Product = prod
			resp = append(resp, cartItem)
		}
		return &pb.Cart{Item: resp}, nil
	}
	log.Println("update cart")
	if err := json.Unmarshal([]byte(result), &itemCart); err != nil {
		log.Println("unmarshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	for _, item := range req.Item {
		itemCart[item.ProductId] = item.Quantity
	}
	byteDataCart, _ := json.Marshal(itemCart)
	if err := p.cache.Set(ctx, keyCartRedis, string(byteDataCart), durationCartRedis).Err(); err != nil {
		log.Println("set cart redis err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	newResult, err := p.cache.Get(ctx, keyCartRedis).Result()
	if err != redis.Nil && err != nil {
		log.Println("redis err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	if err := json.Unmarshal([]byte(newResult), &itemCart); err != nil {
		log.Println("unmarshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	resp := []*pb.ProductOrdered{}
	for prodId, quantity := range itemCart {
		cartItem := &pb.ProductOrdered{
			ProductId: prodId,
			Quantity:  quantity,
		}
		prod, err := p.Db.GetProduct(prodId)
		if err != nil {
			log.Println("get prod err:", err)
			continue
		}
		cartItem.Product = prod
		resp = append(resp, cartItem)
	}
	return &pb.Cart{Item: resp, UserId: req.GetUserId()}, nil
}

func (p *Product) DeleteCart(ctx context.Context, req *pb.Cart) (*pb.Cart, error) {
	if req.GetUserId() == "" {
		return nil, errors.New(utils.E_not_found_user_id)
	}
	keyCartRedis := REDIS_KEY_CART + fmt.Sprint(req.GetUserId())
	if err := p.cache.Del(ctx, keyCartRedis).Err(); err != nil {
		log.Println("delete cart redis err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return &pb.Cart{}, nil
}

func (p *Product) ListCart(ctx context.Context, req *pb.Cart) (*pb.Cart, error) {
	log.Println("ListCart", req)
	if req.GetUserId() == "" {
		return nil, errors.New(utils.E_not_found_user_id)
	}
	keyCartRedis := REDIS_KEY_CART + fmt.Sprint(req.GetUserId())
	result, err := p.cache.Get(ctx, keyCartRedis).Result()
	if err != redis.Nil && err != nil {
		log.Println("redis err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	if err == redis.Nil {
		log.Println("cart not found")
		return &pb.Cart{}, nil
	}
	log.Println("cart redis:", result)
	itemCart := make(map[string]int, 0)
	if err := json.Unmarshal([]byte(result), &itemCart); err != nil {
		log.Println("unmarshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	resp := []*pb.ProductOrdered{}
	for prodId, quantity := range itemCart {
		if len(req.ProductIds) > 0 && !utils.Include(req.GetProductIds(), prodId) {
			continue
		}
		cartItem := &pb.ProductOrdered{
			ProductId: prodId,
			Quantity:  int32(quantity),
		}
		prod, err := p.Db.GetProduct(prodId)
		if err != nil {
			log.Println("get prod err:", err)
			continue
		}
		cartItem.Product = prod
		resp = append(resp, cartItem)
		sort.SliceStable(resp, func(i, j int) bool {
			return resp[i].ProductId < resp[j].ProductId
		})
	}
	log.Println("ok")
	return &pb.Cart{Item: resp, UserId: req.GetUserId()}, nil
}

func (p *Product) DeleteCartItem(ctx context.Context, req *pb.Cart) (*pb.Cart, error) {
	if req.GetUserId() == "" {
		return nil, errors.New(utils.E_not_found_user_id)
	}
	proId := req.GetItem()[0].ProductId
	if proId == "" {
		return nil, errors.New(utils.E_not_found_product_id)
	}
	_, err := p.Db.GetProduct(proId)
	if err != nil {
		return nil, errors.New(utils.E_not_found_product)
	}
	keyCartRedis := REDIS_KEY_CART + fmt.Sprint(req.GetUserId())
	result, err := p.cache.Get(ctx, keyCartRedis).Result()
	if err != redis.Nil && err != nil {
		log.Println("redis err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	if err == redis.Nil {
		return nil, errors.New(utils.E_not_found_item_cart)
	}
	itemCart := make(map[string]int, 0)
	if err := json.Unmarshal([]byte(result), &itemCart); err != nil {
		log.Println("unmarshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	delete(itemCart, proId)
	byteItem, err := json.Marshal(itemCart)
	log.Println("item cart:", itemCart)
	if err != nil {
		log.Println("marshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	cartExpire, _ := strconv.Atoi(os.Getenv("REDIS_CART_EXPIRE"))
	durationCartRedis := time.Duration(cartExpire) * time.Minute
	if err := p.cache.Set(ctx, keyCartRedis, string(byteItem), durationCartRedis).Err(); err != nil {
		log.Println("set cart redis err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return &pb.Cart{Item: req.Item, UserId: req.GetUserId()}, nil
}

func (p *Product) DeleteCartItems(ctx context.Context, req *pb.Cart) (*pb.Cart, error) {
	if req.GetUserId() == "" {
		return nil, errors.New(utils.E_not_found_user_id)
	}
	if len(req.Item) < 1 {
		return nil, errors.New(utils.E_not_found_item_cart)
	}
	keyCartRedis := REDIS_KEY_CART + fmt.Sprint(req.GetUserId())
	result, err := p.cache.Get(ctx, keyCartRedis).Result()
	if err != redis.Nil && err != nil {
		log.Println("redis err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	if err == redis.Nil {
		return nil, errors.New(utils.E_not_found_item_cart)
	}
	itemCart := make(map[string]int, 0)
	if err := json.Unmarshal([]byte(result), &itemCart); err != nil {
		log.Println("unmarshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	for _, item := range req.Item {
		delete(itemCart, item.ProductId)
	}
	byteItem, err := json.Marshal(itemCart)
	if err != nil {
		log.Println("marshal err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	cartExpire, _ := strconv.Atoi(os.Getenv("REDIS_CART_EXPIRE"))
	durationCartRedis := time.Duration(cartExpire) * time.Minute
	if err := p.cache.Set(ctx, keyCartRedis, string(byteItem), durationCartRedis).Err(); err != nil {
		log.Println("set cart redis err:", err)
		return nil, errors.New(utils.E_internal_error)
	}
	return &pb.Cart{Item: req.Item}, nil
}
