// internal/repositories/product_repository.go

package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"product-catalogue-Telkom-LKPP/internal/models"
	"strconv"
	"strings"
)

type ProductRepository interface {
	GetProductByID(productID string) (*models.Product, error)
	SearchProducts(query *models.ProductQuery, page, perPage int) ([]*models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(productID string, product *models.Product) error
}

type productRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (repo *productRepository) GetProductByID(productID string) (*models.Product, error) {
	// Prepare the SQL statement
	query := `
			SELECT
					p.id, p.sku, p.title, p.description, p.category, p.etalase, p.images, p.weight, p.price,
					COALESCE(AVG(pr.rating),0) as rating
				FROM
					products p
				LEFT JOIN
					product_reviews pr on p.id = pr.product_id
			WHERE
				p.id = $1
			GROUP BY p.id
	`

	row := repo.DB.QueryRow(query, productID)

	var product models.Product
	var imagesJSON []byte

	// Scan the retrieved row into the product struct
	err := row.Scan(
		&product.ID,
		&product.SKU,
		&product.Title,
		&product.Description,
		&product.Category,
		&product.Etalase,
		&imagesJSON,
		&product.Weight,
		&product.Price,
		&product.Rating,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSONB images data into the product.Images slice
	err = json.Unmarshal(imagesJSON, &product.Images)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (repo *productRepository) SearchProducts(query *models.ProductQuery, page, perPage int) ([]*models.Product, error) {
	// Prepare the SQL statement
	sql := `
	select
		p.id,
		p.sku,
		p.title,
		p.description,
		p.category,
		p.etalase,
		p.images,
		p.weight,
		p.price,
		COALESCE(AVG(pr.rating),0) as rating
	from
		products p
	left join
			product_reviews pr on p.id = pr.product_id
	where
	`

	var whereConditions []string
	var args []interface{}
	var argCounter = 1

	if query.Title != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("p.title ILIKE $%d", argCounter))
		args = append(args, "%"+query.Title+"%")
		argCounter++
	}

	if query.Etalase != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("p.etalase ILIKE $%d", argCounter))
		args = append(args, "%"+query.Etalase+"%")
		argCounter++
	}

	if query.Category != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("p.category ILIKE $%d", argCounter))
		args = append(args, "%"+query.Category+"%")
		argCounter++
	}

	if query.SKU != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("p.sku ILIKE $%d", argCounter))
		args = append(args, "%"+query.SKU+"%")
		argCounter++
	}

	if len(whereConditions) == 0 {
		return nil, errors.New("Invalid query parameters")
	}

	sql += strings.Join(whereConditions, " AND ")

	var sortField string
	switch query.SortBy {
	case "newest":
		sortField = "p.created_at DESC"
	case "oldest":
		sortField = "p.created_at ASC"
	case "highestRated":
		sortField = "rating DESC"
	case "lowestRated":
		sortField = "rating ASC"
	default:
		sortField = "p.created_at DESC" // Default to newest
	}

	sql += `
				GROUP BY p.id
        ORDER BY ` + sortField + `
        LIMIT $` + strconv.Itoa(argCounter) + `
        OFFSET $` + strconv.Itoa(argCounter+1)

	args = append(args, perPage, page*perPage)

	rows, err := repo.DB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*models.Product{}
	for rows.Next() {
		var product models.Product
		var imagesJSON []byte

		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Title,
			&product.Description,
			&product.Category,
			&product.Etalase,
			&imagesJSON,
			&product.Weight,
			&product.Price,
			&product.Rating,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSONB images data into the product.Images slice
		err = json.Unmarshal(imagesJSON, &product.Images)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (repo *productRepository) CreateProduct(product *models.Product) error {
	// Convert images slice to JSONB data
	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		return fmt.Errorf("failed to marshal images to JSON: %v", err)
	}

	// Insert new product record into the database
	_, err = repo.DB.Exec(`
		INSERT INTO products (id, sku, title, description, category, etalase, images, weight, price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, product.ID, product.SKU, product.Title, product.Description, product.Category, product.Etalase, imagesJSON, product.Weight, product.Price)
	if err != nil {
		return fmt.Errorf("failed to insert product: %v", err)
	}

	return nil
}

func (repo *productRepository) UpdateProduct(productID string, product *models.Product) error {
	// Convert images slice to JSONB data
	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		return err
	}

	// Prepare the SQL statement
	query := `
			UPDATE products
			SET
					sku = $1,
					title = $2,
					description = $3,
					category = $4,
					etalase = $5,
					images = $6,
					weight = $7,
					price = $8
			WHERE
					id = $9
	`

	_, err = repo.DB.Exec(
		query,
		product.SKU,
		product.Title,
		product.Description,
		product.Category,
		product.Etalase,
		imagesJSON,
		product.Weight,
		product.Price,
		productID,
	)
	if err != nil {
		return err
	}

	return nil
}
