package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"product-catalogue-Telkom-LKPP/internal/models"
	"product-catalogue-Telkom-LKPP/internal/repositories"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"fmt"
	"io/ioutil"
)

type ProductHandler struct {
	ProductRepo repositories.ProductRepository
}

func NewProductHandler(productRepo repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{
		ProductRepo: productRepo,
	}
}

func (h *ProductHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	// Extract image ID from the URL parameter
	imageID := chi.URLParam(r, "imageID")

	// Check if the image ID includes a file extension
	if filepath.Ext(imageID) == "" {
		http.Error(w, "Invalid image ID format", http.StatusBadRequest)
		return
	}

	// Construct the file path for the image
	filePath := fmt.Sprintf("internal/repositories/images/%s", imageID)

	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set the appropriate content type based on the image extension
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// Copy the image data to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Failed to serve image", http.StatusInternalServerError)
	}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from the URL parameter
	productID := chi.URLParam(r, "productID")

	// Get the product details from the repository
	product, err := h.ProductRepo.GetProductByID(productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Convert product images to URLs
	imageBaseURL := "http://localhost:8080/products/images/"
	for _, img := range product.Images {
		img.URL = fmt.Sprintf("%s%s%s", imageBaseURL, img.ID, img.Type)
	}

	// Marshal the product data to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		http.Error(w, "Failed to marshal product data", http.StatusInternalServerError)
		return
	}

	// Set the content type and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(productJSON)
}

func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	title := query.Get("title")
	etalase := query.Get("etalase")
	category := query.Get("category")
	sortBy := query.Get("sortBy")
	pageStr := query.Get("page")
	perPageStr := query.Get("perPage")

	// Convert page and perPage parameters to integers with default values
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 0
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 {
		perPage = 10
	}

	// Create the ProductQuery struct
	productQuery := &models.ProductQuery{
		Title:    title,
		Etalase:  etalase,
		Category: category,
		SortBy:   sortBy,
	}

	// Get the list of products from the repository
	products, err := h.ProductRepo.SearchProducts(productQuery, page, perPage)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	response := struct {
		Data []*models.Product `json:"data"`
		Meta struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
		} `json:"meta"`
	}{
		Data: products,
		Meta: struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
		}{
			Page:  page + 1,
			Limit: perPage,
		},
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse JSON data from the request body
	var requestBody models.ProductRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Process images
	var images []*models.ProductImage

	for _, base64Image := range requestBody.Images {
		imageData, err := base64.StdEncoding.DecodeString(base64Image)
		if err != nil {
			http.Error(w, "Failed to decode base64 image", http.StatusInternalServerError)
			return
		}

		// Detect the image type
		ext := detectImageTypeByData(imageData)
		if ext == "" {
			http.Error(w, "Invalid image type", http.StatusBadRequest)
			return
		}

		// Generate UUID for the image
		imgID := uuid.New()

		// Construct the file path for the image
		filePath := filepath.Join("internal", "repositories", "images", fmt.Sprintf("%s%s", imgID.String(), ext))

		// Store the image file locally
		err = ioutil.WriteFile(filePath, imageData, 0644)
		if err != nil {
			http.Error(w, "Failed to store image", http.StatusInternalServerError)
			return
		}

		// Create a ProductImage struct and append to the images slice
		images = append(images, &models.ProductImage{
			ID:       imgID,
			FilePath: filePath,
			Type:     ext,
		})
	}

	// Generate UUID for the product
	productID := uuid.New()

	// Create a Product struct with the extracted data
	product := &models.Product{
		ID:          productID,
		SKU:         requestBody.SKU,
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Category:    requestBody.Category,
		Etalase:     requestBody.Etalase,
		Images:      images,
		Weight:      requestBody.Weight,
		Price:       requestBody.Price,
	}

	// Call the CreateProduct method of the repository to insert the product into the database
	err = h.ProductRepo.CreateProduct(product)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Product created successfully %s", productID)))
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from URL parameter
	productID := chi.URLParam(r, "productID")

	// Parse JSON data from the request body
	var requestBody models.ProductRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Process images
	var images []*models.ProductImage

	for _, base64Image := range requestBody.Images {
		imageData, err := base64.StdEncoding.DecodeString(base64Image)
		if err != nil {
			http.Error(w, "Failed to decode base64 image", http.StatusInternalServerError)
			return
		}

		// Detect the image type
		ext := detectImageTypeByData(imageData)
		if ext == "" {
			http.Error(w, "Invalid image type", http.StatusBadRequest)
			return
		}

		// Generate UUID for the image
		imgID := uuid.New()

		// Construct the file path for the image
		filePath := filepath.Join("internal", "repositories", "images", fmt.Sprintf("%s%s", imgID.String(), ext))

		// Store the image file locally
		err = ioutil.WriteFile(filePath, imageData, 0644)
		if err != nil {
			http.Error(w, "Failed to store image", http.StatusInternalServerError)
			return
		}

		// Create a ProductImage struct and append to the images slice
		images = append(images, &models.ProductImage{
			ID:       imgID,
			FilePath: filePath,
			Type:     ext,
		})
	}

	// Create a Product struct with the extracted data (similar to CreateProduct)
	updatedProduct := &models.Product{
		ID:          uuid.MustParse(productID),
		SKU:         requestBody.SKU,
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Category:    requestBody.Category,
		Etalase:     requestBody.Etalase,
		Images:      images,
		Weight:      requestBody.Weight,
		Price:       requestBody.Price,
	}

	// Update the product in the repository (similar to CreateProduct)
	err = h.ProductRepo.UpdateProduct(productID, updatedProduct)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

func detectImageTypeByData(data []byte) string {
	// Define magic numbers for various image formats
	jpegMagic := []byte{0xFF, 0xD8, 0xFF}
	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47}
	gifMagic := []byte("GIF")

	// Compare the first few bytes of data with magic numbers
	if bytes.HasPrefix(data, jpegMagic) {
		return ".jpg"
	} else if bytes.HasPrefix(data, pngMagic) {
		return ".png"
	} else if bytes.HasPrefix(data, gifMagic) {
		return ".gif"
	}

	return "" // Return empty string for unknown image types
}
