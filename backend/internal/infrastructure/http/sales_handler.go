package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"multi-inventory/internal/application"

	"github.com/go-chi/chi/v5"
)

type SalesHandler struct {
	salesService *application.SalesService
}

func NewSalesHandler(salesService *application.SalesService) *SalesHandler {
	return &SalesHandler{salesService: salesService}
}

type CreateOrderRequest struct {
	UserID string `json:"user_id"` // In real app, get from context/token (UUID or numeric as string)
	Items  []struct {
		ItemID   int64 `json:"item_id"`
		Quantity int   `json:"quantity"`
	} `json:"items"`
}

func (h *SalesHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For MVP, if UserID is missing, use a dummy ID or 1
	// If missing, leave empty to allow NULL FK (column is nullable)

	order, err := h.salesService.CreateOrder(r.Context(), req.UserID, req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *SalesHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.salesService.ListOrders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func (h *SalesHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	order, err := h.salesService.GetOrder(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if order == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

type UpdateFulfillmentRequest struct {
	IsFulfilled bool `json:"is_fulfilled"`
}

func (h *SalesHandler) UpdateItemFulfillment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "itemId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateFulfillmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.salesService.UpdateItemFulfillment(r.Context(), id, req.IsFulfilled); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SalesHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListOrders)
	r.Post("/", h.CreateOrder)
	r.Get("/{id}", h.GetOrder)
	r.Put("/items/{itemId}/fulfillment", h.UpdateItemFulfillment)
	return r
}
