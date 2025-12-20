package gateway

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
	"gorm.io/gorm"
)

type DBConfig struct {
	DriverName string         `yaml:"driverName"`
	MySQL      *MySQLConfig   `yaml:"mysql"`
	SQLite3    *SQLite3Config `yaml:"sqlite3"`
}
type DBConnection struct {
	DriverName string
	Dialect    DialectRDBMS
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

type DialectRDBMS interface {
	Name() string
	BoolDefaultValue() string
}

const MySQLErDupEntry = 1062
const MySQLErNoReferencedRow2 = 1452

const SQLiteConstraintPrimaryKey = 1555
const SQLiteConstraintUnique = 2067

type sqliteError interface {
	error
	Code() int
}

func ConvertDuplicatedError(err error, newErr error) error {
	var mysqlErr *mysql.MySQLError
	if ok := errors.As(err, &mysqlErr); ok {
		switch mysqlErr.Number {
		case MySQLErDupEntry, MySQLErNoReferencedRow2:
			return newErr
		}
	}

	var sqlite3Err sqliteError
	if ok := errors.As(err, &sqlite3Err); ok {
		switch sqlite3Err.Code() {
		case SQLiteConstraintPrimaryKey, SQLiteConstraintUnique:
			return newErr
		}
	}

	return err
}

func MigrateDB(db *gorm.DB, driverName string, sourceDriver source.Driver, getDatabaseDriver func(sqlDB *sql.DB) (database.Driver, error)) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("db.DB in gateway.migrateDB. err: %w", err)
	}

	databaseDriver, err := getDatabaseDriver(sqlDB)
	if err != nil {
		return fmt.Errorf("getDatabaseDriver in gateway.migrateDB. err: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, driverName, databaseDriver)
	if err != nil {
		return fmt.Errorf("NewWithInstance in gateway.migrateDB. err: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to m.Up in gateway.migrateDB. err: %w", err)
	}

	return nil
}
