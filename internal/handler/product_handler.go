package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Zhonghe-zhao/seckill-system/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(ps *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: ps}
}

type InitializeProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gt=0"`
	StartTime   string  `json:"start_time" binding:"required"`
	EndTime     string  `json:"end_time" binding:"required,gtfield=StartTime"`
}

func (h *ProductHandler) HandleInitializeProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		Price        float64 `json:"price"`
		InitialStock int64   `json:"initial_stock"`
		StartTimeStr string  `json:"start_time"` // e.g., "2025-06-01T10:00:00Z"
		EndTimeStr   string  `json:"end_time"`   // e.g., "2025-06-01T10:05:00Z"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.ID == "" || req.Name == "" || req.InitialStock <= 0 {
		http.Error(w, "product_id, name, and positive initial_stock are required", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTimeStr)
	if err != nil {
		http.Error(w, "Invalid start_time format (use RFC3339): "+err.Error(), http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTimeStr)
	if err != nil {
		http.Error(w, "Invalid end_time format (use RFC3339): "+err.Error(), http.StatusBadRequest)
		return
	}
	if !endTime.After(startTime) {
		http.Error(w, "end_time must be after start_time", http.StatusBadRequest)
		return
	}

	product, err := h.productService.InitializeProduct(r.Context(), req.ID, req.Name, req.Price, req.InitialStock, startTime, endTime)
	if err != nil {
		log.Printf("Error initializing product %s: %v", req.ID, err)
		http.Error(w, "Failed to initialize product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}
