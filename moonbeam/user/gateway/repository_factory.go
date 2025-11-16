package gateway

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	libdomain "github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

type repositoryFactory struct {
	dialect    libgateway.DialectRDBMS
	driverName string
	db         *gorm.DB
	location   *time.Location
}

func NewRepositoryFactory(_ context.Context, dialect libgateway.DialectRDBMS, driverName string, db *gorm.DB, location *time.Location) (service.RepositoryFactory, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil. err: %w", libdomain.ErrInvalidArgument)
	}

	return &repositoryFactory{
		dialect:    dialect,
		driverName: driverName,
		db:         db,
		location:   location,
	}, nil
}

func (f *repositoryFactory) NewOrganizationRepository(ctx context.Context) service.OrganizationRepository {
	return NewOrganizationRepository(ctx, f.db)
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
