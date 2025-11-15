package gateway

type DialectMySQL struct {
}

func (d *DialectMySQL) Name() string {
	return "mysql"
}

func (d *DialectMySQL) BoolDefaultValue() string {
	return "0"
}
