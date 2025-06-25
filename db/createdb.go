package db

import (
	"log"

	pb "github.com/huyshop/header/product"
	"xorm.io/xorm"
)

const (
	tblProductType = "product_type"
	tblProduct     = "product"
	tblCategory    = "category"
	tblBanner      = "banner"
	tblOrder       = "order"
	tblOrderShip   = "order_ship"
	tblReview      = "review"
)

func createTable(model interface{}, tblName string, engine *xorm.Engine) error {
	log.Println("createTable", tblName)
	b, err := engine.IsTableExist(model)
	if err != nil {
		return err
	}
	log.Print(b, " ", tblName)
	if b {
		if err = engine.Sync2(model); err != nil {
			return err
		}
		return nil
	}
	if !b {
		if err := engine.CreateTables(model); err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func (d *DB) CreateDb() error {
	if err := createTable(&pb.Product{}, tblProduct, d.engine); err != nil {
		return err
	}
	if err := createTable(&pb.ProductType{}, tblProductType, d.engine); err != nil {
		return err
	}
	if err := createTable(&pb.Category{}, tblCategory, d.engine); err != nil {
		return err
	}
	if err := createTable(&pb.Banner{}, tblBanner, d.engine); err != nil {
		return err
	}
	if err := createTable(&pb.Order{}, tblOrder, d.engine); err != nil {
		return err
	}
	if err := createTable(&pb.OrderShip{}, tblOrderShip, d.engine); err != nil {
		return err
	}
	if err := createTable(&pb.Review{}, tblReview, d.engine); err != nil {
		return err
	}
	return nil
}
