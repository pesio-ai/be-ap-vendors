package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pesio-ai/be-go-common/logger"
	"github.com/pesio-ai/be-vendors-service/internal/service"
)

// HTTPHandler handles HTTP requests
type HTTPHandler struct {
	service *service.VendorService
	log     *logger.Logger
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(service *service.VendorService, log *logger.Logger) *HTTPHandler {
	return &HTTPHandler{
		service: service,
		log:     log,
	}
}

// CreateVendor handles create vendor HTTP requests
func (h *HTTPHandler) CreateVendor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req service.CreateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Get user ID from JWT token
	// req.CreatedBy = "system" // Leave empty for NULL

	vendor, err := h.service.CreateVendor(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vendor)
}

// GetVendor handles get vendor HTTP requests
func (h *HTTPHandler) GetVendor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vendorID := r.URL.Query().Get("id")
	entityID := r.URL.Query().Get("entity_id")

	if vendorID == "" || entityID == "" {
		http.Error(w, "Vendor ID and Entity ID are required", http.StatusBadRequest)
		return
	}

	vendor, err := h.service.GetVendor(r.Context(), vendorID, entityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}

// GetVendorByCode handles get vendor by code HTTP requests
func (h *HTTPHandler) GetVendorByCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vendorCode := r.URL.Query().Get("vendor_code")
	entityID := r.URL.Query().Get("entity_id")

	if vendorCode == "" || entityID == "" {
		http.Error(w, "Vendor Code and Entity ID are required", http.StatusBadRequest)
		return
	}

	vendor, err := h.service.GetVendorByCode(r.Context(), vendorCode, entityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}

// ListVendors handles list vendors HTTP requests
func (h *HTTPHandler) ListVendors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		http.Error(w, "Entity ID is required", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	vendorType := r.URL.Query().Get("vendor_type")
	activeOnly := r.URL.Query().Get("active_only") == "true"

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	var vendorTypePtr *string
	if vendorType != "" {
		vendorTypePtr = &vendorType
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	vendors, total, err := h.service.ListVendors(r.Context(), entityID, statusPtr, vendorTypePtr, activeOnly, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"vendors":  vendors,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// UpdateVendor handles update vendor HTTP requests
func (h *HTTPHandler) UpdateVendor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req service.UpdateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Get user ID from JWT token
	// req.UpdatedBy = "system" // Leave empty for NULL

	vendor, err := h.service.UpdateVendor(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}

// DeleteVendor handles delete vendor HTTP requests
func (h *HTTPHandler) DeleteVendor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vendorID := r.URL.Query().Get("id")
	entityID := r.URL.Query().Get("entity_id")

	if vendorID == "" || entityID == "" {
		http.Error(w, "Vendor ID and Entity ID are required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteVendor(r.Context(), vendorID, entityID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ValidateVendor handles validate vendor HTTP requests
func (h *HTTPHandler) ValidateVendor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vendorID := r.URL.Query().Get("id")
	entityID := r.URL.Query().Get("entity_id")

	if vendorID == "" || entityID == "" {
		http.Error(w, "Vendor ID and Entity ID are required", http.StatusBadRequest)
		return
	}

	valid, message, err := h.service.ValidateVendor(r.Context(), vendorID, entityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   valid,
		"message": message,
	})
}

// GetVendorContacts handles get vendor contacts HTTP requests
func (h *HTTPHandler) GetVendorContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vendorID := r.URL.Query().Get("vendor_id")
	if vendorID == "" {
		http.Error(w, "Vendor ID is required", http.StatusBadRequest)
		return
	}

	contacts, err := h.service.GetVendorContacts(r.Context(), vendorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"contacts": contacts,
	})
}

// AddVendorContact handles add vendor contact HTTP requests
func (h *HTTPHandler) AddVendorContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req service.AddContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	contact, err := h.service.AddVendorContact(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}

// GetPaymentTerms handles get payment terms HTTP requests
func (h *HTTPHandler) GetPaymentTerms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	terms, err := h.service.GetPaymentTerms(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"payment_terms": terms,
	})
}

// UpdateBalance handles update vendor balance HTTP requests
func (h *HTTPHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		VendorID string `json:"vendor_id"`
		EntityID string `json:"entity_id"`
		Amount   int64  `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.VendorID == "" || req.EntityID == "" {
		http.Error(w, "Vendor ID and Entity ID are required", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateBalance(r.Context(), req.VendorID, req.EntityID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}
