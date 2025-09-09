package vehicles

import (
	"context"
	"time"

	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/transnovasi/internal/domains"
	"github.com/bagasunix/transnovasi/pkg/errors"
)

type gormProvider struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewGormVehicle(logger *log.Logger, db *gorm.DB) Repository {
	g := new(gormProvider)
	g.db = db
	g.logger = logger
	return g
}

func (g *gormProvider) Create(ctx context.Context, m *domains.Vehicle) error {
	return errors.ErrDuplicateValue(g.logger, "vehicles", g.db.WithContext(ctx).Create(m).Error)
}
func (g *gormProvider) Delete(ctx context.Context, id int) error {
	return errors.ErrSomethingWrong(g.logger, g.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]interface{}{"is_active": 0, "deleted_at": time.Now()}).Error)
}
func (g *gormProvider) Updates(ctx context.Context, id int, m *domains.Vehicle) error {
	return errors.ErrDuplicateValue(g.logger, "vehicles", g.db.WithContext(ctx).Where("id = ?", id).Updates(m).Error)
}
func (g *gormProvider) GetOneById(ctx context.Context, id int) (result domains.SingleResult[*domains.Vehicle]) {
	result.Error = errors.ErrRecordNotFound(g.logger, "vehicles", g.db.WithContext(ctx).First(&result.Value, "id = ?", id).Error)
	return result
}
func (g *gormProvider) GetOneByParams(ctx context.Context, params map[string]interface{}) (result domains.SingleResult[*domains.Vehicle]) {
	result.Error = errors.ErrRecordNotFound(g.logger, "vehicles", g.db.WithContext(ctx).Where(params).First(&result.Value).Error)
	return result
}
func (g *gormProvider) GetAll(ctx context.Context, limit int, offset int, search string) (result domains.SliceResult[domains.Vehicle]) {
	result.Error = errors.ErrRecordNotFound(g.logger, "vehicles", g.db.WithContext(ctx).Limit(int(limit)).Find(&result.Value).Error)
	return result
}
