package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT * FROM products;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name, quantity, price FROM products WHERE id=%v;", p.ID)
	rows := db.QueryRow(query)
	err := rows.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

// func (p *product) createProduct(db *sql.DB) error {
// 	query := fmt.Sprintf("insert into products(name,quantity,price) values('%v', '%v', '%v')", p.Name, p.Quantity, p.Price)
// 	result, err := db.Exec(query)
// 	if err != nil {
// 		return err
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return err
// 	}
// 	p.ID = int(id)
// 	return nil
// }

func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO products (name, quantity, price) VALUES ($1, $2, $3) RETURNING id")
	var id int
	err := db.QueryRow(query, p.Name, p.Quantity, p.Price).Scan(&id)
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE products SET name='%v', quantity='%v', price='%v' WHERE id='%v'", p.Name, p.Quantity, p.Price, p.ID)
	result, err := db.Exec(query)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such row with the given ID exists")
	}

	return nil
}

func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM products WHERE id='%v'", p.ID)
	result, err := db.Exec(query)

	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such row with the given ID exists")
	}

	return nil
}
