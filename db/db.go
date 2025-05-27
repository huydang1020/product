package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/huyshop/header/product"
	"github.com/huyshop/product/utils"
	"xorm.io/xorm"
)

type DB struct {
	engine *xorm.Engine
}

func (d *DB) ConnectDb(sqlPath, dbName string) error {
	sqlConnStr := fmt.Sprintf("%s/%s", sqlPath, dbName)
	engine, err := xorm.NewEngine("mysql", sqlConnStr)
	if err != nil {
		return err
	}
	tickPingSql := time.NewTicker(15 * time.Minute)
	go func() {
		for {
			select {
			case <-tickPingSql.C:
				if err := engine.Ping(); err != nil {
					log.Print("sql can not ping")
				}
			}
		}
	}()
	d.engine = engine
	d.engine.ShowSQL(false)
	return err
}

func (d *DB) CreateCategory(category *pb.Category) error {
	c, err := d.engine.Insert(category)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_insert)
	}
	return nil
}

func (d *DB) UpdateCategory(updator, selector *pb.Category) error {
	c, err := d.engine.Update(updator, selector)
	if err != nil {
		return err
	}
	if c == 0 {
		log.Println("update category failed")
		return nil
	}
	return nil
}

func (d *DB) DeleteCategory(category *pb.Category) error {
	c, err := d.engine.ID(category.Id).Delete(category)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_delete)
	}
	return nil
}

func (d *DB) GetCategory(id string) (*pb.Category, error) {
	category := &pb.Category{Id: id}
	exist, err := d.engine.Get(category)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(utils.E_not_found_category)
	}
	return category, nil
}

func (d *DB) listCategoryQuery(rq *pb.CategoryRequest) *xorm.Session {
	ss := d.engine.Table(tblCategory)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	}
	if rq.GetName() != "" {
		ss.And("name like ?", "%"+rq.GetName()+"%")
	}
	if rq.GetState() != "" {
		ss.And("state = ?", rq.GetState())
	}
	return ss
}

func (d *DB) ListCategory(rq *pb.CategoryRequest) ([]*pb.Category, error) {
	categories := make([]*pb.Category, 0)
	ss := d.listCategoryQuery(rq)
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}
	if err := ss.Find(&categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (d *DB) CountCategory(rq *pb.CategoryRequest) (int64, error) {
	return d.listCategoryQuery(rq).Count()
}

func (d *DB) CreateProductType(productType *pb.ProductType) error {
	c, err := d.engine.Insert(productType)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_insert)
	}
	return nil
}

func (d *DB) UpdateProductType(updator, selector *pb.ProductType) error {
	c, err := d.engine.Update(updator, selector)
	if err != nil {
		return err
	}
	if c == 0 {
		log.Println("update product type failed")
		return nil
	}
	return nil
}

func (d *DB) DeleteProductType(id string) error {
	productType := &pb.ProductType{}
	c, err := d.engine.ID(id).Delete(productType)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_delete)
	}
	return nil
}

func (d *DB) GetProductType(id string) (*pb.ProductType, error) {
	productType := &pb.ProductType{Id: id}
	exist, err := d.engine.Get(productType)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(utils.E_not_found_product_type)
	}
	return productType, nil
}

func (d *DB) listProductTypeQuery(rq *pb.ProductTypeRequest) *xorm.Session {
	ss := d.engine.Table(tblProductType)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	} else if rq.GetId() != "" {
		ss.And("id = ?", rq.GetId())
	}
	if rq.GetName() != "" {
		ss.And("name like ?", "%"+rq.GetName()+"%")
	}
	if rq.GetCategoryId() != "" {
		ss.And("category_id = ?", rq.GetCategoryId())
	}
	if rq.GetState() != "" {
		ss.And("state = ?", rq.GetState())
	}
	if len(rq.GetPartnerIds()) > 0 {
		ss.In("partner_id", rq.GetPartnerIds())
	} else if rq.GetPartnerId() != "" {
		ss.And("partner_id = ?", rq.GetPartnerId())
	}
	if rq.GetStoreId() != "" {
		ss.And("store_id = ?", rq.GetStoreId())
	}
	if rq.GetQuantitySold() > 0 {
		ss.And("quantity_sold >= ?", rq.GetQuantitySold())
	}
	if rq.GetQuantitySearch() > 0 {
		ss.And("quantity_search >= ?", rq.GetQuantitySearch())
	}
	return ss
}

func (d *DB) ListProductType(rq *pb.ProductTypeRequest) ([]*pb.ProductType, error) {
	productTypes := make([]*pb.ProductType, 0)
	ss := d.listProductTypeQuery(rq)
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}
	if err := ss.Desc(rq.GetOrderBy()).Find(&productTypes); err != nil {
		return nil, err
	}
	return productTypes, nil
}

func (d *DB) CountProductType(rq *pb.ProductTypeRequest) (int64, error) {
	return d.listProductTypeQuery(rq).Count()
}

func (d *DB) CreateProduct(product *pb.Product) error {
	c, err := d.engine.Insert(product)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_insert)
	}
	return nil
}

func (d *DB) UpdateProduct(updator, selector *pb.Product) error {
	c, err := d.engine.Update(updator, selector)
	if err != nil {
		return err
	}
	if c == 0 {
		log.Println("update product failed")
		return nil
	}
	return nil
}

func (d *DB) DeleteProduct(product *pb.Product) error {
	c, err := d.engine.ID(product.Id).Delete(product)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_delete)
	}
	return nil
}

func (d *DB) GetProduct(id string) (*pb.Product, error) {
	product := &pb.Product{Id: id}
	exist, err := d.engine.Get(product)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(utils.E_not_found_product)
	}
	return product, nil
}

func (d *DB) listProductQuery(rq *pb.ProductRequest) *xorm.Session {
	ss := d.engine.Table(tblProduct)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	} else if rq.GetId() != "" {
		ss.And("id = ?", rq.GetId())
	}
	if rq.GetName() != "" {
		ss.And("name like ?", "%"+rq.GetName()+"%")
	}
	if rq.GetState() != "" {
		ss.And("state = ?", rq.GetState())
	}
	if len(rq.GetProductTypeIds()) > 0 {
		ss.In("product_type_id", rq.GetProductTypeIds())
	} else if rq.GetProductTypeId() != "" {
		ss.And("product_type_id = ?", rq.GetProductTypeId())
	}
	return ss
}

func (d *DB) ListProduct(rq *pb.ProductRequest) ([]*pb.Product, error) {
	products := make([]*pb.Product, 0)
	ss := d.listProductQuery(rq)
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}
	if err := ss.Find(&products); err != nil {
		return nil, err
	}
	return products, nil
}

func (d *DB) CountProduct(rq *pb.ProductRequest) (int64, error) {
	return d.listProductQuery(rq).Count()
}

// Create Product
func (d *DB) TransCreateProductType(pt *pb.ProductType) error {
	sess := d.engine.NewSession()
	defer sess.Close()
	// Start transcation.
	if err := sess.Begin(); err != nil {
		return err
	}

	if _, err := sess.Insert(pt); err != nil {
		log.Print(err)
		sess.Commit()
		return errors.New(utils.E_not_found)
	}
	has, err := sess.Get(pt)
	if err != nil {
		log.Print(err)
		sess.Rollback()
		return err
	}
	if !has {
		sess.Rollback()
		return errors.New(utils.E_can_not_insert)
	}

	for _, pro := range pt.Products {
		if pro.GetOriginPrice() < 0 || pro.GetSellPrice() < 0 {
			return errors.New(utils.E_invalid_price)
		}
		if _, err := sess.Insert(pro); err != nil {
			log.Print(err)
			sess.Rollback()
			return errors.New(utils.E_can_not_insert_product)
		}
	}

	return sess.Commit()
}

// update product
func (d *DB) TransUpdateProductType(in *pb.ProductType) error {
	sess := d.engine.NewSession()
	defer sess.Close()
	// Start transcation.
	if err := sess.Begin(); err != nil {
		return err
	}

	count, err := sess.Update(in, &pb.ProductType{Id: in.Id})
	if err != nil {
		log.Print(err)
		sess.Rollback()
		return err
	}
	if count < 1 {
		sess.Rollback()
		return errors.New(utils.E_can_not_update_product_type)
	}
	oldProducts := make([]*pb.Product, 0)

	if err := sess.Table(tblProduct).
		Where("product_type_id = ?", in.Id).
		Find(&oldProducts); err != nil {
		sess.Rollback()
		return err
	}
	mProducts := map[string]*pb.Product{}
	for _, item := range oldProducts {
		mProducts[item.Id] = item
	}

	for _, pro := range in.GetProducts() {
		if pro.GetState() == "" {
			pro.State = in.GetState()
		}
		if pro.GetOriginPrice() < 0 || pro.GetSellPrice() < 0 {
			return errors.New(utils.E_invalid_price)
		}
		_, has := mProducts[pro.Id]
		if has {
			if _, err := sess.Update(pro, &pb.Product{Id: pro.GetId()}); err != nil {
				log.Print(err)
				sess.Rollback()
				return err
			}
			delete(mProducts, pro.GetId())
			continue
		}
		if pro.GetId() == "" {
			pro.Id = utils.MakeProductId()
			pro.CreatedAt = time.Now().Unix()
			pro.ProductTypeId = in.GetId()
		}
		if _, err := sess.Insert(pro); err != nil {
			log.Print(err)
			sess.Rollback()
			return err
		}
		continue
	}
	if len(mProducts) > 0 {
		for id, _ := range mProducts {
			count, err = sess.Where("id = ?", id).Delete(&pb.Product{})
			if err != nil {
				log.Print(err)
				sess.Rollback()
				return err
			}
			if count < 1 {
				sess.Rollback()
				return errors.New(utils.E_can_not_delete_product)
			}
		}
	}
	return sess.Commit()
}

// update state product
func (d *DB) TransUpdateStateProductType(pt *pb.ProductType) error {
	sess := d.engine.NewSession()
	defer sess.Close()
	// Start transcation.
	if err := sess.Begin(); err != nil {
		return err
	}

	count, err := sess.Where("id = ?", pt.Id).Update(&pb.ProductType{
		State:     pt.State,
		UpdatedAt: pt.UpdatedAt,
	})
	if err != nil {
		log.Print(err)
		sess.Rollback()
		return err
	}
	if count < 1 {
		sess.Rollback()
		return errors.New(utils.E_can_not_update_product_type)
	}

	listProduct := []*pb.Product{}

	err = sess.Where("product_type_id = ?", pt.Id).Find(&listProduct)
	if err != nil {
		log.Print(err)
		sess.Rollback()
		return errors.New(utils.E_not_found_product)
	}

	for _, pr := range listProduct {
		count, err := sess.Where("id = ?", pr.Id).Update(&pb.Product{
			State:     pt.State,
			UpdatedAt: pt.UpdatedAt,
		})
		if err != nil {
			log.Print(err)
			sess.Rollback()
			return err
		}
		if count < 1 {
			sess.Rollback()
			return errors.New(utils.E_can_not_update_product)
		}
	}

	return sess.Commit()
}

func (d *DB) TransDeleteProductType(ptid string) error {
	sess := d.engine.NewSession()
	defer sess.Close()
	// Start transcation.
	if err := sess.Begin(); err != nil {
		return err
	}

	count, err := sess.Where("id = ?", ptid).Delete(&pb.ProductType{})
	if err != nil {
		log.Print(err)
		sess.Rollback()
		return err
	}
	if count < 1 {
		sess.Rollback()
		return errors.New(utils.E_can_not_delete_product_type)
	}

	count, err = sess.Where("product_type_id = ?", ptid).Delete(&pb.Product{})
	if err != nil {
		log.Print(err)
		sess.Rollback()
		return err
	}
	if count < 1 {
		sess.Rollback()
		return errors.New(utils.E_can_not_delete_product)
	}

	return sess.Commit()
}

func (d *DB) CreateBanner(banner *pb.Banner) error {
	c, err := d.engine.Insert(banner)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_insert)
	}
	return nil
}

func (d *DB) UpdateBanner(updator, selector *pb.Banner) error {
	c, err := d.engine.Update(updator, selector)
	if err != nil {
		return err
	}
	if c == 0 {
		log.Println("update banner failed")
		return nil
	}
	return nil
}

func (d *DB) DeleteBanner(banner *pb.Banner) error {
	c, err := d.engine.ID(banner.Id).Delete(banner)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_delete)
	}
	return nil
}

func (d *DB) GetBanner(id string) (*pb.Banner, error) {
	banner := &pb.Banner{Id: id}
	exist, err := d.engine.Get(banner)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(utils.E_not_found_banner)
	}
	return banner, nil
}

func (d *DB) listBannerQuery(rq *pb.BannerRequest) *xorm.Session {
	ss := d.engine.Table(tblBanner)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	}
	if rq.GetName() != "" {
		ss.And("name like ?", "%"+rq.GetName()+"%")
	}
	if rq.GetState() != "" {
		ss.And("state = ?", rq.GetState())
	}
	if rq.GetType() != "" {
		ss.And("type = ?", rq.GetType())
	}
	return ss
}

func (d *DB) ListBanner(rq *pb.BannerRequest) ([]*pb.Banner, error) {
	categories := make([]*pb.Banner, 0)
	ss := d.listBannerQuery(rq)
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}
	if err := ss.Find(&categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (d *DB) CountBanner(rq *pb.BannerRequest) (int64, error) {
	return d.listBannerQuery(rq).Count()
}

func (d *DB) CreateOrder(order *pb.Order) error {
	c, err := d.engine.Insert(order)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_insert)
	}
	return nil
}
func (d *DB) UpdateOrder(updator, selector *pb.Order) error {
	c, err := d.engine.Update(updator, selector)
	if err != nil {
		return err
	}
	if c == 0 {
		log.Println("update order failed")
		return nil
	}
	return nil
}
func (d *DB) DeleteOrder(order *pb.Order) error {
	c, err := d.engine.ID(order.Id).Delete(order)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_delete)
	}
	return nil
}
func (d *DB) GetOrder(id string) (*pb.Order, error) {
	order := &pb.Order{Id: id}
	exist, err := d.engine.Get(order)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(utils.E_not_found_order)
	}
	return order, nil
}

func (d *DB) listOrderQuery(rq *pb.OrderRequest) *xorm.Session {
	ss := d.engine.Table(tblOrder)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	} else if rq.GetId() != "" {
		ss.And("id = ?", rq.GetId())
	}
	if rq.GetUserId() != "" {
		ss.And("user_id = ?", rq.GetUserId())
	}
	if rq.GetState() != "" {
		ss.And("state = ?", rq.GetState())
	}
	return ss
}

func (d *DB) ListOrder(rq *pb.OrderRequest) ([]*pb.Order, error) {
	orders := make([]*pb.Order, 0)
	ss := d.listOrderQuery(rq)
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}
	if err := ss.Find(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (d *DB) CountOrder(rq *pb.OrderRequest) (int64, error) {
	return d.listOrderQuery(rq).Count()
}

func (d *DB) TransCreateOrder(order *pb.Order) error {
	sess := d.engine.NewSession()
	defer sess.Close()
	// Start transaction.
	if err := sess.Begin(); err != nil {
		return err
	}

	// Insert the order.
	if _, err := sess.Insert(order); err != nil {
		log.Print(err)
		sess.Rollback()
		return errors.New(utils.E_not_found)
	}
	log.Println("done")
	return sess.Commit()
}

// update quantity of product in order
func (d *DB) TransUpdateOrder(order *pb.Order) error {
	sess := d.engine.NewSession()
	defer sess.Close()
	// Start transaction.
	if err := sess.Begin(); err != nil {
		return err
	}

	// Update the quantity sold for each product in the order.
	for _, orderItem := range order.GetOrderDetail() {
		prod, err := d.GetProduct(orderItem.GetProductId())
		if err != nil {
			log.Print(err)
			sess.Rollback()
			return errors.New(utils.E_not_found_product)
		}
		productType, err := d.GetProductType(prod.GetProductTypeId())
		if err != nil {
			log.Print(err)
			sess.Rollback()
			return errors.New(utils.E_not_found_product_type)
		}
		// Update quantity sold.
		_, err = sess.Exec("UPDATE product_type SET quantity_sold = quantity_sold + ? WHERE id = ?", orderItem.GetQuantity(), productType.GetId())
		if err != nil {
			log.Print(err)
			sess.Rollback()
			return errors.New(utils.E_can_not_update_product)
		}

		// Decrease available quantity.
		_, err = sess.Exec("UPDATE product SET quantity = quantity - ? WHERE id = ?", orderItem.GetQuantity(), orderItem.GetProductId())
		if err != nil {
			log.Print(err)
			sess.Rollback()
			return errors.New(utils.E_can_not_update_product)
		}
	}
	return sess.Commit()
}
