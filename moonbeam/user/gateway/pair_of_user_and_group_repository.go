package gateway

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	libgateway "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

type pairOfUserAndGroupRepository struct {
	dialect  libgateway.DialectRDBMS
	db       *gorm.DB
	rf       service.RepositoryFactory
	rbacRepo service.RBACRepository
}

var _ service.PairOfUserAndGroupRepository = (*pairOfUserAndGroupRepository)(nil)

func NewPairOfUserAndGroupRepository(ctx context.Context, dialect libgateway.DialectRDBMS, db *gorm.DB, rf service.RepositoryFactory) service.PairOfUserAndGroupRepository {
	rbacRepo, err := NewRBACRepository(ctx, db)
	if err != nil {
		panic(fmt.Errorf("new rbac repository: %w", err))
	}

	return &pairOfUserAndGroupRepository{
		dialect:  dialect,
		db:       db,
		rf:       rf,
		rbacRepo: rbacRepo,
	}
}

func (r *pairOfUserAndGroupRepository) CreatePairOfUserAndGroupBySystemAdmin(ctx context.Context, _ domain.SystemAdminInterface, organizationID *domain.OrganizationID, userID *domain.UserID, userGroupID *domain.UserGroupID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndGroupRepository.CreatePairOfUserAndGroupBySystemAdmin")
	defer span.End()

	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)
	rbacUser := domain.NewRBACUserFromUser(userID)
	rbacRole := domain.NewRBACRoleFromGroup(organizationID, userGroupID)

	return r.addSubjectGroupingPolicy(ctx, rbacDomain, rbacUser, rbacRole)
}

func (r *pairOfUserAndGroupRepository) CreatePairOfUserAndGroup(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, userGroupID *domain.UserGroupID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndGroupRepository.CreatePairOfUserAndGroup")
	defer span.End()

	organizationID := operator.GetOrganizationID()
	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)
	rbacUser := domain.NewRBACUserFromUser(userID)
	rbacRole := domain.NewRBACRoleFromGroup(organizationID, userGroupID)

	return r.addSubjectGroupingPolicy(ctx, rbacDomain, rbacUser, rbacRole)
}

func (r *pairOfUserAndGroupRepository) DeletePairOfUserAndGroup(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, userGroupID *domain.UserGroupID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndGroupRepository.DeletePairOfUserAndGroup")
	defer span.End()

	organizationID := operator.GetOrganizationID()
	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)
	rbacUser := domain.NewRBACUserFromUser(userID)
	rbacRole := domain.NewRBACRoleFromGroup(organizationID, userGroupID)

	return r.removeSubjectGroupingPolicy(ctx, rbacDomain, rbacUser, rbacRole)
}

func (r *pairOfUserAndGroupRepository) FindUserGroupsByUserID(ctx context.Context, operator domain.UserInterface, userID *domain.UserID) ([]*domain.UserGroup, error) {
	_, span := tracer.Start(ctx, "pairOfUserAndGroupRepository.FindUserGroupsByUserID")
	defer span.End()

	organizationID := operator.GetOrganizationID()
	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)
	rbacUser := domain.NewRBACUserFromUser(userID)

	roles, err := r.rbacRepo.GetGroupsForSubject(ctx, rbacDomain, rbacUser)
	if err != nil {
		return nil, fmt.Errorf("rbacRepo.GetGroupsForSubject: %w", err)
	}

	if len(roles) == 0 {
		return []*domain.UserGroup{}, nil
	}

	userGroupRepo := r.rf.NewUserGroupRepository(ctx)
	result := make([]*domain.UserGroup, 0, len(roles))
	seen := make(map[int]struct{})
	for _, role := range roles {
		orgID, userGroupID, err := domain.NewOrganizationAndUserGroupIDsFromRole(role)
		if err != nil {
			return nil, fmt.Errorf("domain.NewOrganizationAndUserGroupIDsFromRole: %w", err)
		}
		if orgID.Int() != organizationID.Int() {
			continue
		}
		if _, exists := seen[userGroupID.Int()]; exists {
			continue
		}
		seen[userGroupID.Int()] = struct{}{}

		userGroup, err := userGroupRepo.FindUserGroupByID(ctx, operator, userGroupID)
		if err != nil {
			return nil, fmt.Errorf("userGroupRepo.FindUserGroupByID: %w", err)
		}
		result = append(result, userGroup)
	}

	return result, nil
}

func (r *pairOfUserAndGroupRepository) addSubjectGroupingPolicy(ctx context.Context, rbacDomain domain.RBACDomain, child domain.RBACSubject, parent domain.RBACSubject) error {
	roles, err := r.rbacRepo.GetGroupsForSubject(ctx, rbacDomain, child)
	if err != nil {
		return fmt.Errorf("rbacRepo.GetGroupsForSubject: %w", err)
	}

	for _, role := range roles {
		if role.Role() == parent.Subject() {
			return service.ErrPairOfUserAndGroupAlreadyExists
		}
	}

	if err := r.rbacRepo.CreateSubjectGroupingPolicy(ctx, rbacDomain, child, parent); err != nil {
		return fmt.Errorf("rbacRepo.CreateSubjectGroupingPolicy: %w", err)
	}

	return nil
}

func (r *pairOfUserAndGroupRepository) removeSubjectGroupingPolicy(ctx context.Context, rbacDomain domain.RBACDomain, child domain.RBACSubject, parent domain.RBACSubject) error {
	roles, err := r.rbacRepo.GetGroupsForSubject(ctx, rbacDomain, child)
	if err != nil {
		return fmt.Errorf("rbacRepo.GetGroupsForSubject: %w", err)
	}

	found := false
	for _, role := range roles {
		if role.Role() == parent.Subject() {
			found = true
			break
		}
	}

	if !found {
		return service.ErrPairOfUserAndGroupNotFound
	}

	if err := r.rbacRepo.DeleteSubjectGroupingPolicy(ctx, rbacDomain, child, parent); err != nil {
		return fmt.Errorf("rbacRepo.DeleteSubjectGroupingPolicy: %w", err)
	}

	return nil
}

// helper removed: parsing is centralized in domain.NewOrganizationAndUserGroupIDsFromRole
