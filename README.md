# TRANSNOVASI API


## Technologies Used

- [Go](https://golang.org/)
- [Fiber](github.com/gofiber/fiber/v2)
- [Redis](https://redis.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Logger Phuslu](github.com/phuslu/log)
- [Swagger](github.com/gofiber/swagger)

---

## Prerequisites

- Docker
- Docker Compose
- Go (if running without Docker)
- Redis
- Postgresql

---

## ✨ Features

- User Authentication (JWT)
- Customer CRUD (with soft delete)
- Rate Limiting Windows Sliding 
- Docker & Docker Compose support
- Swagger documentation
- Logging via Phuslu

---


## Configuration

The configuration file `config.yaml` should be placed in the `config` directory. This file will be mounted into the Docker container.

---

## Database Setup

Before running the application, make sure you already have a PostgreSQL database created.  
You **don’t need to create tables/columns manually** since the project uses auto migration.  

### Option 1: Using psql on your host
```sql
CREATE DATABASE transnovasi;
```
### Option 2: Using docker exec (if running with Docker)
If you’re running PostgreSQL from docker-compose, you can create the database like this:
```bash
# masuk ke container database
docker exec -it transnovasi_db bash

# login ke PostgreSQL (ganti user/password sesuai docker-compose)
psql -U postgres

# buat database
CREATE DATABASE transnovasi;

# keluar dari psql
\q
```
Make sure the database name, user, and password match with the configuration in "config.yaml".

---

## 🚀 Getting Started

### 1. Clone & Setup
```bash
git clone https://github.com/bagasunix/transnovasi.git
cd transnovasi
```
### 2. Run with Docker
```bash
docker-compose up --build
```
### 3. Run the application along with its dependencies:
```sh
docker-compose up
```
- App runs on: http://localhost:8080
- Swagger docs: http://localhost:8080/swagger/index.html

### 4. Tech Stack
```bash
Golang
Fiber
GORM
Postgresql
Redis
phuslu/log
Swagger UI
```

### 5. Add Sample Postman Curl
#### Auth API
- Logout User
```bash
curl --location 'http://localhost:8080/api/v1/auth/register' \
--header 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMSIsIm5hbWUiOiJBbGRpbm8gUHJhdGFtYSBCYWdhc2thcmEiLCJzZXgiOiIxIiwiZW1haWwiOiJhbGRpbm9wcmF0YW1hMTVAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiaXNfYWN0aXZlIjoiMSIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE3NTc2MjI5NDV9.37oDXiK--6-anBWYT-Zq5E2jyntpMNY7GEiQ4wW3hExdXw9Y17WWDW0PRSdS3icFTOy9OmxZMsSz8vfocZZBdg'
```
- Register User
```bash
curl --location 'http://localhost:8080/api/v1/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "Aldino Pratama",
  "sex": 1,
  "email": "emailku1@gmail.com",
  "password": "password123",
  "role": ADMIN
}'
```
- Login Admin
```bash
- curl --location 'http://localhost:8080/api/v1/auth' \
--header 'Content-Type: application/json' \
--data-raw '{"email":"emailku1@gmail.com", "password":"password123"}'
```
- Login Operator
```bash
- curl --location 'http://localhost:8080/api/v1/auth' \
--header 'Content-Type: application/json' \
--data-raw '{"email":"emailku2@gmail.com", "password":"password123"}'
```
- Login User
```bash
- curl --location 'http://localhost:8080/api/v1/auth' \
--header 'Content-Type: application/json' \
--data-raw '{"email":"emailku3@gmail.com", "password":"password123"}'
```
#### Customer API
- Create Customer With Vehicle
```bash
curl --location 'http://localhost:8080/api/v1/customer' \
--header 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMSIsIm5hbWUiOiJBbGRpbm8gUHJhdGFtYSBCYWdhc2thcmEiLCJzZXgiOiIxIiwiZW1haWwiOiJhbGRpbm9wcmF0YW1hMTVAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiaXNfYWN0aXZlIjoiMSIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE3NTc2MjI5NDV9.37oDXiK--6-anBWYT-Zq5E2jyntpMNY7GEiQ4wW3hExdXw9Y17WWDW0PRSdS3icFTOy9OmxZMsSz8vfocZZBdg' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "Perusahaan Transportasi Besar Kecil",
  "email": "bigtransport@example.com",
  "phone": "081234567890",
  "address": "Jl. Raya Pusat Kota No. 1000",
  "vehicle": [
    {
      "plate_no": "B 1000 ABC",
      "model": "Truck",
      "brand": "Mitsubishi",
      "color": "Biru",
      "year": "2020",
      "max_speed": 120,
      "fuel_type": "Solar"
    },
    {
      "plate_no": "B 1001 DEF",
      "model": "Van",
      "brand": "Toyota",
      "color": "Putih",
      "year": "2021",
      "max_speed": 150,
      "fuel_type": "Bensin"
    },
    {
      "plate_no": "B 1002 GHI",
      "model": "Pickup",
      "brand": "Isuzu",
      "color": "Hitam",
      "year": "2019",
      "max_speed": 130,
      "fuel_type": "Solar"
    },
    {
      "plate_no": "B 1999 XYZ",
      "model": "Sedan",
      "brand": "Honda",
      "color": "Silver",
      "year": "2022",
      "max_speed": 180,
      "fuel_type": "Bensin"
    }
  ]
}'
```
- Create Customer Without Vehicle
```bash
curl --location 'http://localhost:8080/api/v1/customer' \
--header 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMSIsIm5hbWUiOiJBbGRpbm8gUHJhdGFtYSBCYWdhc2thcmEiLCJzZXgiOiIxIiwiZW1haWwiOiJhbGRpbm9wcmF0YW1hMTVAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiaXNfYWN0aXZlIjoiMSIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE3NTc2MjI5NDV9.37oDXiK--6-anBWYT-Zq5E2jyntpMNY7GEiQ4wW3hExdXw9Y17WWDW0PRSdS3icFTOy9OmxZMsSz8vfocZZBdg' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "Perusahaan Transportasi Besar Kecil",
  "email": "bigtransport@example.com",
  "phone": "081234567890",
  "address": "Jl. Raya Pusat Kota No. 1000"
}'
```
- Get All Customer
```bash
curl --location 'http://localhost:8080/api/v1/customer?page=1&limit=&search' \
--header 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMSIsIm5hbWUiOiJBbGRpbm8gUHJhdGFtYSBCYWdhc2thcmEiLCJzZXgiOiIxIiwiZW1haWwiOiJhbGRpbm9wcmF0YW1hMTVAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiaXNfYWN0aXZlIjoiMSIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE3NTc2MjI5NDV9.37oDXiK--6-anBWYT-Zq5E2jyntpMNY7GEiQ4wW3hExdXw9Y17WWDW0PRSdS3icFTOy9OmxZMsSz8vfocZZBdg'
```
- Get Customer By ID
```bash
curl --location 'http://localhost:8080/api/v1/customer/1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMSIsIm5hbWUiOiJBbGRpbm8gUHJhdGFtYSBCYWdhc2thcmEiLCJzZXgiOiIxIiwiZW1haWwiOiJhbGRpbm9wcmF0YW1hMTVAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiaXNfYWN0aXZlIjoiMSIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE3NTc2MjI5NDV9.37oDXiK--6-anBWYT-Zq5E2jyntpMNY7GEiQ4wW3hExdXw9Y17WWDW0PRSdS3icFTOy9OmxZMsSz8vfocZZBdg'
```
- Get Customer By ID With Vehicle
```bash
curl --location 'http://localhost:8080/api/v1/customer/1/vehicle?page=1&limit=&search' \
--header 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMSIsIm5hbWUiOiJBbGRpbm8gUHJhdGFtYSBCYWdhc2thcmEiLCJzZXgiOiIxIiwiZW1haWwiOiJhbGRpbm9wcmF0YW1hMTVAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiaXNfYWN0aXZlIjoiMSIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE3NTc2MjI5NDV9.37oDXiK--6-anBWYT-Zq5E2jyntpMNY7GEiQ4wW3hExdXw9Y17WWDW0PRSdS3icFTOy9OmxZMsSz8vfocZZBdg'
```
- Delete Customer By ID
```bash
curl --location --request DELETE 'http://localhost:8080/api/v1/customer/1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMSIsIm5hbWUiOiJBbGRpbm8gUHJhdGFtYSBCYWdhc2thcmEiLCJzZXgiOiIxIiwiZW1haWwiOiJhbGRpbm9wcmF0YW1hMTVAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiaXNfYWN0aXZlIjoiMSIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE3NTc2MjI5NDV9.37oDXiK--6-anBWYT-Zq5E2jyntpMNY7GEiQ4wW3hExdXw9Y17WWDW0PRSdS3icFTOy9OmxZMsSz8vfocZZBdg'
```
- Update Customer By ID
```bash
curl --location --request PUT 'http://localhost:8080/api/v1/customer/1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMSIsIm5hbWUiOiJBbGRpbm8gUHJhdGFtYSBCYWdhc2thcmEiLCJzZXgiOiIxIiwiZW1haWwiOiJhbGRpbm9wcmF0YW1hMTVAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiaXNfYWN0aXZlIjoiMSIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE3NTc2MjI5NDV9.37oDXiK--6-anBWYT-Zq5E2jyntpMNY7GEiQ4wW3hExdXw9Y17WWDW0PRSdS3icFTOy9OmxZMsSz8vfocZZBdg'\
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "Perusahaan Transportasi Besar",
  "phone": "081234567890",
  "address": "Jl. Raya Pusat Kota No. 1000"
}'
```
#### Vehicle API
- Get All Vehicle
```bash
curl --location 'http://localhost:8080/api/v1/vehicle?page=1&limit=&search&customer_id=51' \
--header 'Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMSIsIm5hbWUiOiJBbGRpbm8gUHJhdGFtYSBCYWdhc2thcmEiLCJzZXgiOiIxIiwiZW1haWwiOiJhbGRpbm9wcmF0YW1hMTVAZ21haWwuY29tIiwicm9sZSI6IkFETUlOIiwiaXNfYWN0aXZlIjoiMSIsImNyZWF0ZWRfYXQiOiIwMDAxLTAxLTAxVDAwOjAwOjAwWiJ9LCJleHAiOjE3NTc2MjI5NDV9.37oDXiK--6-anBWYT-Zq5E2jyntpMNY7GEiQ4wW3hExdXw9Y17WWDW0PRSdS3icFTOy9OmxZMsSz8vfocZZBdg'
```
## 🔐 Endpoint Access
### Auth Management
```bash
POST /api/auth/logout          → ADMIN, OPERATOR, USER
POST /api/auth                → ADMIN, OPERATOR, USER
```
### Customer Management (CRUD)
```bash
GET    /api/v1/customers                      → ADMIN, OPERATOR, USER
GET    /api/v1/customers/:id/vehicle          → ADMIN, OPERATOR, USER
POST   /api/v1/customers                      → ADMIN, OPERATOR
GET    /api/v1/customers/:id                  → ADMIN, OPERATOR, USER
PUT    /api/v1/customers/:id                  → ADMIN, OPERATOR  
DELETE /api/v1/customers/:id                  → ADMIN only
```
### Vehicle Management
```bash
GET /api/v1/vehicle          → ADMIN, OPERATOR, USER
```