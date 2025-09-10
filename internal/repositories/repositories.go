package repositories

import (
	"github.com/phuslu/log"
	"gorm.io/gorm"

	auditlog "github.com/bagasunix/transnovasi/internal/repositories/audit_log"
	"github.com/bagasunix/transnovasi/internal/repositories/customers"
	"github.com/bagasunix/transnovasi/internal/repositories/users"
	"github.com/bagasunix/transnovasi/internal/repositories/vehicles"
)

type Repositories interface {
	GetCustomer() customers.Repository
	GetUser() users.Repository
	GetVehicle() vehicles.Repository
	GetAuditLog() auditlog.Repository
}

type repo struct {
	customer customers.Repository
	user     users.Repository
	vehicle  vehicles.Repository
	auditlog auditlog.Repository
}

func (r *repo) GetCustomer() customers.Repository {
	return r.customer
}
func (r *repo) GetUser() users.Repository {
	return r.user
}
func (r *repo) GetVehicle() vehicles.Repository {
	return r.vehicle
}
func (r *repo) GetAuditLog() auditlog.Repository {
	return r.auditlog
}

func New(logger *log.Logger, db *gorm.DB) Repositories {
	rs := new(repo)
	rs.user = users.NewGormUser(logger, db)
	rs.customer = customers.NewGormCustomer(logger, db)
	rs.vehicle = vehicles.NewGormVehicle(logger, db)
	rs.auditlog = auditlog.NewGormAuditLog(logger, db)
	return rs
}
