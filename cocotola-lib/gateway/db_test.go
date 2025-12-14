package gateway_test

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"

	gateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
)

func TestConvertDuplicatedError_shouldReturnNewError_whenKnownConstraint(t *testing.T) {
	t.Parallel()

	newErr := errors.New("converted")

	tests := []struct {
		name  string
		input error
	}{
		{name: "mysql duplicate", input: &mysql.MySQLError{Number: gateway.MySQLErDupEntry}},
		{name: "mysql fk", input: &mysql.MySQLError{Number: gateway.MySQLErNoReferencedRow2}},
		{name: "sqlite pk", input: sqliteTestError{code: gateway.SQLiteConstraintPrimaryKey}},
		{name: "sqlite unique", input: sqliteTestError{code: gateway.SQLiteConstraintUnique}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := gateway.ConvertDuplicatedError(tt.input, newErr); !errors.Is(got, newErr) {
				t.Errorf("got %v, want %v", got, newErr)
			}
		})
	}
}

func TestConvertDuplicatedError_shouldReturnOriginal_whenUnknownError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("sentinel")
	if got := gateway.ConvertDuplicatedError(sentinel, errors.New("converted")); !errors.Is(got, sentinel) {
		t.Errorf("expected original error, got %v", got)
	}
}

type sqliteTestError struct {
	code int
}

func (e sqliteTestError) Error() string {
	return "sqlite"
}

func (e sqliteTestError) Code() int {
	return e.code
}
