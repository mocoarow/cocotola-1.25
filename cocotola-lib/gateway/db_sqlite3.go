package gateway

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"

	gorm_sqlite "github.com/glebarez/sqlite"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	slog_gorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	migrate_sqlite3 "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway/sqlite3"
)

type SQLite3Config struct {
	File string `yaml:"file" validate:"required"`
}
type DialectSQLite3 struct {
}

func (d *DialectSQLite3) Name() string {
	return "sqlite3"
}

func (d *DialectSQLite3) BoolDefaultValue() string {
	return "0"
}

func initDBSQLite3(ctx context.Context, cfg *DBConfig, logLevel slog.Level, appName string) (DialectRDBMS, *gorm.DB, *sql.DB, error) {
	db, err := OpenSQLite3(cfg.SQLite3, logLevel, appName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("OpenSQLite file(%s): %w", cfg.SQLite3.File, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("DB. file(%s): %w", cfg.SQLite3.File, err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf("ping. file(%s): %w", cfg.SQLite3.File, err)
	}

	dialect := DialectSQLite3{}
	return &dialect, db, sqlDB, nil
}

func OpenSQLite3(cfg *SQLite3Config, logLevel slog.Level, appName string) (*gorm.DB, error) {
	gormDialector := gorm_sqlite.Open(cfg.File)

	options := make([]slog_gorm.Option, 0)
	options = append(options, slog_gorm.WithHandler(slog.Default().With(slog.String(domain.LoggerNameKey, appName+"-gorm")).Handler()))
	if logLevel == slog.LevelDebug {
		options = append(options, slog_gorm.WithTraceAll()) // trace all messages
	}

	gormConfig := gorm.Config{ //nolint:exhaustruct
		Logger: slog_gorm.New(options...),
	}

	db, err := gorm.Open(gormDialector, &gormConfig)
	if err != nil {
		return nil, fmt.Errorf("open sqlite3: %w", err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	return db, nil
}

func MigrateSQLite3DB(db *gorm.DB, sqlFS fs.FS) error {
	driverName := "sqlite3"
	sourceDriver, err := iofs.New(sqlFS, driverName)
	if err != nil {
		return fmt.Errorf("iofs.New: %w", err)
	}

	return MigrateDB(db, driverName, sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
		return migrate_sqlite3.WithInstance(sqlDB, &migrate_sqlite3.Config{}) //nolint:exhaustruct
	})
}
