-- Migration: create table users

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    sex SMALLINT NOT NULL, -- 1=male, 2=female
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL,
    is_active INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

-- Users Table - Minimal Essential Indexes
-- Note: email already has UNIQUE constraint which creates an index automatically
CREATE INDEX idx_users_role_active ON users(role, is_active) WHERE deleted_at IS NULL;
-- Insert Data Admin
-- Password : password123
INSERT INTO users (name, sex, email, password, role, is_active)
VALUES
('Aldino Pratama Bagaskara', 1, 'emailku1@gmail.com', '$2a$04$n8Ps2Wy9Jf5/7Mc14iK2P.kryqWJeY2AFMCGQW7cl3wumFpR9yBRi', 'ADMIN', 1), 
('Bagaskara', 1, 'emailku2@gmail.com', '$2a$04$n8Ps2Wy9Jf5/7Mc14iK2P.kryqWJeY2AFMCGQW7cl3wumFpR9yBRi', 'OPERATOR', 1),
('Aldino', 1, 'emailku3@gmail.com', '$2a$04$n8Ps2Wy9Jf5/7Mc14iK2P.kryqWJeY2AFMCGQW7cl3wumFpR9yBRi', 'USER', 1);