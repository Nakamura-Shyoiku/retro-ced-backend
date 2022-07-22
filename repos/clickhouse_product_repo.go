package repos

import (
	_ "github.com/ClickHouse/clickhouse-go"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ulventech/retro-ced-backend/models"
)

// Product model
type Product struct {
	Category     string    `db:"category"`
	Brand        string    `db:"brand"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	Price        float64   `db:"price"`
	RetailPrice  float64   `db:"retail_price"`
	Model        string    `db:"model"`
	ItemNumber   string    `db:"item_number"`
	Condition    string    `db:"condition"`
	Accessories  string    `db:"accessories"`
	Measurements string    `db:"measurements"`
	Featured     string    `db:"featured"`
	Img          string    `db:"img"`
	Color        string    `db:"color"`
	Size         string    `db:"size"`
	SubCategory  string    `db:"sub_category"`
	ProductURL   string    `db:"product_url"`
	UpdatedAt    time.Time `db:"updated_at"`
	CreatedAt    time.Time `db:"created_at"`
}

type ClickhouseProductsRepo struct {
	DB *sqlx.DB
}

func NewClickhouseProductsRepo(url string) (*ClickhouseProductsRepo, error) {
	db, err := sqlx.Connect("clickhouse", url)
	if err != nil {
		return nil, err
	}
	return &ClickhouseProductsRepo{
		DB: db,
	}, nil
}

func (c *ClickhouseProductsRepo) Products(category string) ([]models.Product, error) {
	var prod []Product
	err := c.DB.Select(&prod, `select category,
       brand,
       title,
       description,
       price,
       retail_price,
       model,
       item_number,
       condition,
       accessories,
       measurements,
       featured,
       img,
       color,
       size,
       sub_category,
       product_url,
       created_at,
       updated_at
	from products
	where lowerUTF8(category) = ?;`, category)
	if err != nil {
		return nil, err
	}
	var res []models.Product
	for _, p := range prod {
		res = append(res, models.Product{
			CreatedAt:        p.CreatedAt,
			Category:         p.Category,
			Color:            p.Color,
			ProductCondition: p.Condition,
			ProductURL:       p.ProductURL,
			Img:              p.Img,
			Price:            int64(p.Price),
			Description:      p.Description,
			Title:            p.Title,
			Brand:            p.Brand,
			UpdatedAt:        p.UpdatedAt,
			Measurements:     p.Measurements,
			Size:             p.Size,
			Accessories:      p.Accessories,
		})
	}
	return res, nil
}
