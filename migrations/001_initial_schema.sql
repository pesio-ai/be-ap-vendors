-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Vendor type enum
CREATE TYPE vendor_type AS ENUM ('supplier', 'contractor', 'service_provider', 'consultant', 'utility');

-- Vendor status enum
CREATE TYPE vendor_status AS ENUM ('active', 'inactive', 'suspended', 'pending_approval');

-- Payment method enum
CREATE TYPE payment_method AS ENUM ('check', 'ach', 'wire', 'credit_card', 'cash');

-- Contact type enum
CREATE TYPE contact_type AS ENUM ('primary', 'billing', 'shipping', 'technical', 'other');

-- Vendors (Master)
CREATE TABLE vendors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_id UUID NOT NULL,
    vendor_code VARCHAR(50) NOT NULL,
    vendor_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255),
    vendor_type vendor_type NOT NULL,
    status vendor_status NOT NULL DEFAULT 'pending_approval',
    tax_id VARCHAR(50),
    is_tax_exempt BOOLEAN NOT NULL DEFAULT FALSE,
    is_1099_vendor BOOLEAN NOT NULL DEFAULT FALSE,

    -- Contact Information
    email VARCHAR(255),
    phone VARCHAR(50),
    fax VARCHAR(50),
    website VARCHAR(255),

    -- Address
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state_province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2) NOT NULL DEFAULT 'US',

    -- Payment Terms
    payment_terms VARCHAR(50) NOT NULL DEFAULT 'Net 30',
    payment_method payment_method,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    credit_limit BIGINT,  -- Amount in smallest currency unit
    current_balance BIGINT NOT NULL DEFAULT 0,  -- Outstanding balance

    -- Banking Information
    bank_name VARCHAR(255),
    bank_account_number VARCHAR(100),
    bank_routing_number VARCHAR(50),
    swift_code VARCHAR(20),
    iban VARCHAR(50),

    -- Metadata
    notes TEXT,
    tags TEXT[],

    -- Audit fields
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_by UUID,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT vendors_entity_code_unique UNIQUE (entity_id, vendor_code),
    CONSTRAINT vendors_credit_limit_check CHECK (credit_limit IS NULL OR credit_limit >= 0),
    CONSTRAINT vendors_current_balance_check CHECK (current_balance >= 0)
);

-- Vendor Contacts
CREATE TABLE vendor_contacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    contact_type contact_type NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    title VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    mobile VARCHAR(50),
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT vendor_contacts_vendor_type_unique UNIQUE (vendor_id, contact_type, email)
);

-- Vendor Documents (references to document storage)
CREATE TABLE vendor_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vendor_id UUID NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,  -- W9, contract, insurance, etc.
    document_name VARCHAR(255) NOT NULL,
    document_url TEXT NOT NULL,  -- S3/storage URL
    file_size BIGINT,
    mime_type VARCHAR(100),
    expiration_date DATE,
    uploaded_by UUID,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT vendor_documents_vendor_type_unique UNIQUE (vendor_id, document_type, document_name)
);

-- Vendor Payment Terms (predefined common terms)
CREATE TABLE payment_terms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(20) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL,
    net_days INT NOT NULL,
    discount_percent NUMERIC(5, 2),
    discount_days INT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_vendors_entity_id ON vendors(entity_id);
CREATE INDEX idx_vendors_vendor_code ON vendors(vendor_code);
CREATE INDEX idx_vendors_vendor_name ON vendors(vendor_name);
CREATE INDEX idx_vendors_status ON vendors(status);
CREATE INDEX idx_vendors_vendor_type ON vendors(vendor_type);
CREATE INDEX idx_vendors_tax_id ON vendors(tax_id);
CREATE INDEX idx_vendor_contacts_vendor_id ON vendor_contacts(vendor_id);
CREATE INDEX idx_vendor_documents_vendor_id ON vendor_documents(vendor_id);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER trigger_vendors_updated_at
BEFORE UPDATE ON vendors
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_vendor_contacts_updated_at
BEFORE UPDATE ON vendor_contacts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Insert common payment terms
INSERT INTO payment_terms (code, description, net_days, discount_percent, discount_days) VALUES
    ('NET30', 'Net 30 days', 30, NULL, NULL),
    ('NET60', 'Net 60 days', 60, NULL, NULL),
    ('NET90', 'Net 90 days', 90, NULL, NULL),
    ('2/10N30', '2% 10 days, Net 30', 30, 2.00, 10),
    ('1/10N30', '1% 10 days, Net 30', 30, 1.00, 10),
    ('DUE', 'Due on receipt', 0, NULL, NULL),
    ('COD', 'Cash on delivery', 0, NULL, NULL),
    ('CIA', 'Cash in advance', 0, NULL, NULL);

-- Comments
COMMENT ON TABLE vendors IS 'Vendor/supplier master data for accounts payable';
COMMENT ON TABLE vendor_contacts IS 'Contact persons for vendors';
COMMENT ON TABLE vendor_documents IS 'Vendor-related documents (W9, contracts, insurance certificates)';
COMMENT ON TABLE payment_terms IS 'Standard payment terms lookup table';
COMMENT ON COLUMN vendors.current_balance IS 'Total outstanding balance owed to vendor (updated by AP-2 invoices service)';
COMMENT ON COLUMN vendors.credit_limit IS 'Maximum credit allowed with this vendor';
COMMENT ON COLUMN vendors.is_1099_vendor IS 'US tax reporting: vendor receives 1099-MISC/NEC form';
