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
	category := &pb.Category{}
	if _, err := d.engine.ID(id).Get(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (d *DB) listCategoryQuery(rq *pb.CategoryRequest) *xorm.Session {
	ss := d.engine.Table(tblCategory)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	}
	if rq.GetName() != "" {
		ss.Where("name like ?", "%"+rq.GetName()+"%")
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
	productType := &pb.ProductType{}
	if _, err := d.engine.ID(id).Get(productType); err != nil {
		return nil, err
	}
	return productType, nil
}

func (d *DB) listProductTypeQuery(rq *pb.ProductTypeRequest) *xorm.Session {
	ss := d.engine.Table(tblProductType)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	} else if rq.GetId() != "" {
		ss.Where("id = ?", rq.GetId())
	}
	if rq.GetName() != "" {
		ss.Where("name like ?", "%"+rq.GetName()+"%")
	}
	if rq.GetCategoryId() != "" {
		ss.Where("category_id = ?", rq.GetCategoryId())
	}
	if rq.GetBrand() != "" {
		ss.Where("brand = ?", rq.GetBrand())
	}
	if rq.GetOrigin() != "" {
		ss.Where("origin = ?", rq.GetOrigin())
	}
	return ss
}

func (d *DB) ListProductType(rq *pb.ProductTypeRequest) ([]*pb.ProductType, error) {
	productTypes := make([]*pb.ProductType, 0)
	ss := d.listProductTypeQuery(rq)
	if rq.GetLimit() > 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetLimit()*rq.GetSkip()))
	}
	if err := ss.Find(&productTypes); err != nil {
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
	product := &pb.Product{}
	if _, err := d.engine.ID(id).Get(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (d *DB) listProductQuery(rq *pb.ProductRequest) *xorm.Session {
	ss := d.engine.Table(tblProduct)
	if rq.GetIds() != nil {
		ss.In("id", rq.GetIds())
	} else if rq.GetId() != "" {
		ss.Where("id = ?", rq.GetId())
	}
	if rq.GetName() != "" {
		ss.Where("name like ?", "%"+rq.GetName()+"%")
	}
	if rq.GetState() != "" {
		ss.Where("state = ?", rq.GetState())
	}
	if len(rq.GetProductTypeIds()) > 0 {
		ss.In("product_type_id", rq.GetProductTypeIds())
	} else if rq.GetProductTypeId() != "" {
		ss.Where("product_type_id = ?", rq.GetProductTypeId())
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
