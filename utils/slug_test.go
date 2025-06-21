package utils_test

import (
	"log"
	"testing"

	"github.com/huyshop/header/product"
	"github.com/huyshop/product/db"
	"github.com/huyshop/product/utils"
)

func Test_slug(t *testing.T) {
	key := "Nồi điện đa năng"
	a := utils.ToSlug(key)
	log.Println("a: ", a)
}

func Test_convertSlug(t *testing.T) {
	p := &db.DB{}
	if err := p.ConnectDb("root:123456@tcp(localhost:3306)", "product"); err != nil {
		log.Println(err)
		return
	}
	listPty, err := p.ListProductType(&product.ProductTypeRequest{OrderBy: "id"})
	if err != nil {
		log.Println("err", err)
		return
	}
	for _, pty := range listPty {
		// if pty.Slug != "" {
		// 	continue
		// }
		pty.Slug = utils.ToSlug(pty.Name)
		if err := p.UpdateProductType(pty, &product.ProductType{Id: pty.Id}); err != nil {
			log.Println("err: ", err)
			return
		}
	}
	log.Println("done")
}
