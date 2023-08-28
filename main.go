package main

import (
	"fmt"
	"net/http"
	"product-catalogue-Telkom-LKPP/internal/handlers"
	"product-catalogue-Telkom-LKPP/internal/repositories"
	"product-catalogue-Telkom-LKPP/internal/server"
)

func main() {
	// Create a database connection
	db, err := repositories.NewDBConnection()
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return // Stop the application
	}
	defer db.Close() // Close the database connection when the application exits

	productRepo := repositories.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepo)

	router := server.NewRouter(productHandler)

	fmt.Println("Server is running properly")
	http.ListenAndServe(":8080", router)
}
