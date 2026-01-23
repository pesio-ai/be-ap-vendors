# be-vendors-service (AP-1)

Vendor Management Service for Pesio Finance ERP - Manages vendor/supplier master data for accounts payable operations.

## Overview

This service implements vendor/supplier management functionality as defined in FPRD AP-1. It provides:

- **Vendor Master Data**: Complete vendor information management
- **Contact Management**: Multiple contacts per vendor with roles
- **Payment Terms**: Standardized payment terms (Net 30, 2/10 Net 30, etc.)
- **Credit Limits**: Track outstanding balances and credit limits
- **Tax Compliance**: W-9 tracking, 1099 vendor flagging, tax ID management
- **Banking Information**: Secure storage of vendor banking details for payments
- **Document Management**: References to vendor documents (contracts, W-9s, insurance)
- **Multi-entity Support**: Vendors isolated by entity

## Features

### Vendor Types
- **Supplier**: Goods suppliers
- **Contractor**: Independent contractors
- **Service Provider**: Service companies
- **Consultant**: Professional consultants
- **Utility**: Utility companies

### Vendor Status
- **Pending Approval**: Newly created, awaiting approval
- **Active**: Approved and can be used for transactions
- **Inactive**: Deactivated, cannot be used for new transactions
- **Suspended**: Temporarily suspended (payment issues, compliance)

### Payment Methods
- Check
- ACH (Automated Clearing House)
- Wire Transfer
- Credit Card
- Cash

### Contact Types
- **Primary**: Main contact person
- **Billing**: Billing/accounts receivable contact
- **Shipping**: Shipping/receiving contact
- **Technical**: Technical support contact
- **Other**: Custom contact types

### Business Rules
- Vendor codes must be unique within an entity
- New vendors created with "pending_approval" status
- Active vendors can be used for invoice creation
- Credit limit enforcement (if set)
- Current balance tracked (updated by AP-2 invoices service)
- Country codes must be 2-letter ISO (e.g., "US")
- Currency codes must be 3-letter ISO (e.g., "USD")

## API Endpoints

### Health Check
```
GET /health
```

### Vendor Operations

#### List Vendors
```
GET /api/v1/vendors?entity_id={uuid}&status={status}&vendor_type={type}&active_only={bool}&page={int}&page_size={int}
```
**Query Parameters**:
- `entity_id` (required): Entity UUID
- `status` (optional): Filter by active, inactive, suspended, pending_approval
- `vendor_type` (optional): Filter by supplier, contractor, service_provider, consultant, utility
- `active_only` (optional): true/false, default false
- `page` (optional): Page number, default 1
- `page_size` (optional): Items per page, default 50, max 100

**Response**:
```json
{
  "vendors": [
    {
      "id": "uuid",
      "entity_id": "uuid",
      "vendor_code": "VENDOR001",
      "vendor_name": "Acme Corporation",
      "legal_name": "Acme Corporation Inc.",
      "vendor_type": "supplier",
      "status": "active",
      "tax_id": "12-3456789",
      "is_tax_exempt": false,
      "is_1099_vendor": false,
      "email": "ap@acme.com",
      "phone": "+1-555-123-4567",
      "address_line1": "123 Main St",
      "city": "New York",
      "state_province": "NY",
      "postal_code": "10001",
      "country": "US",
      "payment_terms": "NET30",
      "payment_method": "ach",
      "currency": "USD",
      "credit_limit": 5000000,
      "current_balance": 125000,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 142,
  "page": 1,
  "pageSize": 50
}
```

#### Get Vendor by ID
```
GET /api/v1/vendors/get?id={uuid}&entity_id={uuid}
```

#### Get Vendor by Code
```
GET /api/v1/vendors/code?vendor_code={code}&entity_id={uuid}
```

#### Create Vendor
```
POST /api/v1/vendors
Content-Type: application/json

{
  "entity_id": "uuid",
  "vendor_code": "VENDOR001",
  "vendor_name": "Acme Corporation",
  "legal_name": "Acme Corporation Inc.",
  "vendor_type": "supplier",
  "tax_id": "12-3456789",
  "is_tax_exempt": false,
  "is_1099_vendor": false,
  "email": "ap@acme.com",
  "phone": "+1-555-123-4567",
  "address_line1": "123 Main St",
  "city": "New York",
  "state_province": "NY",
  "postal_code": "10001",
  "country": "US",
  "payment_terms": "NET30",
  "payment_method": "ach",
  "currency": "USD",
  "credit_limit": 5000000,
  "bank_name": "Chase Bank",
  "bank_account_number": "123456789",
  "bank_routing_number": "021000021",
  "notes": "Preferred supplier for office supplies",
  "tags": ["office-supplies", "preferred"]
}
```

**Business Rules**:
- Creates vendor in `pending_approval` status
- Vendor code converted to uppercase
- Country code converted to uppercase
- Currency code converted to uppercase

#### Update Vendor
```
PUT /api/v1/vendors/update
Content-Type: application/json

{
  "id": "uuid",
  "entity_id": "uuid",
  "vendor_code": "VENDOR001",
  "vendor_name": "Acme Corporation",
  "status": "active",
  ...
}
```

#### Delete Vendor
```
DELETE /api/v1/vendors/delete?id={uuid}&entity_id={uuid}
```

**Business Rules**:
- Cannot delete vendors with invoices (when AP-2 is implemented)
- Permanently deletes vendor and all related contacts/documents

#### Validate Vendor
```
GET /api/v1/vendors/validate?id={uuid}&entity_id={uuid}
```

**Response**:
```json
{
  "valid": true,
  "message": ""
}
```

**Validation Rules**:
- Vendor must be in "active" status
- If credit limit set, current balance must not exceed limit
- Used by AP-2 (invoices service) before creating invoices

### Contact Operations

#### Get Vendor Contacts
```
GET /api/v1/vendors/contacts?vendor_id={uuid}
```

**Response**:
```json
{
  "contacts": [
    {
      "id": "uuid",
      "vendor_id": "uuid",
      "contact_type": "primary",
      "first_name": "John",
      "last_name": "Smith",
      "title": "Accounts Receivable Manager",
      "email": "john.smith@acme.com",
      "phone": "+1-555-123-4567",
      "mobile": "+1-555-987-6543",
      "is_primary": true,
      "notes": "Preferred contact for payment inquiries",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

#### Add Vendor Contact
```
POST /api/v1/vendors/contacts
Content-Type: application/json

{
  "vendor_id": "uuid",
  "contact_type": "primary",
  "first_name": "John",
  "last_name": "Smith",
  "title": "Accounts Receivable Manager",
  "email": "john.smith@acme.com",
  "phone": "+1-555-123-4567",
  "is_primary": true
}
```

### Payment Terms

#### Get Payment Terms
```
GET /api/v1/payment-terms
```

**Response**:
```json
{
  "payment_terms": [
    {
      "id": "uuid",
      "code": "NET30",
      "description": "Net 30 days",
      "net_days": 30,
      "discount_percent": null,
      "discount_days": null,
      "is_active": true
    },
    {
      "id": "uuid",
      "code": "2/10N30",
      "description": "2% 10 days, Net 30",
      "net_days": 30,
      "discount_percent": 2.00,
      "discount_days": 10,
      "is_active": true
    }
  ]
}
```

**Predefined Terms**:
- NET30 - Net 30 days
- NET60 - Net 60 days
- NET90 - Net 90 days
- 2/10N30 - 2% discount if paid within 10 days, otherwise net 30
- 1/10N30 - 1% discount if paid within 10 days, otherwise net 30
- DUE - Due on receipt
- COD - Cash on delivery
- CIA - Cash in advance

## Database Schema

### Tables

#### vendors
- `id` (UUID, PK): Vendor identifier
- `entity_id` (UUID): Entity this vendor belongs to
- `vendor_code` (VARCHAR): Unique vendor code within entity
- `vendor_name` (VARCHAR): Vendor display name
- `legal_name` (VARCHAR): Legal business name
- `vendor_type` (ENUM): supplier, contractor, service_provider, consultant, utility
- `status` (ENUM): active, inactive, suspended, pending_approval
- `tax_id` (VARCHAR): Tax identification number (EIN, SSN)
- `is_tax_exempt` (BOOLEAN): Tax exempt status
- `is_1099_vendor` (BOOLEAN): Receives 1099 form (US)
- Contact fields: email, phone, fax, website
- Address fields: address_line1, address_line2, city, state_province, postal_code, country
- Payment fields: payment_terms, payment_method, currency, credit_limit, current_balance
- Banking fields: bank_name, bank_account_number, bank_routing_number, swift_code, iban
- Metadata: notes, tags (array)
- Audit fields: created_by, created_at, updated_by, updated_at

**Constraints**:
- `vendors_entity_code_unique`: Unique(entity_id, vendor_code)
- `vendors_credit_limit_check`: credit_limit >= 0 (if set)
- `vendors_current_balance_check`: current_balance >= 0

#### vendor_contacts
- `id` (UUID, PK): Contact identifier
- `vendor_id` (UUID, FK): Parent vendor
- `contact_type` (ENUM): primary, billing, shipping, technical, other
- `first_name`, `last_name`, `title`: Contact person info
- Contact details: email, phone, mobile
- `is_primary` (BOOLEAN): Primary contact flag
- `notes` (TEXT): Additional notes
- Audit fields: created_at, updated_at

**Constraints**:
- Cascading delete when parent vendor deleted

#### vendor_documents
- `id` (UUID, PK): Document reference identifier
- `vendor_id` (UUID, FK): Parent vendor
- `document_type` (VARCHAR): W9, contract, insurance, etc.
- `document_name` (VARCHAR): File name
- `document_url` (TEXT): S3/storage URL
- File metadata: file_size, mime_type
- `expiration_date` (DATE): Document expiration (insurance, certifications)
- Upload tracking: uploaded_by, uploaded_at

**Constraints**:
- Cascading delete when parent vendor deleted

#### payment_terms
- `id` (UUID, PK): Term identifier
- `code` (VARCHAR): Unique code (e.g., "NET30")
- `description` (VARCHAR): Human-readable description
- `net_days` (INT): Number of days until payment due
- `discount_percent` (NUMERIC): Early payment discount percentage
- `discount_days` (INT): Days within which discount applies
- `is_active` (BOOLEAN): Active status
- `created_at` (TIMESTAMP): Creation timestamp

## Configuration

### Environment Variables

```bash
# Service Configuration
SERVICE_NAME=be-vendors-service
SERVICE_VERSION=1.0.0
ENVIRONMENT=development
LOG_LEVEL=info

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=pesio
DB_PASSWORD=dev_password_change_me
DB_NAME=ap_vendors_db
DB_SSL_MODE=disable
DB_MAX_CONNS=25
DB_MIN_CONNS=5

# Server Configuration
SERVER_PORT=8084
```

Copy `.env.example` to `.env` and update values for your environment.

## Dependencies

- **be-go-common**: Shared libraries for config, database, logging, errors, middleware
- **PostgreSQL 16**: Database with pgvector extension
- **Go 1.23+**: Required Go version

## Running Locally

### Prerequisites
- Docker and Docker Compose (for infrastructure)
- Go 1.23+
- PostgreSQL 16 with pgvector extension

### Start Infrastructure
```bash
cd ../infrastructure
make up
```

### Run Service
```bash
# Install dependencies
go mod download

# Run database migrations
psql -h localhost -U pesio -d ap_vendors_db -f migrations/001_initial_schema.sql

# Start service
go run cmd/server/main.go
```

Service will start on port 8084. Health check: http://localhost:8084/health

## Development

### Project Structure
```
be-vendors-service/
├── cmd/
│   └── server/
│       └── main.go                 # Server entry point
├── internal/
│   ├── handler/
│   │   └── http_handler.go         # HTTP REST handlers
│   ├── repository/
│   │   └── vendor_repository.go    # Data access layer
│   └── service/
│       └── vendor_service.go       # Business logic
├── migrations/
│   └── 001_initial_schema.sql      # Database schema
├── .env.example                     # Environment template
├── go.mod                           # Go dependencies
└── README.md                        # This file
```

### Testing
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/service
```

## Integration with Other Services

### be-identity-service (PLT-1)
- Future: JWT authentication for all endpoints
- User ID extraction from token for created_by/updated_by fields

### be-entity-service (PLT-2)
- Validates entity_id references valid entities
- Future: Entity hierarchy permissions

### be-invoices-service (AP-2)
- **CRITICAL**: Validates vendor before invoice creation
- **CRITICAL**: Updates current_balance when invoices created/paid
- Retrieves vendor payment terms
- Retrieves vendor banking details for payment processing

## TODO: Integration Points

The following integration points are marked as TODO in the code:

1. **Invoice Check on Delete** (service/vendor_service.go)
   - Check if vendor has invoices before deletion
   - Call AP-2 to verify no pending/posted invoices

2. **Balance Updates** (repository/vendor_repository.go)
   - current_balance field updated by AP-2 (invoices service)
   - Incremented when invoices created
   - Decremented when payments posted

3. **User Authentication** (handler/http_handler.go)
   - Extract user ID from JWT token
   - Replace "system" with actual user ID
   - For created_by, updated_by fields

## License

Proprietary - Pesio AI

## Support

For questions or issues, contact the development team or file an issue in the repository.
