package main

import (
	"context"
	"sort"
	"time"

	pb "github.com/huyshop/header/product"
)

// GetReportOverview lấy tổng quan về doanh thu, số đơn hàng, trạng thái đơn hàng, số shop hoạt động, số user mới
func (p *Product) GetReportOverview(ctx context.Context, req *pb.ReportRequest) (*pb.ReportOverview, error) {
	// Tổng doanh thu và tổng số đơn hàng (chỉ tính đơn hoàn thành)
	orderReq := &pb.OrderRequest{
		State:   pb.Order_completed.String(),
		From:    req.StartDate,
		To:      req.EndDate,
		OrderBy: req.OrderBy,
	}
	orders, err := p.Db.ListOrder(orderReq)
	if err != nil {
		return nil, err
	}
	totalRevenue := int64(0)
	for _, o := range orders {
		totalRevenue += int64(o.TotalMoney)
	}
	totalOrders := int32(len(orders))

	// Đếm trạng thái đơn hàng
	orderStatus := &pb.OrderStatusCount{}
	// Completed
	orderStatus.Completed = totalOrders
	// Processing (ví dụ: đang chờ xử lý)
	processingReq := &pb.OrderRequest{
		State:   pb.Order_pending.String(),
		From:    req.StartDate,
		To:      req.EndDate,
		OrderBy: req.OrderBy,
	}
	processingOrders, _ := p.Db.ListOrder(processingReq)
	orderStatus.Processing = int32(len(processingOrders))
	// Cancelled
	cancelReq := &pb.OrderRequest{
		State:   pb.Order_cancelled.String(),
		From:    req.StartDate,
		To:      req.EndDate,
		OrderBy: req.OrderBy,
	}
	cancelOrders, _ := p.Db.ListOrder(cancelReq)
	orderStatus.Cancelled = int32(len(cancelOrders))

	return &pb.ReportOverview{
		TotalRevenue: totalRevenue,
		TotalOrders:  totalOrders,
		OrderStatus:  orderStatus,
	}, nil
}

// GetReportRevenue lấy doanh thu theo ngày/tháng/năm
func (p *Product) GetReportRevenue(ctx context.Context, req *pb.ReportRequest) (*pb.ReportRevenue, error) {
	// Gom nhóm doanh thu theo ngày/tháng/năm dựa vào GroupBy
	orderReq := &pb.OrderRequest{
		State:   pb.Order_completed.String(),
		OrderBy: req.OrderBy,
	}
	orders, err := p.Db.ListOrder(orderReq)
	if err != nil {
		return nil, err
	}
	groupBy := req.GroupBy // "day", "month", "year"
	labelMap := map[string]int64{}
	for _, o := range orders {
		t := time.Unix(o.TimeOrder, 0)
		var label string
		switch groupBy {
		case "month":
			label = t.Format("2006-01")
		case "year":
			label = t.Format("2006")
		default:
			label = t.Format("2006-01-02")
		}
		labelMap[label] += int64(o.TotalMoney)
	}
	// Sắp xếp label tăng dần
	labels := make([]string, 0, len(labelMap))
	for k := range labelMap {
		labels = append(labels, k)
	}
	sort.Strings(labels)
	values := make([]int64, 0, len(labels))
	for _, l := range labels {
		values = append(values, labelMap[l])
	}
	return &pb.ReportRevenue{
		Labels: labels,
		Values: values,
	}, nil
}

// GetReportStoreRevenue lấy doanh thu theo từng cửa hàng
func (p *Product) GetReportStoreRevenue(ctx context.Context, req *pb.ReportRequest) (*pb.ReportStoreRevenue, error) {
	// Gom nhóm doanh thu theo từng cửa hàng
	orderReq := &pb.OrderRequest{
		State:   pb.Order_completed.String(),
		From:    req.StartDate,
		To:      req.EndDate,
		OrderBy: req.OrderBy,
	}
	orders, err := p.Db.ListOrder(orderReq)
	if err != nil {
		return nil, err
	}
	storeMap := map[string]*pb.StoreRevenue{}
	for _, o := range orders {
		storeId := o.StoreId
		if storeId == "" {
			continue
		}
		if _, ok := storeMap[storeId]; !ok {
			storeMap[storeId] = &pb.StoreRevenue{
				StoreId: storeId,
			}
		}
		sr := storeMap[storeId]
		sr.TotalOrders++
		sr.TotalRevenue += int64(o.TotalMoney)
	}
	// Tính tỷ lệ hủy đơn cho từng shop (giả lập 0)
	storeRevenues := make([]*pb.StoreRevenue, 0, len(storeMap))
	totalRevenue := int64(0)
	for _, sr := range storeMap {
		sr.CancelledRate = 0
		storeRevenues = append(storeRevenues, sr)
		totalRevenue += sr.TotalRevenue
	}
	return &pb.ReportStoreRevenue{
		StoreRevenues: storeRevenues,
		TotalStores:   int32(len(storeRevenues)),
		TotalRevenue:  totalRevenue,
	}, nil
}
