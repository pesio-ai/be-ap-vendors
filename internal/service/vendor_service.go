package service

import (
	"context"
	"strings"

	"github.com/pesio-ai/be-go-common/errors"
	"github.com/pesio-ai/be-go-common/logger"
	"github.com/pesio-ai/be-vendors-service/internal/repository"
)

// VendorService handles vendor business logic
type VendorService struct {
	vendorRepo *repository.VendorRepository
	log        *logger.Logger
}

// NewVendorService creates a new vendor service
func NewVendorService(
	vendorRepo *repository.VendorRepository,
	log *logger.Logger,
) *VendorService {
	return &VendorService{
		vendorRepo: vendorRepo,
		log:        log,
	}
}

// CreateVendorRequest represents a create vendor request
type CreateVendorRequest struct {
	EntityID          string
	VendorCode        string
	VendorName        string
	LegalName         *string
	VendorType        string
	TaxID             *string
	IsTaxExempt       bool
	Is1099Vendor      bool
	Email             *string
	Phone             *string
	Fax               *string
	Website           *string
	AddressLine1      *string
	AddressLine2      *string
	City              *string
	StateProvince     *string
	PostalCode        *string
	Country           string
	PaymentTerms      string
	PaymentMethod     *string
	Currency          string
	CreditLimit       *int64
	BankName          *string
	BankAccountNumber *string
	BankRoutingNumber *string
	SwiftCode         *string
	IBAN              *string
	Notes             *string
	Tags              []string
	CreatedBy         string
}

// UpdateVendorRequest represents an update vendor request
type UpdateVendorRequest struct {
	ID                string
	EntityID          string
	VendorCode        string
	VendorName        string
	LegalName         *string
	VendorType        string
	Status            string
	TaxID             *string
	IsTaxExempt       bool
	Is1099Vendor      bool
	Email             *string
	Phone             *string
	Fax               *string
	Website           *string
	AddressLine1      *string
	AddressLine2      *string
	City              *string
	StateProvince     *string
	PostalCode        *string
	Country           string
	PaymentTerms      string
	PaymentMethod     *string
	Currency          string
	CreditLimit       *int64
	BankName          *string
	BankAccountNumber *string
	BankRoutingNumber *string
	SwiftCode         *string
	IBAN              *string
	Notes             *string
	Tags              []string
	UpdatedBy         string
}

// AddContactRequest represents an add contact request
type AddContactRequest struct {
	VendorID    string
	ContactType string
	FirstName   string
	LastName    string
	Title       *string
	Email       *string
	Phone       *string
	Mobile      *string
	IsPrimary   bool
	Notes       *string
}

// CreateVendor creates a new vendor
func (s *VendorService) CreateVendor(ctx context.Context, req *CreateVendorRequest) (*repository.Vendor, error) {
	// Validate vendor code is unique for entity
	existing, _ := s.vendorRepo.GetByCode(ctx, req.VendorCode, req.EntityID)
	if existing != nil {
		return nil, errors.AlreadyExists("vendor", req.VendorCode)
	}

	// Validate vendor type
	validTypes := map[string]bool{
		"supplier":         true,
		"contractor":       true,
		"service_provider": true,
		"consultant":       true,
		"utility":          true,
	}
	vendorType := strings.ToLower(req.VendorType)
	if !validTypes[vendorType] {
		return nil, errors.InvalidInput("vendor_type", "invalid vendor type")
	}

	// Validate currency
	if len(req.Currency) != 3 {
		return nil, errors.InvalidInput("currency", "currency must be 3-letter ISO code")
	}

	// Validate credit limit if set
	if req.CreditLimit != nil && *req.CreditLimit < 0 {
		return nil, errors.InvalidInput("credit_limit", "credit limit cannot be negative")
	}

	// Validate country code (should be 2-letter ISO)
	if len(req.Country) != 2 {
		return nil, errors.InvalidInput("country", "country must be 2-letter ISO code")
	}

	// Create vendor with pending approval status
	vendor := &repository.Vendor{
		EntityID:          req.EntityID,
		VendorCode:        strings.ToUpper(req.VendorCode),
		VendorName:        req.VendorName,
		LegalName:         req.LegalName,
		VendorType:        vendorType,
		Status:            "pending_approval",
		TaxID:             req.TaxID,
		IsTaxExempt:       req.IsTaxExempt,
		Is1099Vendor:      req.Is1099Vendor,
		Email:             req.Email,
		Phone:             req.Phone,
		Fax:               req.Fax,
		Website:           req.Website,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
		City:              req.City,
		StateProvince:     req.StateProvince,
		PostalCode:        req.PostalCode,
		Country:           strings.ToUpper(req.Country),
		PaymentTerms:      req.PaymentTerms,
		PaymentMethod:     req.PaymentMethod,
		Currency:          strings.ToUpper(req.Currency),
		CreditLimit:       req.CreditLimit,
		CurrentBalance:    0,
		BankName:          req.BankName,
		BankAccountNumber: req.BankAccountNumber,
		BankRoutingNumber: req.BankRoutingNumber,
		SwiftCode:         req.SwiftCode,
		IBAN:              req.IBAN,
		Notes:             req.Notes,
		Tags:              req.Tags,
		CreatedBy:         &req.CreatedBy,
	}

	if err := s.vendorRepo.Create(ctx, vendor); err != nil {
		return nil, err
	}

	s.log.Info().
		Str("vendor_id", vendor.ID).
		Str("vendor_code", vendor.VendorCode).
		Str("entity_id", req.EntityID).
		Msg("Vendor created")

	return vendor, nil
}

// GetVendor retrieves a vendor by ID
func (s *VendorService) GetVendor(ctx context.Context, id, entityID string) (*repository.Vendor, error) {
	return s.vendorRepo.GetByID(ctx, id, entityID)
}

// GetVendorByCode retrieves a vendor by code
func (s *VendorService) GetVendorByCode(ctx context.Context, code, entityID string) (*repository.Vendor, error) {
	return s.vendorRepo.GetByCode(ctx, code, entityID)
}

// UpdateVendor updates a vendor
func (s *VendorService) UpdateVendor(ctx context.Context, req *UpdateVendorRequest) (*repository.Vendor, error) {
	// Get existing vendor
	vendor, err := s.vendorRepo.GetByID(ctx, req.ID, req.EntityID)
	if err != nil {
		return nil, err
	}

	// Check if code is being changed and if new code is unique
	if req.VendorCode != vendor.VendorCode {
		existing, _ := s.vendorRepo.GetByCode(ctx, req.VendorCode, req.EntityID)
		if existing != nil {
			return nil, errors.AlreadyExists("vendor", req.VendorCode)
		}
	}

	// Validate vendor type
	vendorType := strings.ToLower(req.VendorType)
	if vendorType != "supplier" && vendorType != "contractor" && vendorType != "service_provider" &&
		vendorType != "consultant" && vendorType != "utility" {
		return nil, errors.InvalidInput("vendor_type", "invalid vendor type")
	}

	// Validate status
	status := strings.ToLower(req.Status)
	if status != "active" && status != "inactive" && status != "suspended" && status != "pending_approval" {
		return nil, errors.InvalidInput("status", "invalid vendor status")
	}

	// Validate credit limit if set
	if req.CreditLimit != nil && *req.CreditLimit < 0 {
		return nil, errors.InvalidInput("credit_limit", "credit limit cannot be negative")
	}

	// Update vendor
	vendor.VendorCode = strings.ToUpper(req.VendorCode)
	vendor.VendorName = req.VendorName
	vendor.LegalName = req.LegalName
	vendor.VendorType = vendorType
	vendor.Status = status
	vendor.TaxID = req.TaxID
	vendor.IsTaxExempt = req.IsTaxExempt
	vendor.Is1099Vendor = req.Is1099Vendor
	vendor.Email = req.Email
	vendor.Phone = req.Phone
	vendor.Fax = req.Fax
	vendor.Website = req.Website
	vendor.AddressLine1 = req.AddressLine1
	vendor.AddressLine2 = req.AddressLine2
	vendor.City = req.City
	vendor.StateProvince = req.StateProvince
	vendor.PostalCode = req.PostalCode
	vendor.Country = strings.ToUpper(req.Country)
	vendor.PaymentTerms = req.PaymentTerms
	vendor.PaymentMethod = req.PaymentMethod
	vendor.Currency = strings.ToUpper(req.Currency)
	vendor.CreditLimit = req.CreditLimit
	vendor.BankName = req.BankName
	vendor.BankAccountNumber = req.BankAccountNumber
	vendor.BankRoutingNumber = req.BankRoutingNumber
	vendor.SwiftCode = req.SwiftCode
	vendor.IBAN = req.IBAN
	vendor.Notes = req.Notes
	vendor.Tags = req.Tags
	vendor.UpdatedBy = &req.UpdatedBy

	if err := s.vendorRepo.Update(ctx, vendor); err != nil {
		return nil, err
	}

	s.log.Info().
		Str("vendor_id", vendor.ID).
		Str("vendor_code", vendor.VendorCode).
		Msg("Vendor updated")

	return vendor, nil
}

// DeleteVendor deletes a vendor
func (s *VendorService) DeleteVendor(ctx context.Context, id, entityID string) error {
	// TODO: Check if vendor has invoices (when invoice service is implemented)

	if err := s.vendorRepo.Delete(ctx, id, entityID); err != nil {
		return err
	}

	s.log.Info().
		Str("vendor_id", id).
		Str("entity_id", entityID).
		Msg("Vendor deleted")

	return nil
}

// ListVendors lists vendors with filtering and pagination
func (s *VendorService) ListVendors(ctx context.Context, entityID string, status, vendorType *string, activeOnly bool, page, pageSize int) ([]*repository.Vendor, int64, error) {
	offset := (page - 1) * pageSize
	return s.vendorRepo.List(ctx, entityID, status, vendorType, activeOnly, pageSize, offset)
}

// ActivateVendor activates a vendor
func (s *VendorService) ActivateVendor(ctx context.Context, id, entityID, updatedBy string) error {
	vendor, err := s.vendorRepo.GetByID(ctx, id, entityID)
	if err != nil {
		return err
	}

	vendor.Status = "active"
	vendor.UpdatedBy = &updatedBy

	if err := s.vendorRepo.Update(ctx, vendor); err != nil {
		return err
	}

	s.log.Info().
		Str("vendor_id", id).
		Str("entity_id", entityID).
		Msg("Vendor activated")

	return nil
}

// DeactivateVendor deactivates a vendor
func (s *VendorService) DeactivateVendor(ctx context.Context, id, entityID, updatedBy string) error {
	vendor, err := s.vendorRepo.GetByID(ctx, id, entityID)
	if err != nil {
		return err
	}

	// TODO: Check if vendor has pending invoices

	vendor.Status = "inactive"
	vendor.UpdatedBy = &updatedBy

	if err := s.vendorRepo.Update(ctx, vendor); err != nil {
		return err
	}

	s.log.Info().
		Str("vendor_id", id).
		Str("entity_id", entityID).
		Msg("Vendor deactivated")

	return nil
}

// GetVendorContacts retrieves all contacts for a vendor
func (s *VendorService) GetVendorContacts(ctx context.Context, vendorID string) ([]*repository.VendorContact, error) {
	return s.vendorRepo.GetContacts(ctx, vendorID)
}

// AddVendorContact adds a contact to a vendor
func (s *VendorService) AddVendorContact(ctx context.Context, req *AddContactRequest) (*repository.VendorContact, error) {
	// Validate contact type
	validTypes := map[string]bool{
		"primary":   true,
		"billing":   true,
		"shipping":  true,
		"technical": true,
		"other":     true,
	}
	contactType := strings.ToLower(req.ContactType)
	if !validTypes[contactType] {
		return nil, errors.InvalidInput("contact_type", "invalid contact type")
	}

	contact := &repository.VendorContact{
		VendorID:    req.VendorID,
		ContactType: contactType,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Title:       req.Title,
		Email:       req.Email,
		Phone:       req.Phone,
		Mobile:      req.Mobile,
		IsPrimary:   req.IsPrimary,
		Notes:       req.Notes,
	}

	if err := s.vendorRepo.AddContact(ctx, contact); err != nil {
		return nil, err
	}

	s.log.Info().
		Str("vendor_id", req.VendorID).
		Str("contact_id", contact.ID).
		Msg("Vendor contact added")

	return contact, nil
}

// GetPaymentTerms retrieves all active payment terms
func (s *VendorService) GetPaymentTerms(ctx context.Context) ([]*repository.PaymentTerm, error) {
	return s.vendorRepo.GetPaymentTerms(ctx)
}

// ValidateVendor validates if a vendor can be used for invoice creation
func (s *VendorService) ValidateVendor(ctx context.Context, vendorID, entityID string) (bool, string, error) {
	return s.vendorRepo.ValidateVendor(ctx, vendorID, entityID)
}
