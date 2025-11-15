package config

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	gorm_sqlite "github.com/glebarez/sqlite"
	slog_gorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"

	"github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
)

type SQLite3Config struct {
	File string `yaml:"file" validate:"required"`
}

func initDBSQLite3(ctx context.Context, cfg *DBConfig, logLevel slog.Level, appName string) (*gorm.DB, *sql.DB, error) {
	db, err := OpenSQLite3(cfg.SQLite3, logLevel, appName)
	if err != nil {
		return nil, nil, fmt.Errorf("OpenSQLite file(%s): %w", cfg.SQLite3.File, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("DB. file(%s): %w", cfg.SQLite3.File, err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("ping. file(%s): %w", cfg.SQLite3.File, err)
	}

	return db, sqlDB, nil
}

func OpenSQLite3(cfg *SQLite3Config, logLevel slog.Level, appName string) (*gorm.DB, error) {
	gormDialector := gorm_sqlite.Open(cfg.File)

	options := make([]slog_gorm.Option, 0)
	options = append(options, slog_gorm.WithHandler(slog.Default().With(slog.String(domain.LoggerNameKey, appName+"-gorm")).Handler()))
	if logLevel == slog.LevelDebug {
		options = append(options, slog_gorm.WithTraceAll()) // trace all messages
	}

	gormConfig := gorm.Config{ //nolint:exhaustruct
		Logger: slog_gorm.New(options...), // trace all messages
		// slog_gorm.WithContextFunc(liblog.LoggerNameKey, func(_ context.Context) (slog.Value, bool) {
		// 	return slog.StringValue(appName + "-gorm"), true
		// }),
		// slog_gorm.SetLogLevel(slog_gorm.DefaultLogType, slog.LevelDebug),

	}

	return gorm.Open(gormDialector, &gormConfig) //nolint:wrapcheck
}
