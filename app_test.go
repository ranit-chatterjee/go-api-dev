package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DBUser, DBPassword, "test")
	if err != nil {
		log.Fatal("could not initialise the database")
	}
	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS inventory (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		quantity INT,
		price NUMERIC(10,2)
	);`

	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM inventory")
	a.DB.Exec("ALTER SEQUENCE inventory_id_seq RESTART WITH 1;")
	log.Println("table cleared")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT INTO inventory(name, quantity, price) VALUES('%v', '%v', '%v')", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}
func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 5001.00)
	request, err := http.NewRequest("GET", "/product/1", nil)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("no error")
	}
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected Status: %v, Received: %v", expectedStatusCode, actualStatusCode)
	} else {
		log.Printf("Received Status Code: %v", actualStatusCode)
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var p = []byte(`{"name": "chair", "quantity": "4", "price": "400"}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(p))
	req.Header.Set("Content-Type", "application/json")

	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "chair" {
		t.Errorf("Expected name %v, Got: %v", "chair", m["name"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10.0)
	request, err := http.NewRequest("GET", "/product/1", nil)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("no error")
	}
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("DELETE", "/product/1", nil)

	response = sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("GET", "/product/1", nil)

	response = sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)

}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10.0)
	request, _ := http.NewRequest("GET", "/product/1", nil)

	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	var oldVal map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldVal)

	var p = []byte(`{"name": "connector", "quantity": "4", "price": "40"}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(p))
	req.Header.Set("Content-Type", "application/json")

	response = sendRequest(req)

	var newVal map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newVal)

	if oldVal["id"] != newVal["id"] {
		t.Errorf("Expected: %v | Got: %v", oldVal["id"], newVal["id"])
	}
}
