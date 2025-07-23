package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/jonilsonds9/goexpert-modulo-7-apis/internal/dto"
	"github.com/jonilsonds9/goexpert-modulo-7-apis/internal/entity"
	"github.com/jonilsonds9/goexpert-modulo-7-apis/internal/infra/database"
	entityPkg "github.com/jonilsonds9/goexpert-modulo-7-apis/pkg/entity"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create Product godoc
// @Summary      Create product
// @Description  Create products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateProductInput  true  "product request"
// @Success      201
// @Failure      500 {object}  Error
// @Router       /products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// List Products godoc
// @Summary      List Products
// @Description  get all products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page    query     int  false  "Page number"
// @Param        limit   query     int  false  "Number of products per page"
// @Success      200 {array}   entity.Product
// @Failure      404 {object}  Error
// @Failure      500 {object}  Error
// @Router       /products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	sort := r.URL.Query().Get("sort")

	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// GetProduct godoc
// @Summary      Get a Products
// @Description  Get a products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Product ID" Format(uuid)
// @Success      200 {object}   entity.Product
// @Failure      404
// @Failure      500 {object}  Error
// @Router       /products/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary      Update a Products
// @Description  Update a products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id    	path     string  true  "Product ID" Format(uuid)
// @Param        request  body      dto.CreateProductInput  true  "product request"
// @Success      200
// @Failure      404
// @Failure      500 {object}  Error
// @Router       /products/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
// @Summary      Delete a Product
// @Description  Delete a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id    	path     string  true  "Product ID" Format(uuid)
// @Success      200
// @Failure      404
// @Failure      500 {object}  Error
// @Router       /products/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
