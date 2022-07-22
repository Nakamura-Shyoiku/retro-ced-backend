package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ulventech/retro-ced-backend/models"
)

const limit = 10000

var c = 0

func querybatch(db *sqlx.DB, i int) []models.Product {
	var p []models.Product
	err := db.Select(&p, fmt.Sprintf(`select
        category,
		brand,
		title,
		description,
		price,
		retail_price,
		model,
		item_number,
		product_condition,
		accessories,
		measurements,
		img,
        last_updated,
      	created_at,
      	color,
      	size,
      	shoe_size,
        sub_category,
        product_url
		from Products
		where product_url is not null and created_at <= '2021-09-28 09:52:04'
		order by created_at desc 
		limit %d offset %d `, limit, i))

	if err != nil {
		log.Fatal(err)
	}
	return p
}

func main() {
	db, err := sqlx.Connect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		"root",
		"romsIANIIyHugtGF",
		"35.193.223.28",
		3306,
		"retroced",
	))
	if err != nil {
		log.Fatalf("unable to connect to clickhouse %s", err)
	}

	conn := os.Getenv("CLICKHOUSE_URL")
	ch, err := sqlx.Connect("clickhouse", conn)
	if err != nil {
		log.Fatalf("unable to connect to clickhouse %s", err)
	}

	for i := 0; i < 3_000_000; i += limit {
		log.Printf("inserting batch %d", i)
		res := querybatch(db, i)
		fmt.Printf("number of products %d\n", len(res))
		insertclickhouse(ch, res)
		log.Printf("done inserting batch %d", i)
	}
}

func insertclickhouse(ch *sqlx.DB, prods []models.Product) {

	if len(prods) == 0 {
		return
	}

	tx := ch.MustBegin()

	q := `insert into Products_clone(category,
                           brand,
                           title,
                           description,
                           price,
                           retail_price,
                           model,
                           item_number,
                           product_condition,
                           accessories,
                           measurements,
                           img,
                           last_updated,
                           created_at,
                           color,
                           size,
                           shoe_size,
                           sub_category,
                           product_url)  values
                                                (
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?,
                                                 ?
                                                )`
	stmt, err := tx.Prepare(q)
	if err != nil {
		log.Fatalf("unable to prepare tx %s", err)
	}

	for _, p := range prods {
		_, err = stmt.Exec(
			p.Category,
			p.Brand,
			p.Title,
			p.Description,
			p.Price,
			p.RetailPrice,
			p.Model,
			p.ItemNumber,
			p.ProductCondition,
			p.Accessories,
			p.Measurements,
			p.Img,
			p.UpdatedAt,
			p.CreatedAt,
			p.Color,
			p.Size,
			p.ShoeSize,
			p.SubCategory,
			p.ProductURL)
		if err != nil {
			log.Fatalf("unable to execute insert query  %s", err)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Fatalf("unable to commit tx %s", err)
	}

	c += len(prods)
	log.Printf("%d rows inserted ", c)
}
