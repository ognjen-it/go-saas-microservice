package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ognjen-it/go-saas-microservice/handlers"
)

func main() {
	// Create new Router
	router := mux.NewRouter()

	// route properly to respective handlers
	router.Handle("/products", handlers.GetProductHandler()).Methods("GET")
	router.Handle("/products", handlers.CreateProductHandler()).Methods("POST")

	// Create new server and assign the router
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("Staring Product Catalog server on Port 8080")
	// Start Server on defined port/host.
	server.ListenAndServe()
}
