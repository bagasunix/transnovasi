package users

import (
	"context"

	"github.com/bagasunix/transnovasi/internal/domains"
)

type Repository interface {
	Create(ctx context.Context, m *domains.User) error
	Delete(ctx context.Context, id int) error
	Updates(ctx context.Context, id int, m *domains.User) error
	GetAll(ctx context.Context, limit int, offset int, search string) (result domains.SliceResult[domains.User])

	GetOneByParams(ctx context.Context, params map[string]interface{}) (result domains.SingleResult[*domains.User])
}
