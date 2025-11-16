package gateway

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"

	migrate_sqlite3 "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway/sqlite3"
)

type DialectSQLite3 struct {
}

func (d *DialectSQLite3) Name() string {
	return "sqlite3"
}

func (d *DialectSQLite3) BoolDefaultValue() string {
	return "0"
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
