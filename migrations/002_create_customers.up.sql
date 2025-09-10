-- Migration: create table customers
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(14),
    address TEXT,
    is_active INTEGER DEFAULT 0,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    -- Add Foreign Key Constraint for customers.created_by
CONSTRAINT fk_customers_created_by 
    FOREIGN KEY (created_by) REFERENCES users(id) 
    ON DELETE SET NULL ON UPDATE CASCADE
);

-- Customers Table - Essential Indexes
-- Note: email already has UNIQUE constraint which creates an index automatically
CREATE INDEX idx_customers_created_by ON customers(created_by);
-- Insert Data
INSERT INTO customers (name, email, phone, address, is_active, created_by) VALUES
('Budi Santoso', 'budi.santoso@example.com', '+628112345678', 'Jl. Merdeka No. 123, Jakarta', 1, 2),
('Siti Rahayu', 'siti.rahayu@example.com', '+628123456789', 'Jl. Sudirman No. 45, Bandung', 1, 2),
('Ahmad Hidayat', 'ahmad.hidayat@example.com', '+628134567890', 'Jl. Gatot Subroto No. 67, Surabaya', 0, 2),
('Dewi Lestari', 'dewi.lestari@example.com', '+628145678901', 'Jl. Thamrin No. 89, Medan', 1, 2),
('Joko Widodo', 'joko.widodo@example.com', '+628156789012', 'Jl. Pemuda No. 56, Yogyakarta', 1, 2),
('Rina Wijaya', 'rina.wijaya@example.com', '+628167890123', 'Jl. Diponegoro No. 34, Semarang', 0, 2),
('Fajar Pratama', 'fajar.pratama@example.com', '+628178901234', 'Jl. Ahmad Yani No. 78, Malang', 1, 2),
('Maya Sari', 'maya.sari@example.com', '+628189012345', 'Jl. Pahlawan No. 12, Bali', 1, 2),
('Hendra Setiawan', 'hendra.setiawan@example.com', '+628190123456', 'Jl. Majapahit No. 23, Lombok', 0, 2),
('Lina Marlina', 'lina.marlina@example.com', '+628201234567', 'Jl. Asia Afrika No. 90, Bandung', 1, 2),
('Irfan Syah', 'irfan.syah@example.com', '+628211234567', 'Jl. Cihampelas No. 145, Bandung', 1, 2),
('Rudi Hermawan', 'rudi.hermawan@example.com', '+628221234567', 'Jl. Braga No. 67, Bandung', 0, 2),
('Sari Indah', 'sari.indah@example.com', '+628231234567', 'Jl. Dago No. 89, Bandung', 1, 2),
('Bambang Sutrisno', 'bambang.sutrisno@example.com', '+628241234567', 'Jl. Riau No. 56, Bandung', 1, 2),
('Dian Novita', 'dian.novita@example.com', '+628251234567', 'Jl. Aceh No. 34, Bandung', 0, 2),
('Eko Prasetyo', 'eko.prasetyo@example.com', '+628261234567', 'Jl. Sumatra No. 22, Bandung', 1, 2),
('Fitri Handayani', 'fitri.handayani@example.com', '+628271234567', 'Jl. Jawa No. 11, Bandung', 1, 2),
('Guntur Wibowo', 'guntur.wibowo@example.com', '+628281234567', 'Jl. Kalimantan No. 99, Bandung', 0, 2),
('Hana Saputri', 'hana.saputri@example.com', '+628291234567', 'Jl. Sulawesi No. 88, Bandung', 1, 2),
('Iwan Kurniawan', 'iwan.kurniawan@example.com', '+628301234567', 'Jl. Papua No. 77, Bandung', 1, 2),
('Juli Astuti', 'juli.astuti@example.com', '+628311234567', 'Jl. Maluku No. 66, Bandung', 0, 2),
('Khalid Rahman', 'khalid.rahman@example.com', '+628321234567', 'Jl. Nusa Tenggara No. 55, Bandung', 1, 2),
('Lia Agustina', 'lia.agustina@example.com', '+628331234567', 'Jl. Banten No. 44, Bandung', 1, 2),
('Muhammad Ali', 'muhammad.ali@example.com', '+628341234567', 'Jl. Jakarta No. 33, Bandung', 0, 2),
('Nina Soraya', 'nina.soraya@example.com', '+628351234567', 'Jl. Bogor No. 22, Bandung', 1, 2),
('Oki Setiawan', 'oki.setiawan@example.com', '+628361234567', 'Jl. Depok No. 11, Bandung', 1, 2),
('Putri Anggraini', 'putri.anggraini@example.com', '+628371234567', 'Jl. Tangerang No. 100, Bandung', 0, 2),
('Qomarudin', 'qomarudin@example.com', '+628381234567', 'Jl. Bekasi No. 200, Bandung', 1, 2),
('Rizky Ramadhan', 'rizky.ramadhan@example.com', '+628391234567', 'Jl. Cimahi No. 300, Bandung', 1, 2),
('Sinta Purwanti', 'sinta.purwanti@example.com', '+628401234567', 'Jl. Cibiru No. 400, Bandung', 0, 2),
('Taufik Hidayat', 'taufik.hidayat@example.com', '+628411234567', 'Jl. Ujung Berung No. 500, Bandung', 1, 2),
('Umi Kulsum', 'umi.kulsum@example.com', '+628421234567', 'Jl. Cicaheum No. 600, Bandung', 1, 2),
('Vino Bastian', 'vino.bastian@example.com', '+628431234567', 'Jl. Gedebage No. 700, Bandung', 0, 2),
('Wati Susanti', 'wati.susanti@example.com', '+628441234567', 'Jl. Kopo No. 800, Bandung', 1, 2),
('Xavier Tan', 'xavier.tan@example.com', '+628451234567', 'Jl. Setiabudi No. 900, Bandung', 1, 2),
('Yuni Shara', 'yuni.shara@example.com', '+628461234567', 'Jl. Lembang No. 1000, Bandung', 0, 2),
('Zacky Ahmad', 'zacky.ahmad@example.com', '+628471234567', 'Jl. Ciumbuleuit No. 1100, Bandung', 1, 2),
('Agus Salim', 'agus.salim@example.com', '+628481234567', 'Jl. Sarijadi No. 1200, Bandung', 1, 2),
('Bella Nurmalasari', 'bella.nurmalasari@example.com', '+628491234567', 'Jl. Cijerah No. 1300, Bandung', 0, 2),
('Candra Gunawan', 'candra.gunawan@example.com', '+628501234567', 'Jl. Cibaduyut No. 1400, Bandung', 1, 2),
('Dodi Prabowo', 'dodi.prabowo@example.com', '+628511234567', 'Jl. Cipaganti No. 1500, Bandung', 1, 2),
('Elsa Fitriani', 'elsa.fitriani@example.com', '+628521234567', 'Jl. Surapati No. 1600, Bandung', 0, 2),
('Fandi Ahmad', 'fandi.ahmad@example.com', '+628531234567', 'Jl. Lombok No. 1700, Bandung', 1, 2),
('Gina Melati', 'gina.melati@example.com', '+628541234567', 'Jl. Belitung No. 1800, Bandung', 1, 2),
('Haryanto', 'haryanto@example.com', '+628551234567', 'Jl. Flores No. 1900, Bandung', 0, 2),
('Indra Lesmana', 'indra.lesmana@example.com', '+628561234567', 'Jl. Halmahera No. 2000, Bandung', 1, 2),
('Johan Nawawi', 'johan.nawawi@example.com', '+628571234567', 'Jl. Seram No. 2100, Bandung', 1, 2),
('Kiki Amalia', 'kiki.amalia@example.com', '+628581234567', 'Jl. Timor No. 2200, Bandung', 0, 2),
('Lucky Hakim', 'lucky.hakim@example.com', '+628591234567', 'Jl. Bali No. 2300, Bandung', 1, 2),
('Maman Abdurahman', 'maman.abdurahman@example.com', '+628601234567', 'Jl. Jawa No. 2400, Bandung', 1, 2);