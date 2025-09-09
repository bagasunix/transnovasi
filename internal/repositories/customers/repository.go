package customers

import (
	"context"

	"github.com/bagasunix/transnovasi/internal/domains"
)

type Repository interface {
	Create(ctx context.Context, m *domains.Customer) error
	Delete(ctx context.Context, id int) error
	Updates(ctx context.Context, id int, m *domains.Customer) error
	GetAll(ctx context.Context, limit int, offset int, search string) (result domains.SliceResult[domains.Customer])

	GetOneByParams(ctx context.Context, params map[string]interface{}) (result domains.SingleResult[*domains.Customer])
}
