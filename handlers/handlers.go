package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ognjen-it/go-saas-microservice/entity"
)

// HealthCheck is used to get
func HealthCheckHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data, err := entity.HealthCheck()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Write the body with JSON data
		// rw.Header().Add("content-type", "application/json")
		rw.Write(data)
	}
}

func MaxHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		type Person struct {
			Customer_id   string `json:"customer_id"`
			Customer_name string `json:"customer_name"`
			Age           int    `json:"age"`
		}

		type listOfPersons struct {
			Data []Person `json:"data"`
		}

		var myList listOfPersons
		var s []int

		err := json.NewDecoder(r.Body).Decode(&myList)
		if err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}

		for _, p := range myList.Data {
			s = append(s, p.Age)
		}

		max := s[0]
		for _, v := range s {
			if v > max {
				max = v
			}
		}

		fmt.Fprintf(rw, "Max age of cusomters is: %+v", max)
	}
}

// GetProductsHandler is used to get data inside the products defined on our product catalog
func GetProductsHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data, err := entity.GetProducts()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Write the body with JSON data
		rw.Header().Add("content-type", "application/json")
		rw.WriteHeader(http.StatusFound)
		rw.Write(data)
	}
}

// GetProductHandler is used to get data inside the products defined on our product catalog
func GetProductHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Read product ID
		productID := mux.Vars(r)["id"]
		product, err := entity.GetProduct(productID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		responseData, err := json.Marshal(product)
		if err != nil {
			// Check if it is No product error or any other error
			if errors.Is(err, entity.ErrNoProduct) {
				// Write Header if no related product found.
				rw.WriteHeader(http.StatusNoContent)
			} else {
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		// Write body with found product
		rw.Header().Add("content-type", "application/json")
		rw.WriteHeader(http.StatusFound)
		rw.Write(responseData)
	}
}

// CreateProductHandler is used to create a new product and add to our product store.
func CreateProductHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Read incoming JSON from request body
		data, err := ioutil.ReadAll(r.Body)
		// If no body is associated return with StatusBadRequest
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		// Check if data is proper JSON (data validation)
		var product entity.Product
		err = json.Unmarshal(data, &product)
		if err != nil {
			rw.WriteHeader(http.StatusExpectationFailed)
			rw.Write([]byte("Invalid Data Format"))
			return
		}
		err = entity.AddProduct(product)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// return after writing Body
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte("Added New Product"))
	}
}

// DeleteProductHandler deletes the product with given ID.
func DeleteProductHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Read product ID
		productID := mux.Vars(r)["id"]
		err := entity.DeleteProduct(productID)
		if err != nil {
			// Check if it is No product error or any other error
			if errors.Is(err, entity.ErrNoProduct) {
				// Write Header if no related product found.
				rw.WriteHeader(http.StatusNoContent)
			} else {
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		// Write Header with Accepted Status (done operation)
		rw.WriteHeader(http.StatusAccepted)
	}
}

// UpdateProductHandler deletes the product with given ID.
func UpdateProductHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Read product ID
		productID := mux.Vars(r)["id"]
		err := entity.DeleteProduct(productID)
		if err != nil {
			if errors.Is(err, entity.ErrNoProduct) {
				rw.WriteHeader(http.StatusNoContent)
			} else {
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		// Read incoming JSON from request body
		data, err := ioutil.ReadAll(r.Body)
		// If no body is associated return with StatusBadRequest
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		// Check if data is proper JSON (data validation)
		var product entity.Product
		err = json.Unmarshal(data, &product)
		if err != nil {
			rw.WriteHeader(http.StatusExpectationFailed)
			rw.Write([]byte("Invalid Data Format"))
			return
		}
		// Addproduct with the requested body
		err = entity.AddProduct(product)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Write Header if no related product found.
		rw.WriteHeader(http.StatusAccepted)
	}
}
