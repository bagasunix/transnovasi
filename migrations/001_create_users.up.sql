-- Migration: create table users

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    sex SMALLINT NOT NULL, -- 1=male, 2=female
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    role_id SMALLINT NOT NULL,
    is_active INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

-- Users Table - Minimal Essential Indexes
-- Note: email already has UNIQUE constraint which creates an index automatically
CREATE INDEX idx_users_role_active ON users(role_id, is_active) WHERE deleted_at IS NULL;