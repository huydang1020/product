# Product Service

Service **Product** là một phần của hệ thống Huyshop (Microservices architecture), được phát triển bằng Go (Golang) và giao tiếp chủ yếu qua gRPC. Service này chịu trách nhiệm xử lý các nghiệp vụ liên quan đến sản phẩm, danh mục, đơn hàng, giỏ hàng, đánh giá và báo cáo.

## 🚀 Tính năng chính

- **Quản lý sản phẩm (Products & Product Types)**: Thêm, sửa, xóa, tìm kiếm thông tin sản phẩm.
- **Quản lý danh mục (Categories)**: Phân loại sản phẩm.
- **Quản lý đơn hàng (Orders)**: Xử lý quy trình đặt hàng.
- **Giỏ hàng (Cart)**: Quản lý giỏ hàng của người dùng (sử dụng Redis để tối ưu tốc độ và tự động hết hạn).
- **Quản lý Banner**: Cấu hình banner quảng cáo/trưng bày.
- **Đánh giá (Reviews)**: Hệ thống review/feedback sản phẩm từ người dùng.
- **Báo cáo (Reports)**: Thống kê dữ liệu sản phẩm, đơn hàng.

## 🛠 Yêu cầu hệ thống

- **Go**: 1.18+ (khuyến nghị)
- **Database**: MySQL (lưu trữ dữ liệu chính)
- **Cache/Session**: Redis (lưu trữ giỏ hàng và các dữ liệu cache)
- **Docker & Docker Compose** (tùy chọn cho việc deploy)

## ⚙️ Cấu hình môi trường

Tạo file `.env` ở thư mục gốc của project (có thể tham khảo file `.env.example`) với các cấu hình sau:

```env
# Cổng chạy gRPC server
GRPC_PORT=8000

# Cấu hình kết nối MySQL (Format: user:password@tcp(host:port))
DB_PATH=root:123456@tcp(localhost:3306)
DB_NAME=product

# Cấu hình kết nối Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_CART_EXPIRE=3600 # Thời gian sống của giỏ hàng (giây)

# Cấu hình kết nối tới User Service
USER_HOST=localhost:6001
```

## 📦 Hướng dẫn chạy và phát triển (Development)

Dự án sử dụng `Makefile` và `urfave/cli/v2` để quản lý các lệnh chạy.

### 1. Khởi tạo Database (Tạo bảng tự động)
Trước khi chạy ứng dụng lần đầu tiên, bạn cần khởi tạo các bảng trong Database:
```bash
make cdb
# hoặc chạy lệnh trực tiếp: go build && ./product createDb
```

### 2. Khởi chạy Service
Chạy service ở chế độ bình thường:
```bash
make start
# hoặc chạy lệnh trực tiếp: go build && ./product start
```

### 3. Build file thực thi (Binary)
Build ứng dụng thành file thực thi cho môi trường Linux (sử dụng cho Docker/Production):
```bash
make build
# Lệnh thực thi: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o product .
```

## 🐳 Khởi chạy với Docker

Dự án đã có sẵn `Dockerfile` (sử dụng base image `alpine:3.14`) tối ưu dung lượng cho môi trường production.

1. **Build image**:
```bash
make build # Build file binary (product) trước
docker build -t huyshop/product:latest .
```

2. **Chạy container**:
```bash
docker run -d \
  -p 8000:8000 \
  --name product_service \
  --env-file .env \
  huyshop/product:latest
```

## 🧹 Quản lý bộ nhớ (Memory Management)
Service được tích hợp sẵn một goroutine chạy nền để tự động thu gom rác (Garbage Collection - GC) và giải phóng bộ nhớ cho hệ điều hành (`debug.FreeOSMemory()`) định kỳ mỗi 15 phút, giúp tối ưu RAM khi chạy trên môi trường có tài nguyên hạn chế.
