package config

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
)

type DBConfig struct {
	DriverName string         `yaml:"driverName"`
	MySQL      *MySQLConfig   `yaml:"mysql"`
	SQLite3    *SQLite3Config `yaml:"sqlite3"`
}
type DBConnection struct {
	DriverName string
	Dialect    libgateway.DialectRDBMS
	DB         *gorm.DB
}

func InitDB(ctx context.Context, dbConfig *DBConfig, logConfig *LogConfig, appName string) (*DBConnection, func(), error) {
	initDBFunc, ok := initDBs[dbConfig.DriverName]
	if !ok {
		return nil, nil, fmt.Errorf("invalid database driver: %s", dbConfig.DriverName)
	}
	dbLogLevel := slog.LevelWarn
	if level, ok := logConfig.Levels["db"]; ok {
		dbLogLevel = stringToLogLevel(level)
	}

	dialect, db, sqlDB, err := initDBFunc(ctx, dbConfig, dbLogLevel, appName)
	if err != nil {
		return nil, nil, fmt.Errorf("init DB: %w", err)
	}

	dbConn := DBConnection{
		DriverName: dbConfig.DriverName,
		Dialect:    dialect,
		DB:         db,
	}

	return &dbConn, func() {
		sqlDB.Close()
	}, nil
}
