package main

import (
	"database/sql"
	"fmt"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// getProducts retrieves all products from the database
func getProducts(db *sql.DB) ([]Product, error) {
	query := "SELECT id, name, quantity, price FROM products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil { // Check for any error encountered during the rows iteration
		return nil, err
	}
	return products, nil
}

// getProduct retrieves a single product by its ID from the database
func (p *Product) getProduct(db *sql.DB) error {
	// Use parameterized queries to prevent SQL injection
	query := "SELECT name, quantity, price FROM products WHERE id = ?"
	err := db.QueryRow(query, p.ID).Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product with ID %d not found", p.ID)
		}
		return err
	}
	return nil
}

// createProduct inserts a new product into the database
func (p *Product) createProduct(db *sql.DB) error {
	query := "INSERT INTO products (name, quantity, price) VALUES (?, ?, ?)"
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

// updateProduct updates an existing product in the database
func (p *Product) updateProduct(db *sql.DB) (sql.Result, error) {
	query := "UPDATE products SET name = ?, quantity = ?, price = ? WHERE id = ?"
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price, p.ID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// deleteProduct deletes a product from the database
func (p *Product) deleteProduct(db *sql.DB) error {
	result, err := db.Exec("DELETE FROM products WHERE id=?", p.ID)
	if err != nil {
		return fmt.Errorf("could not delete product: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
