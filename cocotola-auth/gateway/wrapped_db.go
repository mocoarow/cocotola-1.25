package gateway

import (
	"fmt"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
)

type HasTableName interface {
	TableName() string
}

type wrappedDB struct {
	dbc            *libgateway.DBConnection
	organizationID *domain.OrganizationID
}

func (x *wrappedDB) Table(name string, args ...interface{}) *wrappedDB {
	x.dbc.DB = x.dbc.DB.Table(name, args...)

	return x
}

func (x *wrappedDB) Select(query interface{}, args ...interface{}) *wrappedDB {
	x.dbc.DB = x.dbc.DB.Select(query, args...)

	return x
}

func (x *wrappedDB) Where(query interface{}, args ...interface{}) *wrappedDB {
	x.dbc.DB = x.dbc.DB.Where(query, args...)

	return x
}

func (x *wrappedDB) Joins(query string, args ...interface{}) *wrappedDB {
	x.dbc.DB = x.dbc.DB.Joins(query, args...)

	return x
}

func (x *wrappedDB) WhereOrganizationID(table HasTableName, organizationID *domain.OrganizationID) *wrappedDB {
	x.dbc.DB = x.dbc.DB.Where(fmt.Sprintf("%s.organization_id = ?", table.TableName()), organizationID.Int())

	return x
}

func (x *wrappedDB) WhereNotDeleted(table HasTableName) *wrappedDB {
	x.dbc.DB = x.dbc.DB.Where(fmt.Sprintf("%s.deleted = ?", table.TableName()), x.dbc.Dialect.BoolDefaultValue())

	return x
}

func (x *wrappedDB) WhereUser() *wrappedDB {
	return x.WhereOrganizationID(&userEntity{}, x.organizationID).WhereNotDeleted(&userEntity{}) //nolint:exhaustruct
}

func (x *wrappedDB) WhereUserGroup() *wrappedDB {
	return x.WhereOrganizationID(&userGroupEntity{}, x.organizationID).WhereNotDeleted(&userGroupEntity{}) //nolint:exhaustruct
}

// func (x *wrappedDB) WherePairOfUserAndGroup() *wrappedDB {
// 	return x.WhereOrganizationID(&pairOfUserAndGroupEntity{}, x.organizationID) //nolint:exhaustruct
// }
