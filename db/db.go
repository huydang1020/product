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

func (d *DB) GetProductTypeBySlug(key string) (*pb.ProductType, error) {
	productType := &pb.ProductType{Slug: key}
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
	ss := d.engine.Table(tblProductType).Alias("pt").
		Join("INNER", []string{tblProduct, "p"}, "pt.id = p.product_type_id").
		Join("LEFT", []string{"review", "r"}, "r.product_id = p.id")

	// --- Điều kiện lọc trên product_type ---
	if rq.GetIds() != nil {
		ss.In("pt.id", rq.GetIds())
	} else if rq.GetId() != "" {
		ss.And("pt.id = ?", rq.GetId())
	}
	if rq.GetName() != "" {
		ss.And("pt.name LIKE ?", "%"+rq.GetName()+"%")
	}
	if rq.GetCategoryId() != "" {
		ss.And("pt.category_id = ?", rq.GetCategoryId())
	}
	if rq.GetState() != "" {
		ss.And("pt.state = ?", rq.GetState())
	}
	if len(rq.GetPartnerIds()) > 0 {
		ss.In("pt.partner_id", rq.GetPartnerIds())
	} else if rq.GetPartnerId() != "" {
		ss.And("pt.partner_id = ?", rq.GetPartnerId())
	}
	if rq.GetStoreId() != "" {
		ss.And("pt.store_id = ?", rq.GetStoreId())
	}
	if rq.GetQuantitySold() > 0 {
		ss.And("pt.quantity_sold >= ?", rq.GetQuantitySold())
	}
	if rq.GetViews() > 0 {
		ss.And("pt.quantity_search >= ?", rq.GetViews())
	}
	if rq.GetSlug() != "" {
		ss.And("pt.slug = ?", rq.GetSlug())
	}

	// --- Điều kiện lọc theo khoảng giá ---
	if rq.GetPriceFrom() > 0 {
		ss.And("p.sell_price >= ?", rq.GetPriceFrom())
	}
	if rq.GetPriceTo() > 0 {
		ss.And("p.sell_price <= ?", rq.GetPriceTo())
	}

	// --- Nhóm theo product_type ---
	ss.GroupBy("pt.id")

	// --- Lọc theo điểm đánh giá (HAVING) ---
	if rq.GetRatingFrom() > 0 {
		ss.Having(fmt.Sprintf("AVG(r.rating) >= %.2f", rq.GetRatingFrom()))

	}

	// --- Chọn các trường cần thiết ---
	selectCols := `
		pt.*,
		MIN(p.sell_price) AS min_price,
		MAX(p.sell_price) AS max_price,
		AVG(r.rating) AS average_rating,
		COUNT(r.id) AS total_reviews
	`
	ss.Select(selectCols)

	return ss
}

func (d *DB) ListProductType(rq *pb.ProductTypeRequest) ([]*pb.ProductType, error) {
	productTypes := make([]*pb.ProductType, 0)
	ss := d.listProductTypeQuery(rq)

	// Phân trang
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}

	// Sắp xếp
	switch rq.GetOrderBy() {
	case "price_asc":
		ss.OrderBy("min_price ASC")
	case "price_desc":
		ss.OrderBy("max_price DESC")
	case "rating":
		ss.OrderBy("average_rating DESC")
	case "sold":
		ss.OrderBy("pt.quantity_sold DESC")
	default:
		if rq.GetOrderBy() != "" {
			ss.OrderBy(rq.GetOrderBy()) // fallback nếu custom
		}
	}

	// Thực thi
	if err := ss.Find(&productTypes); err != nil {
		return nil, err
	}
	return productTypes, nil
}

func (d *DB) listProductTypeQueryOld(rq *pb.ProductTypeRequest) *xorm.Session {
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
	if rq.GetViews() > 0 {
		ss.And("quantity_search >= ?", rq.GetViews())
	}
	if rq.GetSlug() != "" {
		ss.And("slug = ?", rq.GetSlug())
	}
	return ss
}

func (d *DB) CountProductType(rq *pb.ProductTypeRequest) (int64, error) {
	return d.listProductTypeQueryOld(rq).Count()
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
	if rq.GetStoreIds() != nil {
		ss.In("store_id", rq.GetStoreIds())
	} else if rq.GetStoreId() != "" {
		ss.And("store_id = ?", rq.GetStoreId())
	}
	if rq.GetUserId() != "" {
		ss.And("user_id = ?", rq.GetUserId())
	}
	if rq.GetState() != "" {
		ss.And("state = ?", rq.GetState())
	}
	if rq.GetPartnerId() != "" {
		ss.And("partner_id = ?", rq.GetPartnerId())
	}
	return ss
}

func (d *DB) ListOrder(rq *pb.OrderRequest) ([]*pb.Order, error) {
	log.Println("rq: ", rq)
	orders := make([]*pb.Order, 0)
	ss := d.listOrderQuery(rq)
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}
	if err := ss.Desc("id").Find(&orders); err != nil {
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

	// create order details
	for _, odt := range order.OrderDetails {
		if _, err := sess.Insert(odt); err != nil {
			log.Print(err)
			sess.Rollback()
			return errors.New(utils.E_can_not_insert_order_detail)
		}
	}

	// create order ships
	for _, osh := range order.OrderShips {
		if _, err := sess.Insert(osh); err != nil {
			log.Print(err)
			sess.Rollback()
			return errors.New(utils.E_can_not_insert_order_ship)
		}
	}

	for _, orderItem := range order.GetProductOrdered() {
		// Create order detail
		pro, err := d.GetProduct(orderItem.GetProductId())
		if err != nil {
			log.Print(err)
			sess.Rollback()
			return errors.New(utils.E_not_found_product)
		}
		if pro.GetQuantity() < orderItem.GetQuantity() {
			log.Printf("Not enough quantity for product %s, available: %d, requested: %d", pro.GetId(), pro.GetQuantity(), orderItem.GetQuantity())
			sess.Rollback()
			return errors.New(utils.E_not_enough_quantity)
		}
		if order.GetMethodPayment() == "online" {
			// Update quantity sold.
			_, err = sess.Exec("UPDATE product_type SET quantity_sold = quantity_sold + ? WHERE id = ?", orderItem.GetQuantity(), pro.GetProductTypeId())
			if err != nil {
				log.Print(err)
				sess.Rollback()
				return errors.New(utils.E_can_not_update_product)
			}

			// Update the quantity sold for each product in the order.
			_, err = sess.Exec("UPDATE product SET quantity = quantity - ? WHERE id = ?", orderItem.GetQuantity(), orderItem.GetProductId())
			if err != nil {
				log.Print(err)
				sess.Rollback()
				return errors.New(utils.E_can_not_update_product)
			}
		}
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
	if _, err := sess.Update(order, &pb.Order{Id: order.Id}); err != nil {
		log.Print(err)
		sess.Rollback()
		return errors.New(utils.E_not_found)
	}
	if order.State == pb.Order_canceled.String() {
		if order.GetMethodPayment() == "online" {
			// If the order is canceled, we need to restore the product quantity.
			for _, item := range order.GetProductOrdered() {
				pro, err := d.GetProduct(item.GetProductId())
				if err != nil {
					log.Print(err)
					sess.Rollback()
					return errors.New(utils.E_not_found_product)
				}
				// Restore the product quantity.
				_, err = sess.Exec("UPDATE product SET quantity = quantity + ? WHERE id = ?", item.GetQuantity(), item.GetProductId())
				if err != nil {
					log.Print(err)
					sess.Rollback()
					return errors.New(utils.E_can_not_update_product)
				}
				// Restore the quantity sold.
				_, err = sess.Exec("UPDATE product_type SET quantity_sold = quantity_sold - ? WHERE id = ?", item.GetQuantity(), pro.GetProductTypeId())
				if err != nil {
					log.Print(err)
					sess.Rollback()
					return errors.New(utils.E_can_not_update_product)
				}
			}
		}
	}
	if order.State == pb.Order_completed.String() {
		if order.GetMethodPayment() == "cod" {
			// If the order is successful, we need to update the quantity sold.
			for _, item := range order.GetProductOrdered() {
				pro, err := d.GetProduct(item.GetProductId())
				if err != nil {
					log.Print(err)
					sess.Rollback()
					return errors.New(utils.E_not_found_product)
				}
				// Update the quantity sold.
				_, err = sess.Exec("UPDATE product_type SET quantity_sold = quantity_sold + ? WHERE id = ?", item.GetQuantity(), pro.GetProductTypeId())
				if err != nil {
					log.Print(err)
					sess.Rollback()
					return errors.New(utils.E_can_not_update_product)
				}
				// Update the quantity sold for each product in the order.
				_, err = sess.Exec("UPDATE product SET quantity = quantity - ? WHERE id = ?", item.GetQuantity(), item.GetProductId())
				if err != nil {
					log.Print(err)
					sess.Rollback()
					return errors.New(utils.E_can_not_update_product)
				}
			}
		}
	}
	return sess.Commit()
}

func (d *DB) CreateOrderDetail(req *pb.OrderDetail) error {
	c, err := d.engine.Insert(req)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_insert_order_detail)
	}
	return nil
}

func (d *DB) UpdateOrderDetail(updator, selector *pb.OrderDetail) error {
	c, err := d.engine.Update(updator, selector)
	if err != nil {
		return err
	}
	if c == 0 {
		log.Println("update order detail failed")
		return nil
	}
	return nil
}

func (d *DB) DeleteOrderDetail(req *pb.OrderDetail) error {
	c, err := d.engine.Delete(req)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_delete_order_detail)
	}
	return nil
}

func (d *DB) GetOrderDetail(req *pb.OrderDetail) (*pb.OrderDetail, error) {
	exist, err := d.engine.Get(req)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(utils.E_not_found_order_detail)
	}
	return req, nil
}

func (d *DB) listOrderDetailQuery(rq *pb.OrderDetailRequest) *xorm.Session {
	ss := d.engine.Table(tblOrderDetail)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	} else if rq.GetId() != "" {
		ss.And("id = ?", rq.GetId())
	}
	if rq.GetOrderId() != "" {
		ss.And("order_id = ?", rq.GetOrderId())
	}
	if rq.GetProductId() != "" {
		ss.And("product_id = ?", rq.GetProductId())
	}
	if rq.GetOrderShipId() != "" {
		ss.And("order_ship_id = ?", rq.GetOrderShipId())
	}
	return ss
}

func (d *DB) ListOrderDetail(rq *pb.OrderDetailRequest) ([]*pb.OrderDetail, error) {
	orderDetails := make([]*pb.OrderDetail, 0)
	ss := d.listOrderDetailQuery(rq)
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}
	if err := ss.Desc("id").Find(&orderDetails); err != nil {
		return nil, err
	}
	return orderDetails, nil
}

func (d *DB) CountOrderDetail(rq *pb.OrderDetailRequest) (int64, error) {
	return d.listOrderDetailQuery(rq).Count()
}

func (d *DB) CreateOrderShip(orderShip *pb.OrderShip) error {
	c, err := d.engine.Insert(orderShip)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_insert_order_ship)
	}
	return nil
}

func (d *DB) UpdateOrderShip(updator, selector *pb.OrderShip) error {
	c, err := d.engine.Update(updator, selector)
	if err != nil {
		return err
	}
	if c == 0 {
		log.Println("update order ship failed")
		return nil
	}
	return nil
}

func (d *DB) DeleteOrderShip(orderShip *pb.OrderShip) error {
	c, err := d.engine.ID(orderShip.Id).Delete(orderShip)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_delete_order_ship)
	}
	return nil
}

func (d *DB) GetOrderShip(id string) (*pb.OrderShip, error) {
	orderShip := &pb.OrderShip{Id: id}
	exist, err := d.engine.Get(orderShip)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(utils.E_not_found_order_ship)
	}
	return orderShip, nil
}

func (d *DB) listOrderShipQuery(rq *pb.OrderShipRequest) *xorm.Session {
	ss := d.engine.Table(tblOrderShip)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	} else if rq.GetId() != "" {
		ss.And("id = ?", rq.GetId())
	}
	if rq.GetOrderId() != "" {
		ss.And("order_id = ?", rq.GetOrderId())
	}
	if rq.GetState() != "" {
		ss.And("state = ?", rq.GetState())
	}
	if rq.GetProvince() != "" {
		ss.And("province = ?", rq.GetProvince())
	}
	if rq.GetDistrict() != "" {
		ss.And("district = ?", rq.GetDistrict())
	}
	if rq.GetWard() != "" {
		ss.And("ward = ?", rq.GetWard())
	}
	return ss
}

func (d *DB) ListOrderShip(rq *pb.OrderShipRequest) ([]*pb.OrderShip, error) {
	orderShips := make([]*pb.OrderShip, 0)
	ss := d.listOrderShipQuery(rq)
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}
	if err := ss.Desc("id").Find(&orderShips); err != nil {
		return nil, err
	}
	return orderShips, nil
}

func (d *DB) CountOrderShip(rq *pb.OrderShipRequest) (int64, error) {
	return d.listOrderShipQuery(rq).Count()
}

func (d *DB) CreateReview(review *pb.Review) error {
	c, err := d.engine.Insert(review)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_insert_review)
	}
	return nil
}

func (d *DB) UpdateReview(updator, selector *pb.Review) error {
	c, err := d.engine.Update(updator, selector)
	if err != nil {
		return err
	}
	if c == 0 {
		log.Println("update review failed")
		return nil
	}
	return nil
}

func (d *DB) DeleteReview(review *pb.Review) error {
	c, err := d.engine.ID(review.Id).Delete(review)
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New(utils.E_can_not_delete_review)
	}
	return nil
}

func (d *DB) GetReview(review *pb.Review) (*pb.Review, error) {
	exist, err := d.engine.Get(review)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New(utils.E_not_found_review)
	}
	return review, nil
}

func (d *DB) IsReviewExist(review *pb.Review) bool {
	any, err := d.engine.Exist(review)
	if err != nil {
		return false
	}
	return any
}

func (d *DB) listReviewQuery(rq *pb.ReviewRequest) *xorm.Session {
	ss := d.engine.Table(tblReview)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	} else if rq.GetId() != "" {
		ss.And("id = ?", rq.GetId())
	}
	if len(rq.GetProductIds()) > 0 {
		ss.In("product_id", rq.GetProductId())
	} else if rq.GetProductId() != "" {
		ss.And("product_id = ?", rq.GetProductId())
	}
	if len(rq.GetOrderIds()) > 0 {
		ss.In("order_id", rq.GetOrderIds())
	} else if rq.GetOrderId() != "" {
		ss.And("order_id = ?", rq.GetOrderId())
	}
	if rq.GetUserId() != "" {
		ss.And("user_id = ?", rq.GetUserId())
	}
	if rq.GetRating() != 0 {
		ss.And("rating = ?", rq.GetRating())
	}
	return ss
}

func (d *DB) ListReview(rq *pb.ReviewRequest) ([]*pb.Review, error) {
	reviews := make([]*pb.Review, 0)
	ss := d.listReviewQuery(rq)
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}
	if err := ss.Desc("id").Find(&reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (d *DB) CountReview(rq *pb.ReviewRequest) (int64, error) {
	return d.listReviewQuery(rq).Count()
}
