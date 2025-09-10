package users

import (
	"context"

	"github.com/bagasunix/transnovasi/internal/domains"
	"github.com/bagasunix/transnovasi/internal/repositories/base"
)

type Common interface {
	Create(ctx context.Context, m *domains.User) error
	CreateTx(ctx context.Context, tx any, m *domains.User) error
	Delete(ctx context.Context, id int) error
	Updates(ctx context.Context, id int, m *domains.User) error
	GetAll(ctx context.Context, limit int, offset int, search string) (result domains.SliceResult[domains.User])

	GetOneByParams(ctx context.Context, params map[string]interface{}) (result domains.SingleResult[*domains.User])
}

type Repository interface {
	base.Repository
	Common
}
