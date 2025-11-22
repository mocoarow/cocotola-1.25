package gateway

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

type spaceManager struct {
	dialect libgateway.DialectRDBMS
	db      *gorm.DB
	rf      service.RepositoryFactory
}

var _ service.SpaceManager = (*spaceManager)(nil)

func NewSpaceManager(_ context.Context, dialect libgateway.DialectRDBMS, db *gorm.DB, rf service.RepositoryFactory) (service.SpaceManager, error) {
	return &spaceManager{dialect: dialect, db: db, rf: rf}, nil
}

func (m *spaceManager) CreatePersonalSpace(ctx context.Context, operator domain.UserInterface, param *service.CreatePersonalSpaceParameter) (*domain.SpaceID, error) {
	ctx, span := tracer.Start(ctx, "spaceManager.CreatePersonalSpace")
	defer span.End()

	userRepo := m.rf.NewUserRepository(ctx)
	targetUser, err := userRepo.FindUserByID(ctx, operator, param.UserID)
	if err != nil {
		return nil, fmt.Errorf("FindUserByID: %w", err)
	}

	spaceRepo := m.rf.NewSpaceRepository(ctx)
	createParam := &service.CreateSpaceParameter{ //nolint:exhaustruct
		Key:       param.KeyName,
		Name:      param.Name,
		SpaceType: "personal",
	}
	spaceID, err := spaceRepo.CreateSpace(ctx, targetUser, createParam)
	if err != nil {
		return nil, fmt.Errorf("CreateSpace: %w", err)
	}

	pairRepo := NewPairOfUserAndSpaceRepository(ctx, m.dialect, m.db, m.rf)
	if err := pairRepo.CreatePairOfUserAndSpace(ctx, targetUser, targetUser.GetUserID(), spaceID); err != nil {
		return nil, fmt.Errorf("CreatePairOfUserAndSpace: %w", err)
	}

	return spaceID, nil
}

func (m *spaceManager) AddUserToSpace(ctx context.Context, operator domain.SystemOwnerInterface, userID domain.UserID, spaceID *domain.SpaceID) error {
	ctx, span := tracer.Start(ctx, "spaceManager.AddUserToSpace")
	defer span.End()

	var space spaceEntity
	if err := m.db.WithContext(ctx).
		Where("organization_id = ?", operator.GetOrganizationID().Int()).
		Where("id = ?", spaceID.Int()).
		First(&space).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrSpaceNotFound
		}
		return fmt.Errorf("find space: %w", err)
	}

	pairRepo := NewPairOfUserAndSpaceRepository(ctx, m.dialect, m.db, m.rf)
	userIDCopy := userID
	if err := pairRepo.CreatePairOfUserAndSpace(ctx, operator, &userIDCopy, spaceID); err != nil {
		return fmt.Errorf("CreatePairOfUserAndSpace: %w", err)
	}

	return nil
}

func (m *spaceManager) GetPersonalSpace(ctx context.Context, operator domain.UserInterface) (*domain.Space, error) {
	pairRepo := NewPairOfUserAndSpaceRepository(ctx, m.dialect, m.db, m.rf)
	spaces, err := pairRepo.FindMySpaces(ctx, operator)
	if err != nil {
		return nil, fmt.Errorf("FindMySpaces: %w", err)
	}

	for _, space := range spaces {
		if space.SpaceType == "personal" {
			return space, nil
		}
	}

	return nil, service.ErrSpaceNotFound
}
