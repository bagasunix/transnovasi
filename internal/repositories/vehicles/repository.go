package vehicles

import (
	"context"

	"github.com/bagasunix/transnovasi/internal/domains"
	"github.com/bagasunix/transnovasi/internal/repositories/base"
)

type Common interface {
	Create(ctx context.Context, m []domains.Vehicle) error
	CreateBulkTx(ctx context.Context, tx any, m []domains.Vehicle) error
	Delete(ctx context.Context, id int) error
	Updates(ctx context.Context, id int, m *domains.Vehicle) error
	GetAll(ctx context.Context, limit, offset int, search string) (result domains.SliceResult[domains.Vehicle])
	GetAllByCustomerID(ctx context.Context, customerID, limit, offset int, search string) (result domains.SliceResult[domains.Vehicle])
	CountVehicle(ctx context.Context, search string) (int, error)

	GetOneById(ctx context.Context, id int) (result domains.SingleResult[*domains.Vehicle])
	GetOneByParams(ctx context.Context, params map[string]interface{}) (result domains.SingleResult[*domains.Vehicle])
}

type Repository interface {
	base.Repository
	Common
}
