package product

import (
	"sauravkattel/agri/src/lib"
	"strings"

	"github.com/jmoiron/sqlx"
)

func AddProduct(db *sqlx.DB, products lib.Product, attrib lib.Attrib, userId string) error {
	var id string
	err := db.QueryRowx("INSERT INTO products(name, dec, user_id) VALUES ($1,$2,$3) RETURNING id", products.Name, products.Description, userId).Scan(&id)

	if err != nil {
		return err
	}

	slug := strings.Join(strings.Split(products.Name, " "), "_")
	attrib.Product_id = &id
	err = addAttrib(db, slug, attrib)
	return err
}

func addAttrib(db *sqlx.DB, slug string, attrib lib.Attrib) error {
	_, err := db.Exec(
		"INSERT INTO product_attrib(price, quantity,status, product_id,slug) VALUES($1,$2,$3,$4,$5)",
		attrib.Price,
		attrib.Quantity,
		attrib.Status,
		attrib.Product_id,
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

func GetProductsBySlug() {

}
