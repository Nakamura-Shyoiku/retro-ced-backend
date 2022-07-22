package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

type Product struct {
	Guid             string    `json:"guid" gorm:"column:guid"`
	SiteId           int64     `json:"site_id" gorm:"column:site_id"`
	UrlId            int64     `json:"url_id" gorm:"column:url_id"` // Changed column tag from '-' to 'url_id'
	Url              string    `json:"url" gorm:"column:url"`
	Category         string    `json:"category" gorm:"column:category"`
	Brand            string    `json:"brand" gorm:"column:brand"`
	Title            string    `json:"title" gorm:"column:title"`
	Description      string    `json:"description" gorm:"column:description"`
	Price            int64     `json:"price" gorm:"column:price"`
	RetailPrice      int64     `db:"retail_price" json:"retail_price" gorm:"column:retail_price"`
	Model            string    `json:"model" gorm:"column:model"`
	ItemNumber       string    `db:"item_number" json:"item_number" gorm:"column:item_number"`
	ProductCondition string    `db:"product_condition" json:"product_condition" gorm:"column:product_condition"`
	Accessories      string    `json:"accessories" gorm:"column:accessories"`
	Measurements     string    `json:"measurements" gorm:"column:measurements"`
	Featured         string    `json:"featured" gorm:"column:featured"`
	Img              string    `json:"img" gorm:"column:img"`
	Approved         bool      `json:"approved" gorm:"column:approved"`
	Color            string    `json:"color" gorm:"column:color"`
	Size             string    `json:"size" gorm:"column:size"`
	ShoeSize         string    `db:"shoe_size" json:"shoe_size" gorm:"column:shoe_size"`
	SubCategory      string    `db:"sub_category" gorm:"column:sub_category"`
	ProductURL       string    `db:"product_url" gorm:"column:product_url"`
	UpdatedAt        time.Time `db:"last_updated" json:"last_updated" gorm:"column:last_updated"`
	CreatedAt        time.Time `db:"created_at" json:"created_at" gorm:"column:created_at"`
	//	IsFavourited     bool      `json:"is_favourited" gorm:"column:is_favourited"`
}

// ProductRecord model.
// NOTE: this is a product record, as it should be, it ONLY refers to the product table - no outside dependencies.
// It is also duplicated table from above, which would require cleaning all over the codebase.
type ProductRecord struct {
	Guid             string    `json:"guid" gorm:"column:guid"`
	Id               int64     `json:"Id" gorm:"column:id;primaryKey"`
	SiteId           int64     `json:"site_id" gorm:"column:site_id"`
	UrlId            int64     `json:"url_id" gorm:"column:url_id"`
	Category         string    `json:"category" gorm:"column:category"`
	Brand            string    `json:"brand" gorm:"column:brand"`
	Title            string    `json:"title" gorm:"column:title"`
	Description      string    `json:"description" gorm:"column:description"`
	Price            int64     `json:"price" gorm:"column:price"`
	RetailPrice      int64     `json:"retail_price" gorm:"column:retail_price"`
	Model            string    `json:"model" gorm:"column:model"`
	ItemNumber       string    `json:"item_number" gorm:"column:item_number"`
	ProductCondition string    `json:"product_condition" gorm:"column:product_condition"`
	Accessories      string    `json:"accessories" gorm:"column:accessories"`
	Measurements     string    `json:"measurements" gorm:"column:measurements"`
	Featured         string    `json:"featured" gorm:"column:featured"`
	Img              string    `json:"img" gorm:"column:img"`
	Approved         bool      `json:"approved" gorm:"column:approved"`
	Color            string    `json:"color" gorm:"column:color"`
	Size             string    `json:"size" gorm:"column:size"`
	ShoeSize         string    `json:"shoe_size" gorm:"column:shoe_size"`
	SubCategory      string    `json:"sub_category" gorm:"column:sub_category"`
	ProductURL       *string   `json:"product_url" gorm:"column:product_url"` // TODO: if used for other crawlers, it should be value or NULL
	UpdatedAt        time.Time `json:"last_updated" gorm:"column:last_updated"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at"`
}

func (p ProductRecord) TableName() string {
	return "Products"
}

func (p *Product) GetProduct(guid string) (*Product, error) {
	var prod Product
	err := GetDBv2().
		Table("Products").
		Where("guid = ?", guid).
		Scan(&prod).Error
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("get product error")
		return nil, err
	}
	return &prod, nil
}

func (p *Product) SetFeaturedProducts(guid string, featured string) error {
	err := GetDBv2().Table("Products").Exec(`
		ALTER TABLE Products
		UPDATE featured = ?,
					last_updated = ?
		WHERE guid = ?`, featured, time.Now(), guid).Error

	if err != nil {
		log.WithError(err).Error("failed prepare update user statement")
		return err
	}

	return nil
}

// GetProductsCount returns the number of products matching the specific condition. Used for /admin/products listing.
func GetProductsCount(search string, category string, siteId string, guid string) (count uint64, err error) {

	searchQuery := "%" + search + "%"

	if siteId == "0" {
		siteId = ""
	}

	err = GetDB().QueryRow(
		`SELECT COUNT(*) FROM (
			SELECT * FROM Products 
			WHERE Products.category = IF(? = '',category, ?)
			AND domain(Products.product_url) = IF(? = '', domain(Products.product_url), ?)
			AND Products.guid = IF(? = '', guid, ?)
		) AS Product
		WHERE Product.title ILIKE ?
		OR Product.brand ILIKE ?
		`,
		category,
		category,
		siteId,
		siteId,
		guid,
		guid,
		searchQuery,
		searchQuery,
	).Scan(
		&count,
	)
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("products count row error")
		return count, err
	}
	return count, nil
}

func GetProductsCountByCategory(category string, filterCategory []string, filterBrand []string, filterColor []string, filterSize []string, filterShoeSize []string) (count uint64, err error) {
	categoryString := strings.Join(filterCategory, ",")
	brandString := strings.Join(filterBrand, ",")
	colorString := strings.Join(filterColor, ",")
	sizeString := strings.Join(filterSize, ",")
	shoeSizeString := strings.Join(filterShoeSize, ",")
	err = GetDB().QueryRow(
		`SELECT COUNT(*) FROM Products
		WHERE Products.category = ?
		AND Products.approved = 1
		AND ((FIND_IN_SET(Products.sub_category, ?) > 0) OR ? = "")
		AND ((FIND_IN_SET(Products.brand, ?) > 0) OR ? = "")
		AND ((FIND_IN_SET(Products.color, ?) > 0) OR ? = "")
		AND ((FIND_IN_SET(Products.size, ?) > 0) OR ? = "")
		AND ((FIND_IN_SET(Products.shoe_size, ?) > 0) OR ? = "")`,
		category,
		categoryString,
		categoryString,
		brandString,
		brandString,
		colorString,
		colorString,
		sizeString,
		sizeString,
		shoeSizeString,
		shoeSizeString,
	).Scan(
		&count,
	)
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("products count row error")
		return count, err
	}
	return count, nil
}

func GetProductsCountByBrand(brand string) (count uint64, err error) {
	err = GetDB().QueryRow(
		`SELECT COUNT(*) FROM Products
		WHERE brand = ?
		AND approved = 1`,
		brand,
	).Scan(
		&count,
	)

	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("products count row error")
		return count, err
	}
	return count, nil
}

func GetProductsCountBySearch(search string) (count uint64, err error) {
	searchQuery := "%" + search + "%"
	err = GetDB().QueryRow(
		`SELECT COUNT(*) FROM Products
		WHERE approved = 1 and  title LIKE ?
		OR brand LIKE ?`,
		searchQuery,
		searchQuery,
	).Scan(
		&count,
	)
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("products count row error")
		return count, err
	}
	return count, nil
}

// GetProducts returns the list of products according to certain parameters. This is used by the /admin/products listing at the moment.
func GetProducts(search string,
	offset int,
	itemsPerPage int,
	category string,
	siteID string,
	guid string,
	sortBy int) ([]Product, error) {

	var ps []Product
	db := GetDBv2().
		Table("Products").
		Offset(int(offset)).
		Limit(int(itemsPerPage)).
		Order("created_at desc")

	if siteID == "0" {
		siteID = ""
	}

	if guid != "" {
		db = db.Where("guid = ?", guid)
	}
	if category != "" {
		db = db.Where("category = ?", category)
	}
	if siteID != "" {
		db = db.Where("domain(product_url) = ?", siteID)
	}
	if search != "" {
		searchQuery := "%" + search + "%"
		db = db.Where("title ilike ? or brand ilike ?", searchQuery, searchQuery)
	}

	err := db.Scan(&ps).Error
	if err != nil {
		log.WithError(err).Error("failed to get products")
		return nil, err
	}
	return ps, nil
}

// Create inserts row
func (p *Product) Create() (err error) {
	p.Guid = uuid.New().String()
	return GetDBv2().Table("Products").Create(p).Error
}

// GetProductsByCategory function returns a list of products. This model is tightly coupled to the frontend
// so double check before changing the model here.
func GetProductsByCategory(
	category string,
	offset string,
	filterCategory []string,
	filterBrand []string,
	filterColor []string,
	filterSize []string,
	filterShoeSize []string,
	userId int64) ([]Product, error) {
	categoryString := strings.Join(filterCategory, ",")
	brandString := strings.Join(filterBrand, ",")
	colorString := strings.Join(filterColor, ",")
	sizeString := strings.Join(filterSize, ",")
	shoeSizeString := strings.Join(filterShoeSize, ",")
	ps := make([]Product, 0)
	rows, err := GetDB().Query(
		`SELECT
		Product.guid,
		Product.site_id,
		Product.url_id,
		Product.url,
		Product.category,
		Product.brand,
		Product.title,
		Product.description,
		Product.price,
		Product.img,
		MAX(IF(Favourites.user_id = ?, TRUE, FALSE)) is_favourited
			FROM (SELECT
			Products.id,
			Products.site_id,
			Products.url_id,
			COALESCE(Urls.url, Products.product_url) AS url,
			Products.category,
			Products.brand,
			Products.title,
			Products.description, 
			Products.price,
			Products.img,
			Products.approved
			FROM Products INNER JOIN Urls
			WHERE Products.url_id = Urls.id
			AND Products.category = ?
			AND Products.approved = 1
			AND ((FIND_IN_SET(Products.sub_category, ?) > 0) OR ? = "")
			AND ((FIND_IN_SET(Products.brand, ?) > 0) OR ? = "")
			AND ((FIND_IN_SET(Products.color, ?) > 0) OR ? = "")
			AND ((FIND_IN_SET(Products.size, ?) > 0) OR ? = "")
			AND ((FIND_IN_SET(Products.shoe_size, ?) > 0) OR ? = "")
			LIMIT 18 OFFSET ?
			) AS Product
		LEFT OUTER JOIN Favourites
		ON Product.id = Favourites.product_id
		WHERE Product.approved = 1
		GROUP BY Product.id`,
		userId,
		category,
		categoryString,
		categoryString,
		brandString,
		brandString,
		colorString,
		colorString,
		sizeString,
		sizeString,
		shoeSizeString,
		shoeSizeString,
		offset,
	)
	if err != nil {
		log.WithError(err).Error("failed to query products by category")
	}
	defer rows.Close()
	for rows.Next() {
		var p Product
		err := rows.Scan(
			&p.Guid,
			&p.SiteId,
			&p.UrlId,
			&p.Url,
			&p.Category,
			&p.Brand,
			&p.Title,
			&p.Description,
			&p.Price,
			&p.Img,
			//			&p.IsFavourited,
		)
		if err != nil {
			log.WithError(err).Error("failed to get products by category rows")
			return ps, err
		}
		ps = append(ps, p)
	}
	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("products by category row error")
		return ps, err
	}

	return ps, nil
}

// GetProductsByFeatured will return featured products according to the category.
func GetProductsByFeatured(category string, fetchSize int) ([]Product, error) {

	var products []Product

	db := GetDBv2().Table("Products").Where("featured LIKE ?", category).Order("last_updated DESC").Limit(fetchSize)
	err := db.Scan(&products).Error
	if err != nil {
		log.WithError(err).Error("failed to get products by featured")
		return products, err
	}
	//selectColumns := []string{
	//	"p.id",
	//	"p.site_id",
	//	"p.url_id",
	//	"COALESCE(u.url, p.product_url) AS url",
	//	"p.category",
	//	"p.brand",
	//	"p.title",
	//	"p.description",
	//	"p.featured",
	//	"p.img",
	//	"p.price",
	//	"p.measurements",
	//	"MAX(IF(f.user_id = ?, 1, 0)) is_favourited",
	//}
	//
	//selectStatement := strings.Join(selectColumns, ",")
	//
	//err := GetDBv2().Table("Products p").
	//	Select(selectStatement, userID).
	//	Joins("LEFT JOIN Urls u ON p.url_id = u.id").
	//	Joins("LEFT OUTER JOIN Favourites f ON p.id = f.product_id").
	//	Where("p.featured = ?", category).
	//	Group("p.id").
	//	Order("p.id DESC").
	//	Find(&products).
	//	Error
	//if err != nil {
	//	return nil, fmt.Errorf( "could not lookup favourite products %",err)
	//}

	//return products, nil
	return products, nil
}

func GetProductsByBrand(brand string, offset string, userId int64) ([]Product, error) {
	var ps []Product
	err := GetDBv2().Table("Products").Find(&ps, "brand = ?", brand).Error
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("products by category row error")
		return ps, err
	}

	return ps, nil
}

func GetFavourites(userId int64) ([]Product, error) {
	return []Product{}, nil
}

func SearchProducts(search string, offset string, userId int64) ([]Product, error) {
	ps := make([]Product, 0)
	searchQuery := "%" + search + "%"
	rows, err := GetDB().Query(
		`SELECT * FROM
			(SELECT
			Product.guid,
			Product.site_id,
			Product.url_id,
			Product.url,
			Product.category,
			Product.brand,
			Product.title,
			Product.description,
			Product.price,
			Product.img,
			Product.approved,
			MAX(IF(Favourites.user_id = ?, TRUE, FALSE)) is_favourited
				FROM (SELECT
				Products.guid,
				Products.site_id,
				Products.url_id,
				Urls.url,
				Products.category,
				Products.brand,
				Products.title,
				Products.description,
				Products.price,
				Products.img,
				Products.approved
				FROM Products
				INNER JOIN Urls
				WHERE Products.url_id = Urls.id
				AND Products.approved = 1
				) AS Product
			LEFT OUTER JOIN Favourites
			ON Product.guid = Favourites.product_guid
		GROUP BY Product.guid) AS Result
		WHERE Result.title LIKE ?
		OR Result.brand LIKE ?
		LIMIT 18 OFFSET ?`,
		userId,
		searchQuery,
		searchQuery,
		offset,
	)
	if err != nil {
		log.WithError(err).Error("failed to query products by search")
	}
	defer rows.Close()
	for rows.Next() {
		var p Product
		err := rows.Scan(
			&p.Guid,
			&p.SiteId,
			&p.UrlId,
			&p.Url,
			&p.Category,
			&p.Brand,
			&p.Title,
			&p.Description,
			&p.Price,
			&p.Img,
			&p.Approved,
			//			&p.IsFavourited,
		)
		if err != nil {
			log.WithError(err).Error("failed to get products by search rows")
			return ps, err
		}
		ps = append(ps, p)
	}
	err = rows.Err()
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("products by search row error")
		return ps, err
	}
	return ps, nil
}

// GetProductsByApprovedStatus gets a list of products based on its approved status
// Used by Admin
func GetProductsByApprovedStatus(approvedStatus bool, offset int64, itemsPerPage int64, search, category, siteID string, sortBy int) ([]Product, error) {
	var ps []Product
	db := GetDBv2().
		Table("Products").
		Where("approved = ?", approvedStatus).
		Offset(int(offset)).
		Limit(int(itemsPerPage)).
		Order("created_at desc")

	if category != "" {
		db = db.Where("category = ?", category)
	}
	if siteID != "" {
		db = db.Where("domain(product_url) = ?", siteID)
	}
	if search != "" {
		searchQuery := "%" + search + "%"
		db = db.Where("title ilike ? or brand ilike ?", searchQuery, searchQuery)
	}

	err := db.Scan(&ps).Error
	if err != nil {
		log.WithError(err).Error("failed to get products")
		return nil, err
	}
	return ps, nil
}

// UpdateApprovedStatus retrieves status of product of whether it is approved or not
// Used by Admin
func (p *Product) UpdateApprovedStatus(guid []string, approved bool) error {
	stringID := strings.Join(guid, ", ")
	stmt, err := GetDB().Prepare(`
		ALTER TABLE Products 
		UPDATE 
        approved = ?
		WHERE guid IN ( ? )`)
	if err != nil {
		log.WithError(err).Error("failed prepare update approved status")
		return err
	}
	_, err = stmt.Exec(
		approved,
		stringID,
	)
	if err != nil {
		log.WithError(err).Error("failed to run exec on update product")
		return err
	}

	return nil
}

// GetProductsCountByStatus retrieves the count of approved products from the database
// Used by Admin
func GetProductsCountByStatus(approvedStatus bool, search, category, siteId string) (count int64, err error) {
	searchQuery := "%" + search + "%"
	err = GetDB().QueryRow(
		`SELECT COUNT(*) FROM (
			SELECT * FROM Products 
			WHERE Products.category = IF(? = '',category, ?)
			AND domain(Products.product_url) = IF(? = '', domain(Products.product_url), ?)
			AND approved = ?
		) AS Product
		WHERE Product.title ILIKE ?
		OR Product.brand ILIKE ?
		`,
		category,
		category,
		siteId,
		siteId,
		approvedStatus,
		searchQuery,
		searchQuery,
	).Scan(
		&count,
	)
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("products count row error")
		return count, err
	}
	return count, nil
}

// to add color and size once it got added.

func GetProductByGuid(guid string) (Product, error) {
	var p Product
	err := GetDBv2().Table("Products").Find(&p, "guid = ?", guid).Error
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("fetch product failed")
		return Product{}, err
	}
	return p, nil
}

// GetProductsByGuid returns mutiple products by their IDs
func GetProductsByGuid(guid []string) ([]Product, error) {
	var ps []Product
	stringID := strings.Join(guid, ", ")
	err := GetDBv2().Table("Products").
		Where("guid in (?)", stringID).
		Scan(&ps).Error

	if err != nil {
		log.WithError(err).Error("failed to query products by id")
		return nil, err
	}

	return ps, nil
}

// UpdateProductByGuid will update the specific product.
func UpdateProductByGuid(guid string,
	title, description string,
	price, retailPrice int64,
	color, size, shoeSize, subCategory, featured string,
	approved bool) error {
	err := GetDBv2().Table("Products").Exec(
		`ALTER TABLE Products 
			UPDATE 
				title = ?,
				description = ?,
				price = ?,
				retail_price = ?,
				color = ?,
				size = ?,
				shoe_size = ?,
				sub_category = ?,
				featured = ?,
				approved = ? 
				WHERE guid = ?`,
		title,
		description,
		price,
		retailPrice,
		color,
		size,
		shoeSize,
		subCategory,
		featured,
		approved,
		guid).Error
	if err != nil {
		return errors.Wrap(err, "could not update product")
	}

	return nil
}

func DeleteProductById(id int64) error {
	//
	//log.Infof("delete: %+v", id)
	//stmt, err := GetDB().Prepare(`
	//	DELETE FROM Products
	//	WHERE id = ?
	//`)
	//if err != nil {
	//	log.WithError(err).Error("failed prepare delete Products statement")
	//	return err
	//}
	//_, err = stmt.Exec(
	//	id,
	//)
	//if err != nil {
	//	log.WithError(err).Error("failed to run exec on delete Products")
	//	return err
	//}
	return nil
}

// FIXME: SQL here is funky
// use an SQL builder
func DeleteProductsById(id []uint64) error {
	//var IDs []string
	//for _, i := range id {
	//	IDs = append(IDs, strconv.FormatUint(i, 10))
	//}
	//stringID := strings.Join(IDs, ", ")
	//log.Infof("delete: %+v", id)
	//stmt, err := GetDB().Prepare(`
	//	DELETE FROM Products
	//	WHERE id IN ` + `(` + stringID + `)`)
	//if err != nil {
	//	log.WithError(err).Error("failed prepare delete Products statement")
	//	return err
	//}
	//_, err = stmt.Exec()
	//if err != nil {
	//	log.WithError(err).Error("failed to run exec on delete Products")
	//	return err
	//}
	return nil
}

func GetFavouritesCount(userId int64) (count uint64, err error) {
	//TODO: rewrite

	//err = GetDB().QueryRow(
	//	`SELECT count(*) FROM
	//	(SELECT
	//	Products.id,
	//	Products.approved,
	//	MAX(IF(Favourites.user_id = ?, TRUE, FALSE)) is_favourited
	//	FROM Products
	//	INNER JOIN Favourites
	//	ON Products.id = Favourites.product_id
	//	WHERE Products.approved = 1
	//	GROUP BY Products.id) AS Final
	//	WHERE Final.is_favourited = 1`,
	//	userId,
	//).Scan(
	//	&count,
	//)
	//if err != nil && err != sql.ErrNoRows {
	//	log.WithError(err).Error("products count row error")
	//	return count, err
	//}
	//return count, nil
	return 0, nil
}

// Delete product by product_url
// Not used yet
func DeleteProductByProductUrl(productUrl string) error {
	var products []Product
	err := GetDBv2().Where("product_url = ?", productUrl).Delete(&products)
	if err != nil {
		log.Error("fetch product failed")
		return errors.Errorf("could not delete product by URL %s", err)
	}
	fmt.Println(len(products))
	return nil
}
