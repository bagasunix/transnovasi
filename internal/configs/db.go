package configs

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migPostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/phuslu/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/bagasunix/transnovasi/pkg/configs"
	"github.com/bagasunix/transnovasi/pkg/env"
	"github.com/bagasunix/transnovasi/pkg/errors"
)

func InitDB(ctx context.Context, cfg *env.Cfg, logger *log.Logger) *gorm.DB {
	CfgBuild := &configs.DBConfig{
		Driver:          cfg.Database.Driver,
		Host:            cfg.Database.Host,
		Port:            strconv.Itoa(cfg.Database.Port),
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DatabaseName:    cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxConnection,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.MaxLifeTime,
		ConnMaxIdleTime: cfg.Database.MaxIdleTime,
		Timezone:        cfg.App.TimeZone,
	}
	return NewPostgresDB(ctx, CfgBuild, logger)
}

func NewPostgresDB(ctx context.Context, cfg *configs.DBConfig, logger *log.Logger) *gorm.DB {
	// Membuat koneksi ke database dengan DSN dari dbConfig
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()+cfg.DatabaseName+"?sslmode="+cfg.SSLMode), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	errors.HandlerWithOSExit(logger, err, "init", "database", "config", cfg.GetDSN())
	// Mengatur parameter koneksi
	errors.HandlerWithOSExit(logger, db.WithContext(ctx).Use(dbresolver.Register(dbresolver.Config{}).
		SetMaxOpenConns(cfg.MaxOpenConns).
		SetMaxIdleConns(cfg.MaxIdleConns).
		SetConnMaxLifetime(cfg.ConnMaxLifetime*time.Minute).
		SetConnMaxIdleTime(cfg.ConnMaxIdleTime*time.Minute)),
		"db_resolver")

	sqlDB, _ := db.DB()
	driver, err := migPostgres.WithInstance(sqlDB, &migPostgres.Config{})
	errors.HandlerWithOSExit(logger, err, "failed to initialize postgres driver")

	migrationsPath := "./migrations"
	// Buat instance migrasi
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	errors.HandlerWithOSExit(logger, err, "failed to create migration instance")

	// Jalankan migrasi ke versi terbaru
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		errors.HandlerWithOSExit(logger, err, "Failed to run migrations.")
	}

	fmt.Println("Migration successful")

	// Verifikasi koneksi
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		errors.HandlerWithOSExit(logger, err, "init", "database", "ping", "")
	}
	return db
}
