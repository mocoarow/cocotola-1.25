package gateway

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	gorm_sqlite "github.com/glebarez/sqlite"
	slog_gorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"

	"github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway"
)

const testSQLite3File = "./test_sqlite3.db"

type SQLite3Config struct {
	File string `yaml:"file" validate:"required"`
}

func openSQLiteForTest() (*gorm.DB, error) {
	return OpenSQLite3(&SQLite3Config{
		File: testSQLite3File,
	}, slog.LevelInfo, "test")
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

	return db, nil
}

// func OpenSQLiteInMemory(sqlFS embed.FS) (*gorm.DB, error) {
// 	logger := slog.Default()
// 	db, err := gorm.Open(gormSQLite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{
// 		Logger: slog_gorm.New(
// 			slog_gorm.WithLogger(logger), // Optional, use slog.Default() by default
// 			slog_gorm.WithTraceAll(),     // trace all messages
// 		),
// 	})
// 	if err != nil {
// 		return nil, liberrors.Errorf("gorm.Open. err: %w", err)
// 	}
// 	if err := setupSQLite(sqlFS, db); err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }

// func setupSQLite(sqlFS embed.FS, db *gorm.DB) error {
// 	driverName := "sqlite3"
// 	sourceDriver, err := iofs.New(sqlFS, driverName)
// 	if err != nil {
// 		return err
// 	}
// 	return setupDB(db, driverName, sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
// 		return migrate_sqlite3.WithInstance(sqlDB, &migrate_sqlite3.Config{})
// 	})
// }

func InitSQLit3InFile(fs embed.FS) (*gorm.DB, error) {
	os.Remove(testSQLite3File)
	db, err := openSQLiteForTest()
	if err != nil {
		return nil, fmt.Errorf("openSQLiteForTest: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("sqlDB.ping: %w", err)
	}

	if err := libgateway.MigrateSQLite3DB(db, fs); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}
