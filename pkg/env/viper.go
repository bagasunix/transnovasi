package env

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/phuslu/log"
	"github.com/spf13/viper"

	"github.com/bagasunix/transnovasi/pkg/errors"
)

type Cfg struct {
	App struct {
		Name        string `mapstructure:"name"`
		Version     string `mapstructure:"version"`
		Environment string `mapstructure:"environment"`
		TimeZone    string `mapstructure:"time_zone"`
	} `mapstructure:"app"`

	Server struct {
		Port        int    `mapstructure:"port"`
		Version     string `mapstructure:"version"`
		RateLimiter struct {
			Enabled  bool          `mapstructure:"enabled"`
			Limit    int           `mapstructure:"limit"`
			Duration time.Duration `mapstructure:"duration"`
		} `mapstructure:"rate_limiter"`
		Token struct {
			JWTKey string `mapstructure:"jwt_key"`
		} `mapstructure:"token"`
	} `mapstructure:"server"`

	Database struct {
		Driver        string        `mapstructure:"driver"`
		Host          string        `mapstructure:"host"`
		Port          int           `mapstructure:"port"`
		User          string        `mapstructure:"user"`
		Password      string        `mapstructure:"password"`
		DBName        string        `mapstructure:"dbname"`
		SSLMode       string        `mapstructure:"sslmode"`
		MaxConnection int           `mapstructure:"max_connection"`
		MaxIdleConns  int           `mapstructure:"max_idle"`
		MaxLifeTime   time.Duration `mapstructure:"max_life_time"`
		MaxIdleTime   time.Duration `mapstructure:"max_idle_time"`
	} `mapstructure:"database"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
		Type     string `mapstructure:"type"`
	} `mapstructure:"redis"`

	Logging struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"logging"`
}

func LoadCfg(ctx context.Context, path string) (*Cfg, error) {
	// Here you can implement loading configuration from file, environment variables, etc.
	// For simplicity, we'll just return nil for now.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	viper.AddConfigPath(absPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var config Cfg
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if config.Database.Driver == "" || config.Database.Host == "" || config.Database.Port == 0 || config.Database.User == "" || config.Database.Password == "" || config.Database.DBName == "" {
		return nil, errors.CustomError("database configuration is missing")
	}

	if config.Redis.Host == "" || config.Redis.Port == "" {
		return nil, errors.CustomError("redis configuration is missing")
	}

	return &config, nil
}

func (c *Cfg) GetRedisDSN() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// Fungsi untuk menginisialisasi konfigurasi
func InitConfig(ctx context.Context, logger *log.Logger) *Cfg {
	config, err := LoadCfg(ctx, "../")
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot load config")
		os.Exit(1)
	}

	logger.Info().Str("service", "config").Msg("Configuration loaded successfully")
	return config
}
