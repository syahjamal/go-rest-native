package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//Struct Product
type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

var (
	database = make(map[string]Product)
)

//Fungsi respon json
func SetJSONResp(res http.ResponseWriter, message []byte, httpCode int) {
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(httpCode)
	res.Write(message)
}

func main() {
	//init db
	database["001"] = Product{ID: "001", Name: "Samsung Galaxy S7", Quantity: 10}
	database["002"] = Product{ID: "002", Name: "Nokia A3", Quantity: 5}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		message := []byte(`{"message": "server is up"}`)
		SetJSONResp(res, message, http.StatusOK)
	})

	//Get All Product
	http.HandleFunc("/get-products", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			message := []byte(`{"message": "Invalid http method"}`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}
		//Untuk mengkonfersi data map database ke list
		var products []Product

		for _, product := range database {
			products = append(products, product)
		}

		productJson, err := json.Marshal(&products)
		if err != nil {
			message := []byte(`{"message": "Error when parsing data"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}
		SetJSONResp(res, productJson, http.StatusOK)
	})

	//Add Product
	http.HandleFunc("/add-product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			message := []byte(`{"message": "Invalid http method"}`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}
		//Menampung body dari depan
		var product Product

		payload := req.Body

		defer req.Body.Close()

		err := json.NewDecoder(payload).Decode(&product)
		if err != nil {
			message := []byte(`{"message": "Error Parsing Data"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}
		database[product.ID] = product
		message := []byte(`{"message": "Success Create Product"}`)
		SetJSONResp(res, message, http.StatusCreated)
	})

	//Get Detail by ID
	http.HandleFunc("/get-product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			message := []byte(`{"message": "Invalid http method"}`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		//Validasi untuk manggil route harus pakai id
		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message": "Required product id"}`)
			SetJSONResp(res, message, http.StatusBadRequest)
			return
		}
		id := req.URL.Query()["id"][0]
		//Validasi data product jika tidak sesuai id nya
		product, ok := database[id]
		if !ok {
			message := []byte(`{"message": "product not found"}`)
			SetJSONResp(res, message, http.StatusOK)
			return
		}

		productJSON, err := json.Marshal(&product)
		if err != nil {
			message := []byte(`{"message": "some error when parsing data"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONResp(res, productJSON, http.StatusOK)
	})

	//Delete
	http.HandleFunc("/delete-products", func(res http.ResponseWriter, req *http.Request) {

		if req.Method != "DELETE" {
			message := []byte(`{"message": "Invalid http method"}`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message": "Required product id"}`)
			SetJSONResp(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]
		product, ok := database[id]
		if !ok {
			message := []byte(`{"message": "product not found"}`)
			SetJSONResp(res, message, http.StatusOK)
			return
		}

		delete(database, id)

		productJSON, err := json.Marshal(&product)
		if err != nil {
			message := []byte(`{"message": "some error when parsing data"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONResp(res, productJSON, http.StatusOK)

	})

	//Update CRUD
	http.HandleFunc("/update-products", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			message := []byte(`{"message": "Invalid http method"}`)
			SetJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message": "Required product id"}`)
			SetJSONResp(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]
		product, ok := database[id]
		if !ok {
			message := []byte(`{"message": "product not found"}`)
			SetJSONResp(res, message, http.StatusOK)
			return
		}

		var newProduct Product

		payload := req.Body

		defer req.Body.Close()

		err := json.NewDecoder(payload).Decode(&newProduct)
		if err != nil {
			message := []byte(`{"message": "error when parsing data"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		product.Name = newProduct.Name
		product.Quantity = newProduct.Quantity

		database[product.ID] = product

		productJSON, err := json.Marshal(&product)
		if err != nil {
			message := []byte(`{"message": "some error when parsing data"}`)
			SetJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONResp(res, productJSON, http.StatusOK)

	})

	err := http.ListenAndServe(":9080", nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
