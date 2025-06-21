package utils

import (
	"log"
	"testing"
)

func Test_slug(t *testing.T) {
	key := "Xịt/lăn khử mùi Rexona 72H kháng khuẩn thể thao dành cho nam 135/45ml"
	a := ToSlug(key)
	log.Println("a: ", a)
}

// func Test_convertSlug(*testing.T) {
// 	p := &Db{}
// 	if err := d.ConnectDb("admin_exchange:36b1c9722055507ac63cba13aa85d3fa17a555904df0c234@tcp(52.221.218.37:3306)", "voucher"); err != nil {
// 		log.Println(err)
// 	}
// }
