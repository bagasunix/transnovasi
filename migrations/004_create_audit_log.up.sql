-- Migration: create table audit_logs

CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255),
    method VARCHAR(50) NOT NULL,
    url TEXT NOT NULL,
    status_code VARCHAR(10),
    request TEXT,
    response TEXT,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index penting untuk performa query