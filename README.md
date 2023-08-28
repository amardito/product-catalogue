# Product Catalogue API

This project is a RESTful API for managing a product catalog.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (Golang) installed on your machine
- PostgreSQL database
- Postman (for testing the API)

## Getting Started

To get started, follow these steps:

1. Clone this repository
2. Run this PostgreSQL script

```
-- Create products table
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50),
    etalase VARCHAR(50),
    images JSONB, -- Store image metadata as JSONB
    weight DECIMAL(10, 2),
    price DECIMAL(10, 2),
    created_at TIMESTAMP WITH TIME ZONE CURRENT_TIMESTAMP
);

-- Create product_reviews table
CREATE TABLE product_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES products(id),
    rating INT CHECK (rating >= 1 AND rating <= 5),
    review_comment TEXT
);

```

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

POST /review - Create a new review for product

## Postman Documentation

For detailed usage and examples, please refer to the Postman Documentation.
https://documenter.getpostman.com/view/9503318/2s9Y5YSi71#intro
