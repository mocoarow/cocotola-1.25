package gateway

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

type spaceManager struct {
	dialect  libgateway.DialectRDBMS
	db       *gorm.DB
	rf       service.RepositoryFactory
	rbacRepo service.RBACRepository
}

var _ service.SpaceManager = (*spaceManager)(nil)

func NewSpaceManager(ctx context.Context, dialect libgateway.DialectRDBMS, db *gorm.DB, rf service.RepositoryFactory) (service.SpaceManager, error) {
	rbacRepo, err := NewRBACRepository(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("new rbac repository: %w", err)
	}

	return &spaceManager{
		dialect:  dialect,
		db:       db,
		rf:       rf,
		rbacRepo: rbacRepo,
	}, nil
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

	if err := m.attachUserToSpace(ctx, targetUser.GetOrganizationID(), targetUser.GetUserID(), spaceID); err != nil {
		return nil, fmt.Errorf("attachUserToSpace: %w", err)
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

	userIDCopy := userID
	userRepo := m.rf.NewUserRepository(ctx)
	targetUser, err := userRepo.FindUserByID(ctx, operator, &userIDCopy)
	if err != nil {
		return fmt.Errorf("FindUserByID: %w", err)
	}
	if targetUser.GetOrganizationID().Int() != operator.GetOrganizationID().Int() {
		return service.ErrUserNotFound
	}

	if err := m.attachUserToSpace(ctx, targetUser.GetOrganizationID(), targetUser.GetUserID(), spaceID); err != nil {
		return fmt.Errorf("attachUserToSpace: %w", err)
	}

	return nil
}

func (m *spaceManager) GetPersonalSpace(ctx context.Context, operator domain.UserInterface) (*domain.Space, error) {
	spaces, err := m.findSpacesByUser(ctx, operator)
	if err != nil {
		return nil, fmt.Errorf("findSpacesByUser: %w", err)
	}

	for _, space := range spaces {
		if space.SpaceType == "personal" {
			return space, nil
		}
	}

	return nil, service.ErrSpaceNotFound
}

func (m *spaceManager) attachUserToSpace(ctx context.Context, organizationID *domain.OrganizationID, userID *domain.UserID, spaceID *domain.SpaceID) error {
	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)
	rbacUser := domain.NewRBACUserFromUser(userID)
	rbacSpace := domain.NewRBACRoleFromSpace(organizationID, spaceID)

	roles, err := m.rbacRepo.GetGroupsForSubject(ctx, rbacDomain, rbacUser)
	if err != nil {
		return fmt.Errorf("rbacRepo.GetGroupsForSubject: %w", err)
	}

	for _, role := range roles {
		if role.Role() == rbacSpace.Role() {
			return nil
		}
	}

	if err := m.rbacRepo.CreateSubjectGroupingPolicy(ctx, rbacDomain, rbacUser, rbacSpace); err != nil {
		return fmt.Errorf("rbacRepo.CreateSubjectGroupingPolicy: %w", err)
	}

	return nil
}

func (m *spaceManager) findSpacesByUser(ctx context.Context, operator domain.UserInterface) ([]*domain.Space, error) {
	rbacDomain := domain.NewRBACDomainFromOrganization(operator.GetOrganizationID())
	rbacUser := domain.NewRBACUserFromUser(operator.GetUserID())

	roles, err := m.rbacRepo.GetGroupsForSubject(ctx, rbacDomain, rbacUser)
	if err != nil {
		return nil, fmt.Errorf("rbacRepo.GetGroupsForSubject: %w", err)
	}

	spaceIDs := make([]int, 0, len(roles))
	seen := make(map[int]struct{})
	for _, role := range roles {
		if !strings.Contains(role.Role(), ",space:") {
			continue
		}
		orgID, spaceID, err := domain.NewOrganizationAndSpaceIDsFromRole(role)
		if err != nil {
			return nil, fmt.Errorf("domain.NewOrganizationAndSpaceIDsFromRole: %w", err)
		}
		if orgID.Int() != operator.GetOrganizationID().Int() {
			continue
		}
		if _, exists := seen[spaceID.Int()]; exists {
			continue
		}
		seen[spaceID.Int()] = struct{}{}
		spaceIDs = append(spaceIDs, spaceID.Int())
	}

	if len(spaceIDs) == 0 {
		return []*domain.Space{}, nil
	}

	var spacesE spaceEntities
	if err := m.db.WithContext(ctx).
		Where("organization_id = ?", operator.GetOrganizationID().Int()).
		Where("deleted = ?", m.dialect.BoolDefaultValue()).
		Where("id IN ?", spaceIDs).
		Order("key_name").
		Find(&spacesE).Error; err != nil {
		return nil, fmt.Errorf("find spaces: %w", err)
	}

	spaces, err := spacesE.toSpaces()
	if err != nil {
		return nil, fmt.Errorf("spacesE.toSpaces: %w", err)
	}

	return spaces, nil
}
