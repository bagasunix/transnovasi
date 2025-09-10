package auditlog

import (
	"context"

	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/transnovasi/internal/domains"
	"github.com/bagasunix/transnovasi/pkg/errors"
)

type auditlogProvider struct {
	db     *gorm.DB
	logger *log.Logger
}

// Create implements Repository.
func (g *auditlogProvider) Create(ctx context.Context, model *domains.AuditLog) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Create(&model).Error)
}

// GetModelName returns the model name.
func (t *auditlogProvider) GetModelName() string {
	return "audit_log"
}

func NewGormAuditLog(logger *log.Logger, db *gorm.DB) Repository {
	t := new(auditlogProvider)
	t.db = db
	t.logger = logger
	return t
}
