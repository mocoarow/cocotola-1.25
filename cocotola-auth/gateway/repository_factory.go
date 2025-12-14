package gateway

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
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
func (f *repositoryFactory) NewUserRepository(ctx context.Context) service.UserRepository {
	return NewUserRepository(ctx, f.dialect, f.db, f)
}

func (f *repositoryFactory) NewUserGroupRepository(ctx context.Context) service.UserGroupRepository {
	return NewUserGroupRepository(ctx, f.dialect, f.db)
}

func (f *repositoryFactory) NewSpaceRepository(ctx context.Context) service.SpaceRepository {
	return NewSpaceRepository(ctx, f.dialect, f.db)
}

func (f *repositoryFactory) NewSpaceManager(ctx context.Context) (service.SpaceManager, error) {
	return NewSpaceManager(ctx, f.dialect, f.db, f)
}

func (f *repositoryFactory) NewAuthorizationManager(ctx context.Context) (service.AuthorizationManager, error) {
	return NewAuthorizationManager(ctx, f.dialect, f.db, f)
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
