package product

import (
	"fmt"
	"log"
	"sauravkattel/agri/src/lib"

	"github.com/jmoiron/sqlx"
)

func AddProduct(db *sqlx.DB, products lib.Product, attrib lib.Attrib, userId, slug string) error {
	var id string
	err := db.QueryRowx("INSERT INTO products(name, dec, user_id) VALUES ($1,$2,$3) RETURNING id", products.Name, products.Description, userId).Scan(&id)

	if err != nil {
		return err
	}

	err = addAttrib(db, slug, attrib, id)
	return err
}

func addAttrib(db *sqlx.DB, slug string, attrib lib.Attrib, productId string) error {
	_, err := db.Exec(
		"INSERT INTO product_attrib(price, quantity,status, product_id,slug) VALUES($1,$2,$3,$4,$5)",
		attrib.Price,
		attrib.Quantity,
		attrib.Status,
		productId,
		slug,
	)
	return err
}

func GetProducts(db *sqlx.DB, pageNumber, pageSize int) (*[]lib.ProductDetails, error) {
	var data []lib.ProductDetails

	offset := (pageNumber - 1) * pageSize
	err := db.Select(&data, `
		SELECT 
		products.id as id,
		products.name as name,
		products.dec as dec,
		products.created_at as created_at,
		pa.id as attrib_id,
		pa.price as price,
		pa.quantity as quantity,
		pa.slug as slug,
		pa.status as status,
		pa.product_id as product_id
		FROM products JOIN product_attrib pa ON products.id = pa.product_id WHERE pa.status <> '0' ORDER BY products.created_at LIMIT $1 OFFSET $2
	`, pageSize, offset)
	return &data, err
}

func GetProductsBySlug(db *sqlx.DB, slug string) (*lib.ProductDetails, error) {
	var data lib.ProductDetails
	err := db.QueryRowx(`
		SELECT 
		products.id as id,
		products.name as name,
		products.dec as dec,
		products.created_at as created_at,
		pa.id as attrib_id,
		pa.price as price,
		pa.quantity as quantity,
		pa.slug as slug,
		pa.status as status,
		pa.product_id as product_id
		FROM products JOIN product_attrib pa ON products.id = pa.product_id WHERE pa.status <> '0' AND pa.slug = $1
		`, slug).StructScan(&data)
	if err != nil {
		return nil, err
	}
	return &data, err
}

func UpdateProductAttrib[T any](db *sqlx.DB, slug, column string, value T) error {
	query := fmt.Sprintf("UPDATE product_attrib SET %s = $1 WHERE slug = $2", column)
	_, err := db.Exec(query, value, slug)
	return err
}

func UpdateProduct[T any](db *sqlx.DB, slug, column string, value T) error {
	query := fmt.Sprintf("UPDATE products SET %s = $1 WHERE id = (SELECT product_id FROM product_attrib WHERE slug = $2)", column)
	_, err := db.Exec(query, value, slug)
	return err
}

func DeleteProduct(db *sqlx.DB, userId, slug string) error {
	_, err := db.Exec(`DELETE FROM products WHERE  user_id = $1 AND id = (
	SELECT product_id FROM product_attrib WHERE slug = $2
	)`, userId, slug)

	log.Print(err)
	return err
}
