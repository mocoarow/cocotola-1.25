package config

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type DBConfig struct {
	DriverName string         `yaml:"driverName"`
	MySQL      *MySQLConfig   `yaml:"mysql"`
	SQLite3    *SQLite3Config `yaml:"sqlite3"`
}

func InitDB(ctx context.Context, dbConfig *DBConfig, logConfig *LogConfig, appName string) (*gorm.DB, func(), error) {
	initDBFunc, ok := initDBs[dbConfig.DriverName]
	if !ok {
		return nil, nil, fmt.Errorf("invalid database driver: %s", dbConfig.DriverName)
	}
	dbLogLevel := slog.LevelWarn
	if level, ok := logConfig.Levels["db"]; ok {
		dbLogLevel = stringToLogLevel(level)
	}

	db, sqlDB, err := initDBFunc(ctx, dbConfig, dbLogLevel, appName)
	if err != nil {
		return nil, nil, fmt.Errorf("init DB: %w", err)
	}

	return db, func() {
		sqlDB.Close()
	}, nil
}
