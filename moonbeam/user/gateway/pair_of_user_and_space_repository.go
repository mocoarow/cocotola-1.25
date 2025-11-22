package gateway

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

type pairOfUserAndSpaceRepository struct {
	dialect  libgateway.DialectRDBMS
	db       *gorm.DB
	rf       service.RepositoryFactory
	rbacRepo service.RBACRepository
}

var _ service.PairOfUserAndSpaceRepository = (*pairOfUserAndSpaceRepository)(nil)

func NewPairOfUserAndSpaceRepository(ctx context.Context, dialect libgateway.DialectRDBMS, db *gorm.DB, rf service.RepositoryFactory) service.PairOfUserAndSpaceRepository {
	rbacRepo, err := NewRBACRepository(ctx, db)
	if err != nil {
		panic(fmt.Errorf("new rbac repository: %w", err))
	}

	return &pairOfUserAndSpaceRepository{
		dialect:  dialect,
		db:       db,
		rf:       rf,
		rbacRepo: rbacRepo,
	}
}

func (r *pairOfUserAndSpaceRepository) CreatePairOfUserAndSpace(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, spaceID *domain.SpaceID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndSpaceRepository.CreatePairOfUserAndSpace")
	defer span.End()

	rbacDomain := domain.NewRBACDomainFromOrganization(operator.GetOrganizationID())
	rbacUser := domain.NewRBACUserFromUser(userID)
	rbacSpace := domain.NewRBACRoleFromSpace(operator.GetOrganizationID(), spaceID)

	roles, err := r.rbacRepo.GetGroupsForSubject(ctx, rbacDomain, rbacUser)
	if err != nil {
		return fmt.Errorf("rbacRepo.GetGroupsForSubject: %w", err)
	}
	for _, role := range roles {
		if role.Role() == rbacSpace.Role() {
			return service.ErrPairOfUserAndSpaceAlreadyExists
		}
	}

	if err := r.rbacRepo.CreateSubjectGroupingPolicy(ctx, rbacDomain, rbacUser, rbacSpace); err != nil {
		return fmt.Errorf("rbacRepo.CreateSubjectGroupingPolicy: %w", err)
	}

	return nil
}

func (r *pairOfUserAndSpaceRepository) FindMySpaces(ctx context.Context, operator domain.UserInterface) ([]*domain.Space, error) {
	_, span := tracer.Start(ctx, "pairOfUserAndSpaceRepository.FindMySpaces")
	defer span.End()

	rbacDomain := domain.NewRBACDomainFromOrganization(operator.GetOrganizationID())
	rbacUser := domain.NewRBACUserFromUser(operator.GetUserID())

	roles, err := r.rbacRepo.GetGroupsForSubject(ctx, rbacDomain, rbacUser)
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
	if err := r.db.WithContext(ctx).
		Where("organization_id = ?", operator.GetOrganizationID().Int()).
		Where("id IN ?", spaceIDs).
		Find(&spacesE).Error; err != nil {
		return nil, fmt.Errorf("find spaces: %w", err)
	}

	spaces, err := spacesE.toSpaces()
	if err != nil {
		return nil, fmt.Errorf("spacesE.toSpaces: %w", err)
	}

	return spaces, nil
}
