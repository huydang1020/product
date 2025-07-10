package main

import (
	"context"
	"log"
	"strconv"
	"time"

	pb "github.com/huyshop/header/product"
)

// GetReportOverview lấy tổng quan về doanh thu, số đơn hàng, trạng thái đơn hàng, số shop hoạt động, số user mới
func (p *Product) GetReportOverview(ctx context.Context, req *pb.ReportRequest) (*pb.ReportOverview, error) {
	log.Println("GetReportOverview called with request:", req)
	// Tổng doanh thu và tổng số đơn hàng (chỉ tính đơn hoàn thành)
	orderReq := &pb.OrderRequest{
		From:      req.StartDate,
		To:        req.EndDate,
		OrderBy:   req.OrderBy,
		PartnerId: req.PartnerId,
		StoreId:   req.StoreId,
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
	orderStatus := &pb.OrderStatusCount{}
	// Đếm trạng thái đơn hàng
	for _, ord := range orders {
		switch ord.State {
		case pb.Order_pending.String():
			orderStatus.Pending++
		case pb.Order_confirmed.String():
			orderStatus.Confirmed++
		case pb.Order_shipping.String():
			orderStatus.Shipping++
		case pb.Order_completed.String():
			orderStatus.Completed++
		case pb.Order_processing.String():
			orderStatus.Processing++
		case pb.Order_cancelled.String():
			orderStatus.Cancelled++
		}
	}
	log.Println("orderStatus:", orderStatus)
	listPty, err := p.Db.CountProductType(&pb.ProductTypeRequest{
		PartnerId: req.PartnerId,
		From:      req.StartDate,
		To:        req.EndDate,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReportOverview{
		TotalRevenue: totalRevenue,
		TotalOrders:  totalOrders,
		OrderStatus:  orderStatus,
		TotalProduct: int32(listPty),
	}, nil
}

// GetReportRevenue lấy doanh thu theo ngày/tháng/năm
func (p *Product) GetReportRevenue(ctx context.Context, req *pb.ReportRequest) (*pb.ReportRevenue, error) {
	month := req.Month
	year := req.Year
	now := time.Now()
	if month == 0 {
		month = int32(now.Month())
	}
	if year == 0 {
		year = int32(now.Year())
	}

	orderReq := &pb.OrderRequest{
		State:     pb.Order_completed.String(),
		PartnerId: req.PartnerId,
	}
	orders, err := p.Db.ListOrder(orderReq)
	if err != nil {
		return nil, err
	}
	groupBy := req.GroupBy // "day", "month", "year"

	labels := make([]string, 0)
	values := make([]int64, 0)

	switch groupBy {
	case "month":
		// Chỉ tạo labels cho tháng/năm được truyền vào, label là số ngày
		daysInMonth := time.Date(int(year), time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
		for day := 1; day <= daysInMonth; day++ {
			label := ""
			if day < 10 {
				label = "0" + strconv.Itoa(day)
			} else {
				label = strconv.Itoa(day)
			}
			labels = append(labels, label)
			// Tìm doanh thu theo ngày cụ thể từ orders
			dateStr := time.Date(int(year), time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
			dailyRevenue := int64(0)
			for _, order := range orders {
				orderTime := time.Unix(order.TimeOrder, 0)
				if orderTime.Format("2006-01-02") == dateStr {
					dailyRevenue += int64(order.TotalMoney)
				}
			}
			values = append(values, dailyRevenue)
		}

	case "year":
		// Label là số tháng ("01", "02", ... "12")
		for m := 1; m <= 12; m++ {
			label := ""
			if m < 10 {
				label = "0" + strconv.Itoa(m)
			} else {
				label = strconv.Itoa(m)
			}
			labels = append(labels, label)
			monthStr := time.Date(int(year), time.Month(m), 1, 0, 0, 0, 0, time.UTC).Format("2006-01")
			monthlyRevenue := int64(0)
			for _, order := range orders {
				orderTime := time.Unix(order.TimeOrder, 0)
				if orderTime.Format("2006-01") == monthStr {
					monthlyRevenue += int64(order.TotalMoney)
				}
			}
			values = append(values, monthlyRevenue)
		}

	default: // "day" hoặc không truyền
		// Lấy 7 ngày gần nhất
		last7Days := make([]string, 0, 7)
		for i := 6; i >= 0; i-- {
			label := now.AddDate(0, 0, -i).Format("2006-01-02")
			last7Days = append(last7Days, label)
		}
		labels = append(labels, last7Days...)
		for _, label := range last7Days {
			dailyRevenue := int64(0)
			for _, order := range orders {
				orderTime := time.Unix(order.TimeOrder, 0)
				if orderTime.Format("2006-01-02") == label {
					dailyRevenue += int64(order.TotalMoney)
				}
			}
			values = append(values, dailyRevenue)
		}
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
	// Thêm map đếm trạng thái đơn hàng cho từng shop
	statusMap := map[string]*pb.OrderStatusCount{}
	for _, o := range orders {
		storeId := o.StoreId
		if storeId == "" {
			continue
		}
		if _, ok := storeMap[storeId]; !ok {
			storeMap[storeId] = &pb.StoreRevenue{
				StoreId: storeId,
			}
			statusMap[storeId] = &pb.OrderStatusCount{}
		}
		sr := storeMap[storeId]
		sr.TotalOrders++
		sr.TotalRevenue += int64(o.TotalMoney)

		// Đếm trạng thái đơn hàng
		st := o.State
		if st == pb.Order_completed.String() {
			statusMap[storeId].Completed++
		} else if st == pb.Order_pending.String() {
			statusMap[storeId].Processing++
		} else if st == pb.Order_cancelled.String() {
			statusMap[storeId].Cancelled++
		}
	}
	// Tính tỷ lệ hủy đơn cho từng shop (giả lập 0)
	storeRevenues := make([]*pb.StoreRevenue, 0, len(storeMap))
	totalRevenue := int64(0)
	for storeId, sr := range storeMap {
		sr.CancelledRate = 0
		sr.OrderStatus = statusMap[storeId]
		storeRevenues = append(storeRevenues, sr)
		totalRevenue += sr.TotalRevenue
	}
	return &pb.ReportStoreRevenue{
		StoreRevenues: storeRevenues,
		TotalStores:   int32(len(storeRevenues)),
		TotalRevenue:  totalRevenue,
	}, nil
}
