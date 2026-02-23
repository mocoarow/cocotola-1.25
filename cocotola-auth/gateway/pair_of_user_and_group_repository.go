package gateway

import (
	"context"
	"fmt"
	"strings"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type PairOfUserAndGroupRepository struct {
	dbc      *libgateway.DBConnection
	rbacRepo *RBACRepository
}

var _ service.PairOfUserAndGroupRepository = (*PairOfUserAndGroupRepository)(nil)

func NewPairOfUserAndGroupRepository(ctx context.Context, dbc *libgateway.DBConnection) *PairOfUserAndGroupRepository {
	rbacRepo, err := NewRBACRepository(ctx, dbc)
	if err != nil {
		panic(fmt.Errorf("new rbac repository: %w", err))
	}

	return &PairOfUserAndGroupRepository{
		dbc:      dbc,
		rbacRepo: rbacRepo,
	}
}

func (r *PairOfUserAndGroupRepository) CreatePairOfUserAndGroupBySystemAdmin(ctx context.Context, _ domain.SystemAdminInterface, organizationID *domain.OrganizationID, userID *domain.UserID, userGroupID *domain.UserGroupID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndGroupRepository.CreatePairOfUserAndGroupBySystemAdmin")
	defer span.End()

	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)
	rbacUser := domain.NewRBACUserFromUser(userID)
	rbacRole := domain.NewRBACRoleFromGroup(organizationID, userGroupID)

	return r.addSubjectGroupingPolicy(ctx, rbacDomain, rbacUser, rbacRole)
}

func (r *PairOfUserAndGroupRepository) CreatePairOfUserAndGroup(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, userGroupID *domain.UserGroupID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndGroupRepository.CreatePairOfUserAndGroup")
	defer span.End()

	organizationID := operator.GetOrganizationID()
	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)
	rbacUser := domain.NewRBACUserFromUser(userID)
	rbacRole := domain.NewRBACRoleFromGroup(organizationID, userGroupID)

	return r.addSubjectGroupingPolicy(ctx, rbacDomain, rbacUser, rbacRole)
}

func (r *PairOfUserAndGroupRepository) DeletePairOfUserAndGroup(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, userGroupID *domain.UserGroupID) error {
	_, span := tracer.Start(ctx, "pairOfUserAndGroupRepository.DeletePairOfUserAndGroup")
	defer span.End()

	organizationID := operator.GetOrganizationID()
	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)
	rbacUser := domain.NewRBACUserFromUser(userID)
	rbacRole := domain.NewRBACRoleFromGroup(organizationID, userGroupID)

	return r.removeSubjectGroupingPolicy(ctx, rbacDomain, rbacUser, rbacRole)
}

func (r *PairOfUserAndGroupRepository) FindUserGroupsByUserID(ctx context.Context, operator domain.UserInterface, userID *domain.UserID) ([]*domain.UserGroup, error) {
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

	userGroupRepo := NewUserGroupRepository(r.dbc)
	// userGroupRepo := r.rf.NewUserGroupRepository(ctx)
	result := make([]*domain.UserGroup, 0, len(roles))
	seen := make(map[int]struct{})
	for _, role := range roles {
		if !strings.Contains(role.Role(), ",role:") {
			continue
		}
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

func (r *PairOfUserAndGroupRepository) addSubjectGroupingPolicy(ctx context.Context, rbacDomain libdomain.RBACDomainInterface, child libdomain.RBACSubjectInterface, parent libdomain.RBACSubjectInterface) error {
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

func (r *PairOfUserAndGroupRepository) removeSubjectGroupingPolicy(ctx context.Context, rbacDomain libdomain.RBACDomainInterface, child libdomain.RBACSubjectInterface, parent libdomain.RBACSubjectInterface) error {
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
