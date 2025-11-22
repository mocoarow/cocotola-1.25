package service

import (
	"context"
	"errors"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
)

var ErrPairOfUserAndGroupAlreadyExists = errors.New("pair of user and group already exists")
var ErrPairOfUserAndGroupNotFound = errors.New("pair of user and group not found")

type PairOfUserAndGroupRepository interface {
	CreatePairOfUserAndGroupBySystemAdmin(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID, userID *domain.UserID, userGroupID *domain.UserGroupID) error

	CreatePairOfUserAndGroup(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, userGroupID *domain.UserGroupID) error

	DeletePairOfUserAndGroup(ctx context.Context, operator domain.UserInterface, userID *domain.UserID, userGroupID *domain.UserGroupID) error

	FindUserGroupsByUserID(ctx context.Context, operator domain.UserInterface, userID *domain.UserID) ([]*domain.UserGroup, error)
}
