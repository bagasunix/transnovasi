-- Migration: create table customers
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    address TEXT,
    password VARCHAR(100) NOT NULL,
    role_id VARCHAR NOT NULL,
    is_active INTEGER DEFAULT 0,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL

    -- Add Foreign Key Constraint for customers.created_by
CONSTRAINT fk_customers_created_by 
    FOREIGN KEY (created_by) REFERENCES users(id) 
    ON DELETE SET NULL ON UPDATE CASCADE
);

-- Customers Table - Essential Indexes
-- Note: email already has UNIQUE constraint which creates an index automatically
CREATE INDEX idx_customers_active_role ON customers(is_active, role_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_created_by ON customers(created_by);
