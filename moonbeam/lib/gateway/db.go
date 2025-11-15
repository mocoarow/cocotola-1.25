package gateway

import (
	"errors"

	"github.com/go-sql-driver/mysql"
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
