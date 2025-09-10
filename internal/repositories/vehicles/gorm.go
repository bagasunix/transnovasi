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

func (g *gormProvider) Create(ctx context.Context, m []domains.Vehicle) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Create(&m).Error)
}
func (g *gormProvider) CreateBulkTx(ctx context.Context, tx any, m []domains.Vehicle) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), tx.(*gorm.DB).WithContext(ctx).Create(&m).Error)
}
func (g *gormProvider) Delete(ctx context.Context, id int) error {
	return errors.ErrSomethingWrong(g.logger, g.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]interface{}{"is_active": 0, "deleted_at": time.Now()}).Error)
}
func (g *gormProvider) Updates(ctx context.Context, id int, m *domains.Vehicle) error {
	return errors.ErrDuplicateValue(g.logger, g.GetModelName(), g.db.WithContext(ctx).Where("id = ?", id).Updates(m).Error)
}
func (g *gormProvider) GetOneById(ctx context.Context, id int) (result domains.SingleResult[*domains.Vehicle]) {
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), g.db.WithContext(ctx).First(&result.Value, "id = ?", id).Error)
	return result
}
func (g *gormProvider) GetOneByParams(ctx context.Context, params map[string]interface{}) (result domains.SingleResult[*domains.Vehicle]) {
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), g.db.WithContext(ctx).Where(params).First(&result.Value).Error)
	return result
}
func (g *gormProvider) GetAll(ctx context.Context, limit, offset int, search string) (result domains.SliceResult[domains.Vehicle]) {
	query := g.db.WithContext(ctx).Limit(limit).Offset(offset)
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), query.Find(&result.Value).Error)
	return result
}
func (g *gormProvider) GetAllByCustomerID(ctx context.Context, customerID, limit, offset int, search string) (result domains.SliceResult[domains.Vehicle]) {
	query := g.db.WithContext(ctx).Limit(limit).Offset(offset)
	if search != "" {
		query = query.Where("brand LIKE ? OR model LIKE ? OR plate_no LIKE ? OR year LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if customerID != 0 {
		query = query.Where("customer_id = ?", customerID)
	}
	result.Error = errors.ErrRecordNotFound(g.logger, g.GetModelName(), query.Find(&result.Value).Error)
	return result
}
func (g *gormProvider) GetConnection() (T any) {
	return g.db
}
func (g *gormProvider) GetModelName() string {
	return "vehicles"
}
func (g *gormProvider) CountVehicle(ctx context.Context, search string) (int, error) {
	var count int64
	query := g.db.WithContext(ctx).Model(&domains.Vehicle{})
	if search != "" {
		query = query.Where("brand LIKE ? OR model LIKE ? OR plate_no LIKE ? OR year LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, errors.ErrSomethingWrong(g.logger, err)
	}
	return int(count), nil
}
