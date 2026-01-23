package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pesio-ai/be-go-common/database"
	"github.com/pesio-ai/be-go-common/errors"
)

// Vendor represents a vendor/supplier
type Vendor struct {
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
	CurrentBalance    int64
	BankName          *string
	BankAccountNumber *string
	BankRoutingNumber *string
	SwiftCode         *string
	IBAN              *string
	Notes             *string
	Tags              []string
	CreatedBy         *string
	CreatedAt         string
	UpdatedBy         *string
	UpdatedAt         string
}

// VendorContact represents a vendor contact person
type VendorContact struct {
	ID          string
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
	CreatedAt   string
	UpdatedAt   string
}

// VendorDocument represents a vendor document reference
type VendorDocument struct {
	ID             string
	VendorID       string
	DocumentType   string
	DocumentName   string
	DocumentURL    string
	FileSize       *int64
	MimeType       *string
	ExpirationDate *string
	UploadedBy     *string
	UploadedAt     string
}

// PaymentTerm represents payment terms
type PaymentTerm struct {
	ID              string
	Code            string
	Description     string
	NetDays         int
	DiscountPercent *float64
	DiscountDays    *int
	IsActive        bool
	CreatedAt       string
}

// VendorRepository handles vendor data operations
type VendorRepository struct {
	db *database.DB
}

// NewVendorRepository creates a new vendor repository
func NewVendorRepository(db *database.DB) *VendorRepository {
	return &VendorRepository{db: db}
}

// Create creates a new vendor
func (r *VendorRepository) Create(ctx context.Context, vendor *Vendor) error {
	query := `
		INSERT INTO vendors (entity_id, vendor_code, vendor_name, legal_name, vendor_type,
		                     status, tax_id, is_tax_exempt, is_1099_vendor,
		                     email, phone, fax, website,
		                     address_line1, address_line2, city, state_province, postal_code, country,
		                     payment_terms, payment_method, currency, credit_limit,
		                     bank_name, bank_account_number, bank_routing_number, swift_code, iban,
		                     notes, tags, created_by)
		VALUES ($1, $2, $3, $4, $5::vendor_type, $6::vendor_status, $7, $8, $9,
		        $10, $11, $12, $13,
		        $14, $15, $16, $17, $18, $19,
		        $20, $21::payment_method, $22, $23,
		        $24, $25, $26, $27, $28,
		        $29, $30, $31)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		vendor.EntityID,
		vendor.VendorCode,
		vendor.VendorName,
		vendor.LegalName,
		vendor.VendorType,
		vendor.Status,
		vendor.TaxID,
		vendor.IsTaxExempt,
		vendor.Is1099Vendor,
		vendor.Email,
		vendor.Phone,
		vendor.Fax,
		vendor.Website,
		vendor.AddressLine1,
		vendor.AddressLine2,
		vendor.City,
		vendor.StateProvince,
		vendor.PostalCode,
		vendor.Country,
		vendor.PaymentTerms,
		vendor.PaymentMethod,
		vendor.Currency,
		vendor.CreditLimit,
		vendor.BankName,
		vendor.BankAccountNumber,
		vendor.BankRoutingNumber,
		vendor.SwiftCode,
		vendor.IBAN,
		vendor.Notes,
		vendor.Tags,
		vendor.CreatedBy,
	).Scan(&vendor.ID, &vendor.CreatedAt, &vendor.UpdatedAt)

	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "failed to create vendor")
	}

	return nil
}

// GetByID retrieves a vendor by ID
func (r *VendorRepository) GetByID(ctx context.Context, id, entityID string) (*Vendor, error) {
	vendor := &Vendor{}

	query := `
		SELECT id, entity_id, vendor_code, vendor_name, legal_name, vendor_type,
		       status, tax_id, is_tax_exempt, is_1099_vendor,
		       email, phone, fax, website,
		       address_line1, address_line2, city, state_province, postal_code, country,
		       payment_terms, payment_method, currency, credit_limit, current_balance,
		       bank_name, bank_account_number, bank_routing_number, swift_code, iban,
		       notes, tags,
		       created_by, created_at, updated_by, updated_at
		FROM vendors
		WHERE id = $1 AND entity_id = $2
	`

	err := r.db.QueryRow(ctx, query, id, entityID).Scan(
		&vendor.ID,
		&vendor.EntityID,
		&vendor.VendorCode,
		&vendor.VendorName,
		&vendor.LegalName,
		&vendor.VendorType,
		&vendor.Status,
		&vendor.TaxID,
		&vendor.IsTaxExempt,
		&vendor.Is1099Vendor,
		&vendor.Email,
		&vendor.Phone,
		&vendor.Fax,
		&vendor.Website,
		&vendor.AddressLine1,
		&vendor.AddressLine2,
		&vendor.City,
		&vendor.StateProvince,
		&vendor.PostalCode,
		&vendor.Country,
		&vendor.PaymentTerms,
		&vendor.PaymentMethod,
		&vendor.Currency,
		&vendor.CreditLimit,
		&vendor.CurrentBalance,
		&vendor.BankName,
		&vendor.BankAccountNumber,
		&vendor.BankRoutingNumber,
		&vendor.SwiftCode,
		&vendor.IBAN,
		&vendor.Notes,
		&vendor.Tags,
		&vendor.CreatedBy,
		&vendor.CreatedAt,
		&vendor.UpdatedBy,
		&vendor.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("vendor", id)
	}
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to get vendor")
	}

	return vendor, nil
}

// GetByCode retrieves a vendor by vendor code
func (r *VendorRepository) GetByCode(ctx context.Context, code, entityID string) (*Vendor, error) {
	vendor := &Vendor{}

	query := `
		SELECT id, entity_id, vendor_code, vendor_name, legal_name, vendor_type,
		       status, tax_id, is_tax_exempt, is_1099_vendor,
		       email, phone, fax, website,
		       address_line1, address_line2, city, state_province, postal_code, country,
		       payment_terms, payment_method, currency, credit_limit, current_balance,
		       bank_name, bank_account_number, bank_routing_number, swift_code, iban,
		       notes, tags,
		       created_by, created_at, updated_by, updated_at
		FROM vendors
		WHERE vendor_code = $1 AND entity_id = $2
	`

	err := r.db.QueryRow(ctx, query, code, entityID).Scan(
		&vendor.ID,
		&vendor.EntityID,
		&vendor.VendorCode,
		&vendor.VendorName,
		&vendor.LegalName,
		&vendor.VendorType,
		&vendor.Status,
		&vendor.TaxID,
		&vendor.IsTaxExempt,
		&vendor.Is1099Vendor,
		&vendor.Email,
		&vendor.Phone,
		&vendor.Fax,
		&vendor.Website,
		&vendor.AddressLine1,
		&vendor.AddressLine2,
		&vendor.City,
		&vendor.StateProvince,
		&vendor.PostalCode,
		&vendor.Country,
		&vendor.PaymentTerms,
		&vendor.PaymentMethod,
		&vendor.Currency,
		&vendor.CreditLimit,
		&vendor.CurrentBalance,
		&vendor.BankName,
		&vendor.BankAccountNumber,
		&vendor.BankRoutingNumber,
		&vendor.SwiftCode,
		&vendor.IBAN,
		&vendor.Notes,
		&vendor.Tags,
		&vendor.CreatedBy,
		&vendor.CreatedAt,
		&vendor.UpdatedBy,
		&vendor.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("vendor", code)
	}
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to get vendor by code")
	}

	return vendor, nil
}

// Update updates a vendor
func (r *VendorRepository) Update(ctx context.Context, vendor *Vendor) error {
	query := `
		UPDATE vendors
		SET vendor_code = $3, vendor_name = $4, legal_name = $5, vendor_type = $6::vendor_type,
		    status = $7::vendor_status, tax_id = $8, is_tax_exempt = $9, is_1099_vendor = $10,
		    email = $11, phone = $12, fax = $13, website = $14,
		    address_line1 = $15, address_line2 = $16, city = $17, state_province = $18,
		    postal_code = $19, country = $20,
		    payment_terms = $21, payment_method = $22::payment_method, currency = $23, credit_limit = $24,
		    bank_name = $25, bank_account_number = $26, bank_routing_number = $27,
		    swift_code = $28, iban = $29,
		    notes = $30, tags = $31, updated_by = $32, updated_at = NOW()
		WHERE id = $1 AND entity_id = $2
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		vendor.ID,
		vendor.EntityID,
		vendor.VendorCode,
		vendor.VendorName,
		vendor.LegalName,
		vendor.VendorType,
		vendor.Status,
		vendor.TaxID,
		vendor.IsTaxExempt,
		vendor.Is1099Vendor,
		vendor.Email,
		vendor.Phone,
		vendor.Fax,
		vendor.Website,
		vendor.AddressLine1,
		vendor.AddressLine2,
		vendor.City,
		vendor.StateProvince,
		vendor.PostalCode,
		vendor.Country,
		vendor.PaymentTerms,
		vendor.PaymentMethod,
		vendor.Currency,
		vendor.CreditLimit,
		vendor.BankName,
		vendor.BankAccountNumber,
		vendor.BankRoutingNumber,
		vendor.SwiftCode,
		vendor.IBAN,
		vendor.Notes,
		vendor.Tags,
		vendor.UpdatedBy,
	).Scan(&vendor.UpdatedAt)

	if err == pgx.ErrNoRows {
		return errors.NotFound("vendor", vendor.ID)
	}
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "failed to update vendor")
	}

	return nil
}

// Delete deletes a vendor
func (r *VendorRepository) Delete(ctx context.Context, id, entityID string) error {
	query := `DELETE FROM vendors WHERE id = $1 AND entity_id = $2`

	tag, err := r.db.Exec(ctx, query, id, entityID)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "failed to delete vendor")
	}

	if tag.RowsAffected() == 0 {
		return errors.NotFound("vendor", id)
	}

	return nil
}

// List retrieves vendors with filtering and pagination
func (r *VendorRepository) List(ctx context.Context, entityID string, status, vendorType *string, activeOnly bool, limit, offset int) ([]*Vendor, int64, error) {
	query := `
		SELECT id, entity_id, vendor_code, vendor_name, legal_name, vendor_type,
		       status, tax_id, is_tax_exempt, is_1099_vendor,
		       email, phone, fax, website,
		       address_line1, address_line2, city, state_province, postal_code, country,
		       payment_terms, payment_method, currency, credit_limit, current_balance,
		       bank_name, bank_account_number, bank_routing_number, swift_code, iban,
		       notes, tags,
		       created_by, created_at, updated_by, updated_at
		FROM vendors
		WHERE entity_id = $1
	`

	countQuery := `SELECT COUNT(*) FROM vendors WHERE entity_id = $1`

	args := []interface{}{entityID}
	argCount := 2

	if status != nil {
		query += fmt.Sprintf(" AND status = $%d::vendor_status", argCount)
		countQuery += fmt.Sprintf(" AND status = $%d::vendor_status", argCount)
		args = append(args, *status)
		argCount++
	}

	if vendorType != nil {
		query += fmt.Sprintf(" AND vendor_type = $%d::vendor_type", argCount)
		countQuery += fmt.Sprintf(" AND vendor_type = $%d::vendor_type", argCount)
		args = append(args, *vendorType)
		argCount++
	}

	if activeOnly {
		query += fmt.Sprintf(" AND status = $%d::vendor_status", argCount)
		countQuery += fmt.Sprintf(" AND status = $%d::vendor_status", argCount)
		args = append(args, "active")
		argCount++
	}

	query += " ORDER BY vendor_name"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)

	queryArgs := append(args, limit, offset)

	// Get total count
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrCodeInternal, "failed to count vendors")
	}

	// Get vendors
	rows, err := r.db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrCodeInternal, "failed to list vendors")
	}
	defer rows.Close()

	vendors := make([]*Vendor, 0)
	for rows.Next() {
		vendor := &Vendor{}
		err := rows.Scan(
			&vendor.ID,
			&vendor.EntityID,
			&vendor.VendorCode,
			&vendor.VendorName,
			&vendor.LegalName,
			&vendor.VendorType,
			&vendor.Status,
			&vendor.TaxID,
			&vendor.IsTaxExempt,
			&vendor.Is1099Vendor,
			&vendor.Email,
			&vendor.Phone,
			&vendor.Fax,
			&vendor.Website,
			&vendor.AddressLine1,
			&vendor.AddressLine2,
			&vendor.City,
			&vendor.StateProvince,
			&vendor.PostalCode,
			&vendor.Country,
			&vendor.PaymentTerms,
			&vendor.PaymentMethod,
			&vendor.Currency,
			&vendor.CreditLimit,
			&vendor.CurrentBalance,
			&vendor.BankName,
			&vendor.BankAccountNumber,
			&vendor.BankRoutingNumber,
			&vendor.SwiftCode,
			&vendor.IBAN,
			&vendor.Notes,
			&vendor.Tags,
			&vendor.CreatedBy,
			&vendor.CreatedAt,
			&vendor.UpdatedBy,
			&vendor.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, errors.ErrCodeInternal, "failed to scan vendor")
		}

		vendors = append(vendors, vendor)
	}

	return vendors, total, nil
}

// GetContacts retrieves all contacts for a vendor
func (r *VendorRepository) GetContacts(ctx context.Context, vendorID string) ([]*VendorContact, error) {
	query := `
		SELECT id, vendor_id, contact_type, first_name, last_name, title,
		       email, phone, mobile, is_primary, notes,
		       created_at, updated_at
		FROM vendor_contacts
		WHERE vendor_id = $1
		ORDER BY is_primary DESC, first_name, last_name
	`

	rows, err := r.db.Query(ctx, query, vendorID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to get vendor contacts")
	}
	defer rows.Close()

	contacts := make([]*VendorContact, 0)
	for rows.Next() {
		contact := &VendorContact{}
		err := rows.Scan(
			&contact.ID,
			&contact.VendorID,
			&contact.ContactType,
			&contact.FirstName,
			&contact.LastName,
			&contact.Title,
			&contact.Email,
			&contact.Phone,
			&contact.Mobile,
			&contact.IsPrimary,
			&contact.Notes,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to scan vendor contact")
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// AddContact adds a contact to a vendor
func (r *VendorRepository) AddContact(ctx context.Context, contact *VendorContact) error {
	query := `
		INSERT INTO vendor_contacts (vendor_id, contact_type, first_name, last_name, title,
		                             email, phone, mobile, is_primary, notes)
		VALUES ($1, $2::contact_type, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		contact.VendorID,
		contact.ContactType,
		contact.FirstName,
		contact.LastName,
		contact.Title,
		contact.Email,
		contact.Phone,
		contact.Mobile,
		contact.IsPrimary,
		contact.Notes,
	).Scan(&contact.ID, &contact.CreatedAt, &contact.UpdatedAt)

	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "failed to add vendor contact")
	}

	return nil
}

// GetPaymentTerms retrieves all active payment terms
func (r *VendorRepository) GetPaymentTerms(ctx context.Context) ([]*PaymentTerm, error) {
	query := `
		SELECT id, code, description, net_days, discount_percent, discount_days, is_active, created_at
		FROM payment_terms
		WHERE is_active = TRUE
		ORDER BY net_days
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to get payment terms")
	}
	defer rows.Close()

	terms := make([]*PaymentTerm, 0)
	for rows.Next() {
		term := &PaymentTerm{}
		err := rows.Scan(
			&term.ID,
			&term.Code,
			&term.Description,
			&term.NetDays,
			&term.DiscountPercent,
			&term.DiscountDays,
			&term.IsActive,
			&term.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to scan payment term")
		}

		terms = append(terms, term)
	}

	return terms, nil
}

// ValidateVendor validates if a vendor can be used for invoice creation
func (r *VendorRepository) ValidateVendor(ctx context.Context, vendorID, entityID string) (bool, string, error) {
	vendor, err := r.GetByID(ctx, vendorID, entityID)
	if err != nil {
		return false, "vendor not found", err
	}

	if vendor.Status != "active" {
		return false, fmt.Sprintf("vendor status is '%s', must be active", vendor.Status), nil
	}

	// Check credit limit if set
	if vendor.CreditLimit != nil && vendor.CurrentBalance >= *vendor.CreditLimit {
		return false, fmt.Sprintf("vendor has exceeded credit limit: balance=%d, limit=%d",
			vendor.CurrentBalance, *vendor.CreditLimit), nil
	}

	return true, "", nil
}

// UpdateBalance updates the vendor's current balance
func (r *VendorRepository) UpdateBalance(ctx context.Context, vendorID, entityID string, amount int64) error {
	query := `
		UPDATE vendors
		SET current_balance = current_balance + $3,
		    updated_at = NOW()
		WHERE id = $1 AND entity_id = $2
		RETURNING id
	`

	var returnedID string
	err := r.db.QueryRow(ctx, query, vendorID, entityID, amount).Scan(&returnedID)

	if err == pgx.ErrNoRows {
		return errors.NotFound("vendor", vendorID)
	}
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "failed to update vendor balance")
	}

	return nil
}
