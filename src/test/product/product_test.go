package test

import (
	"sauravkattel/agri/src/lib"
	"sauravkattel/agri/src/product"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a SQLx DB instance from the mockDB
	db := sqlx.NewDb(mockDB, "sqlmock")

	// Mock data
	pageNumber := 1
	pageSize := 10

	// Mock expected query and return values
	rows := sqlmock.NewRows([]string{
		"products.id", "products.name", "products.description", "products.created_at",
		"pa.id as attrib_id", "pa.price", "pa.quantity", "pa.status", "pa.slug",
	}).AddRow(
		1, "Product 1", "Description 1", "2024-07-11 10:00:00",
		101, 19.99, 50, "1", "product-xyz",
	)

	mock.ExpectQuery(`
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
		FROM products JOIN product_attrib pa ON products.id = pa.product_id WHERE pa.status <> '0' ORDER BY products.created_at LIMIT \$1 OFFSET \$2
	`).WithArgs(pageSize, (pageNumber-1)*pageSize).WillReturnRows(rows)

	// Call GetProducts function
	result, err := product.GetProducts(db, pageNumber, pageSize)

	// Verify results
	assert.NoError(t, err, "GetProducts should not return an error")
	assert.NotNil(t, result, "Result should not be nil")

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "Expectations were not met")
}

func TestAddProduct(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a SQLx DB instance from the mockDB
	db := sqlx.NewDb(mockDB, "sqlmock")

	// Mock data
	products := lib.Product{
		Name:        "Test Product",
		Description: "Test description",
	}

	attrib := lib.Attrib{
		Price:    10.99,
		Quantity: 100,
		Status:   "active",
	}

	userID := "user123"

	// Mock expected query and return values for AddProduct function
	mock.ExpectQuery(`INSERT INTO products`).WithArgs(products.Name, products.Description, userID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("product123"))

	// Mock expected query and return values for addAttrib function
	mock.ExpectExec(`INSERT INTO product_attrib`).WithArgs(attrib.Price, attrib.Quantity, attrib.Status, "product123", "Test_Product").WillReturnResult(sqlmock.NewResult(1, 1))

	slug := "Test_Product"
	// Call AddProduct function
	err = product.AddProduct(db, products, attrib, userID, slug)

	// Verify results
	assert.NoError(t, err, "AddProduct should not return an error")

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "Expectations were not met")
}

func TestGetProductsBySlug(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a SQLx DB instance from the mockDB
	db := sqlx.NewDb(mockDB, "sqlmock")

	// Mock expected query and return values
	rows := sqlmock.NewRows([]string{
		"products.id", "products.name", "products.description", "products.created_at",
		"pa.id as attrib_id", "pa.price", "pa.quantity", "pa.status", "pa.slug",
	}).AddRow(
		1, "Product 1", "Description 1", "2024-07-11 10:00:00",
		101, 19.99, 50, "1", "product-xyz",
	)

	mock.ExpectQuery(`
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
		FROM products JOIN product_attrib pa ON products.id = pa.product_id WHERE pa.status <> '0' AND pa.slug = \$1
	`).WithArgs("product-xyz").WillReturnRows(rows)

	// Call GetProducts function
	result, err := product.GetProductsBySlug(db, "product-xyz")

	// Verify results
	assert.NoError(t, err, "GetProducts should not return an error")
	assert.NotNil(t, result, "Result should not be nil")

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "Expectations were not met")
}

func TestUpdateAttrib(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %s", err)
	}
	defer db.Close()

	// Create a new sqlx DB instance
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Define test cases
	tests := []struct {
		slug     string
		table    string
		column   string
		value    interface{}
		expected string
	}{
		{"test-slug", "product_attrib", "quantity", 10, `UPDATE product_attrib SET quantity = \$1 WHERE slug = \$2`},
		{"test-slug", "product_attrib", "price", 19.99, `UPDATE product_attrib SET price = \$1 WHERE slug = \$2`},
		{"test-slug", "product_attrib", "status", "available", `UPDATE product_attrib SET status = \$1 WHERE slug = \$2`},
	}

	for _, tt := range tests {
		mock.ExpectExec(tt.expected).
			WithArgs(tt.value, tt.slug).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := product.UpdateProductAttrib(sqlxDB, tt.slug, tt.table, tt.column, tt.value)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
