package gateway

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
	"gorm.io/gorm"
)

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
