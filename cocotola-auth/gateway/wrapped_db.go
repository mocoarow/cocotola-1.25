package gateway

import (
	"gorm.io/gorm"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
)

type HasTableName interface {
	TableName() string
}

type wrappedDB struct {
	dbc            *libgateway.DBConnection
	db             *gorm.DB
	organizationID *domain.OrganizationID
}

func newWrappedDB(dbc *libgateway.DBConnection, organizationID *domain.OrganizationID) *wrappedDB {
	return &wrappedDB{
		dbc:            dbc,
		db:             dbc.DB,
		organizationID: organizationID,
	}
}

func (x *wrappedDB) Table(name string, args ...any) *wrappedDB {
	x.db = x.db.Table(name, args...)

	return x
}

func (x *wrappedDB) Select(query any, args ...any) *wrappedDB {
	x.db = x.db.Select(query, args...)

	return x
}

func (x *wrappedDB) Where(query any, args ...any) *wrappedDB {
	x.db = x.db.Where(query, args...)

	return x
}

func (x *wrappedDB) Joins(query string, args ...any) *wrappedDB {
	x.db = x.db.Joins(query, args...)

	return x
}

func (x *wrappedDB) WhereOrganizationID(table HasTableName, organizationID *domain.OrganizationID) *wrappedDB {
	x.db = x.db.Where(table.TableName()+".organization_id = ?", organizationID.Int())

	return x
}

func (x *wrappedDB) WhereNotDeleted(table HasTableName) *wrappedDB {
	x.db = x.db.Where(table.TableName()+".deleted = ?", x.dbc.Dialect.BoolDefaultValue())

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
