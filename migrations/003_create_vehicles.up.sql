-- Migration: create table vehicles
CREATE TABLE vehicles (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    plate_no VARCHAR UNIQUE NOT NULL,
    model VARCHAR,
    brand VARCHAR,
    color VARCHAR,
    year INTEGER,
    is_active INTEGER DEFAULT 1,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Foreign Key Constraints
    CONSTRAINT fk_vehicles_customer 
        FOREIGN KEY (customer_id) REFERENCES customers(id) 
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_vehicles_created_by 
        FOREIGN KEY (created_by) REFERENCES users(id) 
        ON DELETE SET NULL ON UPDATE CASCADE
);

-- Vehicles Table - Most Critical Indexes
-- Note: plate_no already has UNIQUE constraint which creates an index automatically
CREATE INDEX idx_vehicles_customer_active ON vehicles(customer_id, is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicles_brand_model ON vehicles(brand, model) WHERE deleted_at IS NULL;