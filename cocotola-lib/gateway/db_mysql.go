package gateway

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-sql-driver/mysql"
	slog_gorm "github.com/orandin/slog-gorm"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
)

type DialectMySQL struct {
}

func (d *DialectMySQL) Name() string {
	return "mysql"
}

func (d *DialectMySQL) BoolDefaultValue() string {
	return "0"
}

type MySQLConfig struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	Database string `yaml:"database" validate:"required"`
}

func initDBMySQL(ctx context.Context, cfg *DBConfig, logLevel slog.Level, appName string) (DialectRDBMS, *gorm.DB, *sql.DB, error) {
	db, err := OpenMySQL(cfg.MySQL, logLevel, appName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("OpenMySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf("ping: %w", err)
	}

	dialect := DialectMySQL{}
	return &dialect, db, sqlDB, nil
}

func OpenMySQLWithDSN(dsn string, logLevel slog.Level, appName string) (*gorm.DB, error) {
	gormDialector := gorm_mysql.Open(dsn)

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
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	return db, nil
}

func OpenMySQL(cfg *MySQLConfig, logLevel slog.Level, appName string) (*gorm.DB, error) {
	c := mysql.Config{ //nolint:exhaustruct
		DBName:               cfg.Database,
		User:                 cfg.Username,
		Passwd:               cfg.Password,
		Addr:                 fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Net:                  "tcp",
		ParseTime:            true,
		MultiStatements:      true,
		Params:               map[string]string{"charset": "utf8mb4"},
		Collation:            "utf8mb4_bin",
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		MaxAllowedPacket:     64 << 20, // 64 MiB.
		Loc:                  time.UTC,
	}

	return OpenMySQLWithDSN(c.FormatDSN(), logLevel, appName)
}
