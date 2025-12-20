package initialize

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	authdomain "github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	authservice "github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
	libservice "github.com/mocoarow/cocotola-1.25/cocotola-lib/service"
)

func findOrganizationByName(ctx context.Context, systemAdmin authdomain.SystemAdminInterface, mbNonTxManager authservice.TransactionManager, organizationName string) (*authdomain.Organization, error) {
	fn := func(mbrf authservice.RepositoryFactory) (*authdomain.Organization, error) {
		orgRepo := mbrf.NewOrganizationRepository(ctx)
		org, err := orgRepo.FindOrganizationByName(ctx, systemAdmin, organizationName)
		if err != nil {
			if errors.Is(err, authservice.ErrOrganizationNotFound) {
				return nil, fmt.Errorf("organization not found(%s): %w", organizationName, err)
			}
			return nil, fmt.Errorf("find organization by name(%s): %w", organizationName, err)
		}
		return org, nil
	}
	org, err := libservice.Do1(ctx, mbNonTxManager, fn)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return org, nil
}

// func findSystemOwnerByOrganizationID(ctx context.Context, systemAdmin authdomain.SystemAdminInterface, mbNonTxManager authservice.TransactionManager, organizationID *authdomain.OrganizationID) (*authdomain.SystemOwner, error) {
// 	fn := func(rf authservice.RepositoryFactory) (*authdomain.SystemOwner, error) {
// 		userRepo := rf.NewUserRepository(ctx)
// 		sysOwner, err := userRepo.FindSystemOwnerByOrganizationID(ctx, systemAdmin, organizationID)
// 		if err != nil {
// 			return nil, fmt.Errorf("find system owner by organization id(%d): %w", organizationID.Int(), err)
// 		}

// 		return sysOwner, nil
// 	}
// 	sysOwner, err := libservice.Do1(ctx, mbNonTxManager, fn)
// 	if err != nil {
// 		return nil, err //nolint:wrapcheck
// 	}
// 	return sysOwner, nil
// }

func findSystemOwnerByOrganizationName(ctx context.Context, systemAdmin authdomain.SystemAdminInterface, mbNonTxManager authservice.TransactionManager, organizationName string) (*authdomain.SystemOwner, error) {
	fn := func(rf authservice.RepositoryFactory) (*authdomain.SystemOwner, error) {
		userRepo := rf.NewUserRepository(ctx)
		sysOwner, err := userRepo.FindSystemOwnerByOrganizationName(ctx, systemAdmin, organizationName)
		if err != nil {
			return nil, fmt.Errorf("find system owner by organization name(%s): %w", organizationName, err)
		}

		return sysOwner, nil
	}
	sysOwner, err := libservice.Do1(ctx, mbNonTxManager, fn)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return sysOwner, nil
}

func findPublicSpaceByKey(ctx context.Context, systemOwner authdomain.SystemOwnerInterface, nonTxManager authservice.TransactionManager, key string) (*authdomain.Space, error) {
	fn := func(rf authservice.RepositoryFactory) (*authdomain.Space, error) {
		spaceRepo := rf.NewSpaceRepository(ctx)
		publicDefaultSpace, err := spaceRepo.FindPublicSpaceByKey(ctx, systemOwner, key)
		if err != nil {
			return nil, fmt.Errorf("find public default space by key(%s): %w", key, err)
		}

		return publicDefaultSpace, nil
	}
	publicDefaultSpace, err := libservice.Do1(ctx, nonTxManager, fn)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return publicDefaultSpace, nil
}

func findUserByLoginID(ctx context.Context, systemOwner authdomain.SystemOwnerInterface, mbNonTxManager authservice.TransactionManager, loginID string) (*authdomain.User, error) {
	fn := func(mbrf authservice.RepositoryFactory) (*authdomain.User, error) {
		userRepo := mbrf.NewUserRepository(ctx)
		user, err := userRepo.FindUserByLoginID(ctx, systemOwner, loginID)
		if err != nil {
			return nil, fmt.Errorf("find user by login id(%s): %w", loginID, err)
		}

		return user, nil
	}
	user, err := libservice.Do1(ctx, mbNonTxManager, fn)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return user, nil
}

func initTransactionManager(db *gorm.DB, rff func(ctx context.Context, db *gorm.DB) (authservice.RepositoryFactory, error)) (authservice.TransactionManager, error) {
	txManager, err := libgateway.NewTransactionManagerT(db, rff)
	if err != nil {
		return nil, fmt.Errorf("NewTransactionManagerT: %w", err)
	}
	return txManager, nil
}

func initNonTransactionManager(rf authservice.RepositoryFactory) (authservice.TransactionManager, error) {
	nonTxManager, err := libgateway.NewNonTransactionManagerT(rf)
	if err != nil {
		return nil, fmt.Errorf("NewNonTransactionManagerT: %w", err)
	}
	return nonTxManager, nil
}
