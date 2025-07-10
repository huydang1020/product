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
		PartnerId: req.PartnerId,
		
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
	// Xác định tháng/năm cần lấy
	month := req.Month
	year := req.Year
	now := time.Now()
	if month == 0 {
		month = int32(now.Month())
	}
	if year == 0 {
		year = int32(now.Year())
	}

	// Gom nhóm doanh thu theo ngày/tháng/năm dựa vào GroupBy
	orderReq := &pb.OrderRequest{
		State:     pb.Order_completed.String(),
		PartnerId: req.PartnerId,
	}
	orders, err := p.Db.ListOrder(orderReq)
	if err != nil {
		return nil, err
	}
	groupBy := req.GroupBy // "day", "month", "year"

	// Tạo map để lưu trữ doanh thu theo từng đơn vị thời gian
	revenueMap := make(map[string]int64)

	// Tính toán doanh thu theo từng đơn vị thời gian
	for _, o := range orders {
		t := time.Unix(o.TimeOrder, 0)
		var key string
		switch groupBy {
		case "month":
			// Lọc đúng tháng/năm
			if int32(t.Month()) != month || int32(t.Year()) != year {
				continue
			}
			key = t.Format("2006-01")
		case "year":
			if int32(t.Year()) != year {
				continue
			}
			key = t.Format("2006")
		default:
			key = t.Format("2006-01-02")
		}
		revenueMap[key] += int64(o.TotalMoney)
	}

	// Tạo labels và values dựa trên groupBy
	labels := make([]string, 0)
	values := make([]int64, 0)

	switch groupBy {
	case "month":
		// Chỉ tạo labels cho tháng/năm được truyền vào
		daysInMonth := time.Date(int(year), time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
		for day := 1; day <= daysInMonth; day++ {
			dayLabel := time.Date(int(year), time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
			labels = append(labels, dayLabel)
			// Tìm doanh thu theo ngày cụ thể từ orders
			dailyRevenue := int64(0)
			for _, order := range orders {
				orderTime := time.Unix(order.TimeOrder, 0)
				if orderTime.Format("2006-01-02") == dayLabel {
					dailyRevenue += int64(order.TotalMoney)
				}
			}
			values = append(values, dailyRevenue)
		}

	case "year":
		// Tạo tất cả các tháng trong năm cho mỗi năm có dữ liệu
		for month := 1; month <= 12; month++ {
			monthLabel := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("2006-01")
			labels = append(labels, monthLabel)
			// Tìm doanh thu theo tháng cụ thể từ orders
			monthlyRevenue := int64(0)
			for _, order := range orders {
				orderTime := time.Unix(order.TimeOrder, 0)
				if orderTime.Format("2006-01") == monthLabel {
					monthlyRevenue += int64(order.TotalMoney)
				}
			}
			values = append(values, monthlyRevenue)
		}

	default: // "day"
		// Tạo tất cả các ngày có dữ liệu
		for dayKey := range revenueMap {
			labels = append(labels, dayKey)
			values = append(values, revenueMap[dayKey])
		}
	}

	// Sắp xếp label tăng dần
	sort.Strings(labels)
	// Sắp xếp lại values theo thứ tự của labels đã sort
	sortedValues := make([]int64, len(labels))
	for i, label := range labels {
		// Tìm lại giá trị tương ứng
		for j, originalLabel := range labels {
			if originalLabel == label {
				sortedValues[i] = values[j]
				break
			}
		}
	}

	return &pb.ReportRevenue{
		Labels: labels,
		Values: sortedValues,
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
