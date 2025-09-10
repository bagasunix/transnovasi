package customers

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

func NewGormCustomer(logger *log.Logger, db *gorm.DB) Repository {
	g := new(gormProvider)
	g.db = db
	g.logger = logger
	return g
}

func (g *gormProvider) Create(ctx context.Context, m *domains.Customer) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Create(m).Error)
}
func (g *gormProvider) CreateTx(ctx context.Context, tx any, m *domains.Customer) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), tx.(*gorm.DB).WithContext(ctx).Create(m).Error)
}
func (g *gormProvider) Delete(ctx context.Context, id int) error {
	return errors.ErrSomethingWrong(g.logger, g.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]interface{}{"is_active": 0, "deleted_at": time.Now()}).Error)
}
func (g *gormProvider) Updates(ctx context.Context, id int, m *domains.Customer) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Where("id = ?", id).Updates(m).Error)
}
func (g *gormProvider) GetOneByParams(ctx context.Context, params map[string]interface{}) (result domains.SingleResult[domains.Customer]) {
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), g.db.WithContext(ctx).Preload("Vehicles").Where(params).First(&result.Value).Error)
	return result
}
func (g *gormProvider) GetAll(ctx context.Context, limit int, offset int, search string) (result domains.SliceResult[domains.Customer]) {
	query := g.db.WithContext(ctx).Limit(limit).Offset(offset)
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), query.Find(&result.Value).Error)
	return result
}
func (g *gormProvider) GetConnection() (T any) {
	return g.db
}
func (g *gormProvider) GetModelName() string {
	return "customers"
}
func (g *gormProvider) CountCustomer(ctx context.Context, search string) (int, error) {
	var count int64
	query := g.db.WithContext(ctx).Model(&domains.Customer{})
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, errors.ErrSomethingWrong(g.logger, err)
	}
	return int(count), nil
}
