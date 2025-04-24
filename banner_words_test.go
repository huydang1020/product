package main

import (
	"log"
	"testing"

	"github.com/huyshop/product/utils"
)

func Test_checkBannerWords(t *testing.T) {
	a := "VD LUXURY - Ví da cao cấp: Chuyên cung cấp túi xách, ví da, balo, túi xách du lịch, phụ kiện thắt lưng nam nữ đạt tiêu chuẩn cao, chất liệu cao cấp, hàng loại một. Tuyệt đối nói KHÔNG với hàng gia công, kém chất lượng. VD LUXURY - Ví da cao cấp cam kết: Sản phẩm giống 100% ảnh về kiểu dáng, hàng loại 1 nhập trực tiếp tại xưởng sản xuất; Có thể có sai lệch về màu sắc nhưng không đáng kể; Bảo hành đổi trả hoàn tiền trong 03 ngày nếu phát sai về kiểu dáng, màu sắc, số lượng sản phẩm. Mô tả sản phẩm: Chất liệu da cao cấp - Sản phẩm ví da được làm từ da cao cấp mềm mại, sang trọng và có thời gian sử dụng lâu dài. Màu sắc đen nâu nam tính, dễ phối đồ. Thiết kế ví tinh tế, nam tính. Ví có kiểu dáng ngang đơn giản và nam tính nhưng tiện lợi trong việc lưu giữ tiền bạc, giấy tờ xe, thẻ ATM, hình ảnh lưu niệm. Kích thước: 12 x 10 cm (gấp lại); 24 x 10 cm (mở ra). Số ngăn: 02 ngăn lớn và 08 ngăn nhỏ đựng thẻ, giấy tờ. Hướng dẫn bảo quản: 1. Chú ý đến thời tiết - Tránh để ví ra ngoài khi trời mưa hoặc thời tiết xấu để chúng không bị ướt dẫn đến bong tróc. Mua thêm túi xách plastic bao bọc khi thời tiết xấu đột ngột. Cất giữ ví da ở nơi thoáng mát. 2. Hạn chế ngồi đè lên ví - Ngồi đè có thể làm biến dạng ví và giấy tờ, nên để ở túi trước thay vì túi sau. 3. Rửa tay trước khi cầm ví da - Mồ hôi và chất nhờn có thể làm hỏng ví. 4. Xịt dung dịch bảo vệ da - Trước khi sử dụng, xịt toàn bộ bề mặt ví, để qua đêm cho khô. 5. Không để quá nhiều tiền và thẻ trong ví - Kiểm soát trọng lượng khi sử dụng. 6. Sửa chữa ngay các chi tiết bị hỏng - Chú ý các chi tiết nhỏ như dây kéo, nút gài, đường chỉ may để tránh hư hỏng nặng hơn."
	check := utils.ContainsBannedWords(a)
	if check {
		log.Println("false")
		return
	}
	log.Println("oke")
}
