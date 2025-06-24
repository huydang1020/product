package db_test

import (
	"log"
	"testing"

	"github.com/huyshop/product/db"
)

func Test_addStoreIdToOrder(t *testing.T) {
	p := &db.DB{}
	if err := p.ConnectDb("root:123456@tcp(localhost:3306)", "product"); err != nil {
		log.Println(err)
		return
	}
	// listord, err := p.ListOrder(&product.OrderRequest{})
	// if err != nil {
	// 	log.Println("err", err)
	// 	return
	// }
	// for _, or := range listord {
	// 	// if pty.Slug != "" {
	// 	// 	continue
	// 	// }
	// 	for _, pro := range or.ProductOrdered {
	// 		pr, err := p.GetProduct(pro.ProductId)
	// 		if err != nil {
	// 			log.Println("err: ", err, pro.ProductId)
	// 			return
	// 		}
	// 		pty, err := p.GetProductType(pr.ProductTypeId)
	// 		if err != nil {
	// 			log.Println("err: ", err, pr.Id)
	// 			return
	// 		}
	// 		if pty != nil {
	// 			or.PartnerId = pty.PartnerId
	// 			or.StoreId = pty.StoreId
	// 			break
	// 		}
	// 	}
	// 	if err = p.UpdateOrder(or, &product.Order{Id: or.Id}); err != nil {
	// 		log.Println("err: ", err, or.Id)
	// 		return
	// 	}

	// }
	log.Println("done")
}
