// internal/server/server.go

package server

import (
	"net/http"

	"product-catalogue-Telkom-LKPP/internal/handlers"

	"github.com/go-chi/chi"
)

func NewRouter(productHandler *handlers.ProductHandler) http.Handler {
	r := chi.NewRouter()

	// Add a handler for the root path
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Server is running properly"))
	})

	// Group the routes under "/products"
	r.Route("/products", func(productRouter chi.Router) {
		productRouter.Post("/", productHandler.CreateProduct)
		productRouter.Put("/{productID}", productHandler.UpdateProduct)
		productRouter.Get("/", productHandler.SearchProducts)
		productRouter.Get("/{productID}", productHandler.GetProduct)
		productRouter.Get("/images/{imageID}", productHandler.ServeImage)
	})

	// Group the routes under "/reviews"
	// TODO: create review product

	return r
}
