package gateway

type DialectSQLite3 struct {
}

func (d *DialectSQLite3) Name() string {
	return "sqlite3"
}

func (d *DialectSQLite3) BoolDefaultValue() string {
	return "0"
}
