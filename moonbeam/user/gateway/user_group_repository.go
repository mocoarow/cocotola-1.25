package gateway

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	libdomain "github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/moonbeam/lib/gateway"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
)

type userGroupEntity struct {
	BaseModelEntity
	ID             int
	OrganizationID int
	KeyName        string
	Name           string
	Description    string
	Deleted        bool
}

func (e *userGroupEntity) TableName() string {
	return UserGroupTableName
}

func (e *userGroupEntity) toUserGroup() (*domain.UserGroup, error) {
	baseModel, err := e.ToBaseModel()
	if err != nil {
		return nil, fmt.Errorf("toBaseModel: %w", err)
	}

	userGroupID, err := domain.NewUserGroupID(e.ID)
	if err != nil {
		return nil, fmt.Errorf("domain.NewUser: %w", err)
	}

	organizationID, err := domain.NewOrganizationID(e.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("domain.NewOrganizationID: %w", err)
	}

	userGroupModel, err := domain.NewUserGroup(baseModel, userGroupID, organizationID, e.KeyName, e.Name, e.Description)
	if err != nil {
		return nil, fmt.Errorf("domain.NewUserGroup: %w", err)
	}

	return userGroupModel, nil
}

type UserGroupRepository struct {
	dialect libgateway.DialectRDBMS
	db      *gorm.DB
	logger  *slog.Logger
}

var _ service.UserGroupRepository = (*UserGroupRepository)(nil)

func NewUserGroupRepository(_ context.Context, dialect libgateway.DialectRDBMS, db *gorm.DB) service.UserGroupRepository {
	return &UserGroupRepository{
		dialect: dialect,
		db:      db,
		logger:  slog.Default().With(slog.String(libdomain.LoggerNameKey, "UserGroupRepository")),
	}
}

func (r *UserGroupRepository) FindAllUserGroups(ctx context.Context, operator domain.UserInterface) ([]*domain.UserGroup, error) {
	_, span := tracer.Start(ctx, "UserGroupRepository.FindAllUserGroups")
	defer span.End()

	userGroups := []userGroupEntity{}
	if result := r.db.Where(&userGroupEntity{ //nolint:exhaustruct
		OrganizationID: operator.GetOrganizationID().Int(),
	}).Find(&userGroups); result.Error != nil {
		return nil, result.Error
	}

	userGroupModels := make([]*domain.UserGroup, len(userGroups))
	for i, e := range userGroups {
		m, err := e.toUserGroup()
		if err != nil {
			return nil, fmt.Errorf("toUserGroup: %w", err)
		}
		userGroupModels[i] = m
	}

	return userGroupModels, nil
}

func (r *UserGroupRepository) FindSystemOwnerGroup(ctx context.Context, _ domain.SystemAdminInterface, organizationID *domain.OrganizationID) (*domain.UserGroup, error) {
	_, span := tracer.Start(ctx, "UserGroupRepository.FindSystemOwnerGroup")
	defer span.End()

	var userGroup userGroupEntity
	if result := r.db.Where(&userGroupEntity{ //nolint:exhaustruct
		OrganizationID: organizationID.Int(),
		KeyName:        service.SystemOwnerGroupKey,
	}).First(&userGroup); result.Error != nil {
		return nil, result.Error
	}

	return userGroup.toUserGroup()
}

func (r *UserGroupRepository) FindUserGroupByID(ctx context.Context, operator domain.UserInterface, userGroupID *domain.UserGroupID) (*domain.UserGroup, error) {
	_, span := tracer.Start(ctx, "UserGroupRepository.FindUserGroupByID")
	defer span.End()

	var userGroup userGroupEntity
	if result := r.db.Where("organization_id = ?", operator.GetOrganizationID().Int()).
		Where("id = ? and deleted = ?", userGroupID.Int(), r.dialect.BoolDefaultValue()).
		First(&userGroup); result.Error != nil {
		return nil, result.Error
	}

	return userGroup.toUserGroup()
}

func (r *UserGroupRepository) FindUserGroupByKey(ctx context.Context, operator domain.UserInterface, key string) (*domain.UserGroup, error) {
	_, span := tracer.Start(ctx, "UserGroupRepository.FindUserGroupByKey")
	defer span.End()

	var userGroup userGroupEntity
	if result := r.db.Where("organization_id = ?", operator.GetOrganizationID().Int()).
		Where("key_name = ? and deleted = ?", key, r.dialect.BoolDefaultValue()).
		First(&userGroup); result.Error != nil {
		return nil, result.Error
	}

	return userGroup.toUserGroup()
}

func (r *UserGroupRepository) createUserGroup(userID *domain.UserID, organizationID *domain.OrganizationID, key, name string) (*domain.UserGroupID, error) {
	r.logger.InfoContext(context.Background(), "createUserGroup", "key", key, "name", name, "organizationID", organizationID.Int(), "userID", userID.Int())
	userGroup := userGroupEntity{ //nolint:exhaustruct
		BaseModelEntity: BaseModelEntity{ //nolint:exhaustruct
			Version:   1,
			CreatedBy: userID.Int(),
			UpdatedBy: userID.Int(),
		},
		OrganizationID: organizationID.Int(),
		KeyName:        key,
		Name:           name,
	}
	if result := r.db.Create(&userGroup); result.Error != nil {
		return nil, fmt.Errorf("create user group(%s): %w", key, libgateway.ConvertDuplicatedError(result.Error, service.ErrUserGroupAlreadyExists))
	}

	userGroupID, err := domain.NewUserGroupID(userGroup.ID)
	if err != nil {
		return nil, fmt.Errorf("NewUserGroupID: %w", err)
	}

	return userGroupID, nil
}

func (r *UserGroupRepository) CreateSystemOwnerGroup(ctx context.Context, operator domain.SystemAdminInterface, organizationID *domain.OrganizationID) (*domain.UserGroupID, error) {
	_, span := tracer.Start(ctx, "UserGroupRepository.CreateSystemOwnerGroup")
	defer span.End()

	r.logger.InfoContext(context.Background(), "CreateSystemOwnerGroup", "organizationID", organizationID.Int())
	userGroupID, err := r.createUserGroup(operator.GetUserID(), organizationID, service.SystemOwnerGroupKey, service.SystemOwnerGroupName)
	if err != nil {
		return nil, fmt.Errorf("createUserGroup: %w", err)
	}

	return userGroupID, nil
}

func (r *UserGroupRepository) CreateOwnerGroup(ctx context.Context, operator domain.SystemOwnerInterface, organizationID *domain.OrganizationID) (*domain.UserGroupID, error) {
	_, span := tracer.Start(ctx, "UserGroupRepository.CreateOwnerGroup")
	defer span.End()

	r.logger.InfoContext(ctx, "CreateOwnerGroup", "organizationID", organizationID.Int())
	userGroupID, err := r.createUserGroup(operator.GetUserID(), organizationID, service.OwnerGroupKey, service.OwnerGroupName)
	if err != nil {
		return nil, fmt.Errorf("createUserGroup: %w", err)
	}

	return userGroupID, nil
}

func (r *UserGroupRepository) CreatePublicGroup(ctx context.Context, operator domain.SystemOwnerInterface, organizationID *domain.OrganizationID) (*domain.UserGroupID, error) {
	_, span := tracer.Start(ctx, "UserGroupRepository.CreatePublicGroup")
	defer span.End()

	r.logger.InfoContext(context.Background(), "CreatePublicGroup", "organizationID", organizationID.Int())
	userGroupID, err := r.createUserGroup(operator.GetUserID(), organizationID, service.PublicGroupKey, service.PublicGroupName)
	if err != nil {
		return nil, fmt.Errorf("createUserGroup: %w", err)
	}

	return userGroupID, nil
}

func (r *UserGroupRepository) AddUserGroup(ctx context.Context, operator domain.OwnerInterface, param *service.AddUserGroupParameter) (*domain.UserGroupID, error) {
	_, span := tracer.Start(ctx, "UserGroupRepository.AddUserGroup")
	defer span.End()

	userGroup := userGroupEntity{ //nolint:exhaustruct
		BaseModelEntity: BaseModelEntity{ //nolint:exhaustruct
			Version:   1,
			CreatedBy: operator.GetUserID().Int(),
			UpdatedBy: operator.GetUserID().Int(),
		},
		OrganizationID: operator.GetOrganizationID().Int(),
		KeyName:        param.Key,
		Name:           param.Name,
		Description:    param.Description,
	}
	if result := r.db.Create(&userGroup); result.Error != nil {
		return nil, fmt.Errorf(": %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrUserGroupAlreadyExists))
	}

	userGroupID, err := domain.NewUserGroupID(userGroup.ID)
	if err != nil {
		return nil, fmt.Errorf("NewUserGroupID: %w", err)
	}

	return userGroupID, nil
}
