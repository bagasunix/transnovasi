package auditlog

import (
	"context"

	"github.com/bagasunix/transnovasi/internal/domains"
)

type Repository interface {
	Create(ctx context.Context, auditLog *domains.AuditLog) error
}
