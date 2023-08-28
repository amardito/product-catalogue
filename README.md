# Product Catalogue API

This project is a RESTful API for managing a product catalog.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (Golang) installed on your machine
- PostgreSQL database
- Postman (for testing the API)

## Getting Started

To get started, follow these steps:

1. Clone this repository:
2. run this command
   `go mod tidy`
3. run the project with this command
   `go run main.go`

The API server should now be running at http://localhost:8080.

## Endpoints

GET /products - Search for products with optional query parameters.

POST /products - Create a new product.

PUT /products/{productID} - Update an existing product by ID.

GET /products/{productID} - Get a product by ID.

GET /products/images/{imageID} - Get an image by ID.

## Postman Documentation

For detailed usage and examples, please refer to the Postman Documentation.
https://documenter.getpostman.com/view/9503318/2s9Y5YSi71#intro
