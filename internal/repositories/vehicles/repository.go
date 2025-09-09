package vehicles

import (
	"context"

	"github.com/bagasunix/transnovasi/internal/domains"
)

type Repository interface {
	Create(ctx context.Context, m *domains.Vehicle) error
	Delete(ctx context.Context, id int) error
	Updates(ctx context.Context, id int, m *domains.Vehicle) error
	GetAll(ctx context.Context, limit int, offset int, search string) (result domains.SliceResult[domains.Vehicle])

	GetOneById(ctx context.Context, id int) (result domains.SingleResult[*domains.Vehicle])
	GetOneByParams(ctx context.Context, params map[string]interface{}) (result domains.SingleResult[*domains.Vehicle])
}
