package gateway

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type authorizationManager struct {
	dialect  libgateway.DialectRDBMS
	db       *gorm.DB
	rf       service.RepositoryFactory
	rbacRepo service.RBACRepository
	pairRepo service.PairOfUserAndGroupRepository
}

func NewAuthorizationManager(ctx context.Context, dialect libgateway.DialectRDBMS, db *gorm.DB, rf service.RepositoryFactory) (service.AuthorizationManager, error) {
	rbacRepo, err := NewRBACRepository(ctx, db)
	if err != nil {
		return nil, err
	}
	pairRepo := NewPairOfUserAndGroupRepository(ctx, dialect, db, rf)

	return &authorizationManager{
		dialect:  dialect,
		db:       db,
		rf:       rf,
		rbacRepo: rbacRepo,
		pairRepo: pairRepo,
	}, nil
}

// func (m *authorizationManager) Init(ctx context.Context) error {
// 	rbacRepo, err := newRBACRepository(ctx, m.db)
// 	if err != nil {
// 		return err
// 	}
// 	m.rbacRepo = rbacRepo
// 	return m.rbacRepo.Init()
// }

func (m *authorizationManager) AddUserToGroupBySystemAdmin(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID, userID *domain.UserID, userGroupID *domain.UserGroupID) error {
	if err := m.pairRepo.CreatePairOfUserAndGroupBySystemAdmin(ctx, operator, organizationID, userID, userGroupID); err != nil {
		return fmt.Errorf("CreatePairOfUserAndGroupBySystemAdmin: %w", err)
	}

	return nil
}

func (m *authorizationManager) AddUserToGroup(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, userGroupID *domain.UserGroupID) error {
	if err := m.pairRepo.CreatePairOfUserAndGroup(ctx, operator, userID, userGroupID); err != nil {
		return fmt.Errorf("CreatePairOfUserAndGroup: %w", err)
	}

	return nil
}

func (m *authorizationManager) AttachPolicyToUser(ctx context.Context, operator domain.UserInterface, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error {
	ctx, span := tracer.Start(ctx, "authorizationManager.AttachPolicyToUser")
	defer span.End()

	rbacDomain := domain.NewRBACDomainFromOrganization(operator.GetOrganizationID())

	if err := m.rbacRepo.CreatePolicy(ctx, rbacDomain, subject, action, object, effect); err != nil {
		return fmt.Errorf("CreatePolicy: %w", err)
	}

	return nil
}

func (m *authorizationManager) AttachPolicyToUserBySystemAdmin(ctx context.Context, _ domain.SystemAdminInterface, organizationID *domain.OrganizationID, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error {
	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)

	if err := m.rbacRepo.CreatePolicy(ctx, rbacDomain, subject, action, object, effect); err != nil {
		return fmt.Errorf("CreatePolicy: %w", err)
	}

	return nil
}
func (m *authorizationManager) AttachPolicyToUserBySystemOwner(ctx context.Context, operator domain.SystemOwnerInterface, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error {
	organizationID := operator.GetOrganizationID()
	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)

	if err := m.rbacRepo.CreatePolicy(ctx, rbacDomain, subject, action, object, effect); err != nil {
		return fmt.Errorf("CreatePolicy: %w", err)
	}

	return nil
}

func (m *authorizationManager) AttachPolicyToGroup(ctx context.Context, operator domain.UserInterface, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error {
	rbacDomain := domain.NewRBACDomainFromOrganization(operator.GetOrganizationID())

	if err := m.rbacRepo.CreatePolicy(ctx, rbacDomain, subject, action, object, effect); err != nil {
		return fmt.Errorf("CreatePolicy: %w", err)
	}

	return nil
}

func (m *authorizationManager) AttachPolicyToGroupBySystemAdmin(ctx context.Context, _ domain.SystemAdminInterface, organizationID *domain.OrganizationID, subject libdomain.RBACSubject, action libdomain.RBACAction, object libdomain.RBACObject, effect libdomain.RBACEffect) error {
	rbacDomain := domain.NewRBACDomainFromOrganization(organizationID)

	if err := m.rbacRepo.CreatePolicy(ctx, rbacDomain, subject, action, object, effect); err != nil {
		return fmt.Errorf("CreatePolicy: %w", err)
	}

	return nil
}

func (m *authorizationManager) AddObjectToObject(ctx context.Context, operator domain.SystemOwnerInterface, child, parent libdomain.RBACObject) error {
	rbacDomain := domain.NewRBACDomainFromOrganization(operator.GetOrganizationID())

	if err := m.rbacRepo.CreateObjectGroupingPolicy(ctx, rbacDomain, child, parent); err != nil {
		return fmt.Errorf("CreateObjectGroupingPolicy. priv: read: %w", err)
	}

	return nil
}

func (m *authorizationManager) CheckAuthorization(ctx context.Context, operator domain.UserInterface, rbacAction libdomain.RBACAction, rbacObject libdomain.RBACObject) (bool, error) {
	rbacDomain := domain.NewRBACDomainFromOrganization(operator.GetOrganizationID())

	userGroups, err := m.pairRepo.FindUserGroupsByUserID(ctx, operator, operator.GetUserID())
	if err != nil {
		return false, fmt.Errorf("FindUserGroupsByUserID: %w", err)
	}

	rbacRoles := make([]libdomain.RBACRole, 0, len(userGroups))
	for _, userGroup := range userGroups {
		rbacRoles = append(rbacRoles, domain.NewRBACRoleFromGroup(operator.GetOrganizationID(), userGroup.UserGroupID))
	}

	rbacOperator := domain.NewRBACUserFromUser(operator.GetUserID())
	e, err := m.rbacRepo.NewEnforcerWithGroupsAndUsers(ctx, rbacRoles, []libdomain.RBACUser{rbacOperator})
	if err != nil {
		return false, fmt.Errorf("NewEnforcerWithGroupsAndUsers: %w", err)
	}

	ok, err := e.Enforce(rbacOperator.Subject(), rbacObject.Object(), rbacAction.Action(), rbacDomain.Domain())
	if err != nil {
		return false, fmt.Errorf("enforce: %w", err)
	}

	return ok, nil
}
