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

type organizationEntity struct {
	BaseModelEntity
	ID   int
	Name string
}

func (e *organizationEntity) TableName() string {
	return OrganizationTableName
}

func (e *organizationEntity) toModel() (*domain.Organization, error) {
	baseModel, err := e.ToBaseModel()
	if err != nil {
		return nil, fmt.Errorf("to base model: %w", err)
	}

	organizationID, err := domain.NewOrganizationID(e.ID)
	if err != nil {
		return nil, fmt.Errorf("new organization ID: %w", err)
	}

	organization, err := domain.NewOrganization(baseModel, organizationID, e.Name)
	if err != nil {
		return nil, fmt.Errorf("new organization: %w", err)
	}

	return organization, nil
}

type OrganizationRepository struct {
	db *gorm.DB
}

var _ service.OrganizationRepository = (*OrganizationRepository)(nil)

func NewOrganizationRepository(_ context.Context, db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{
		db: db,
	}
}

func (r *OrganizationRepository) GetOrganization(ctx context.Context, operator domain.UserInterface) (*domain.Organization, error) {
	_, span := tracer.Start(ctx, "organizationRepository.GetOrganization")
	defer span.End()

	var organization organizationEntity
	if result := r.db.WithContext(ctx).Where(organizationEntity{ //nolint:exhaustruct
		ID: operator.GetOrganizationID().Int(),
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrOrganizationNotFound
		}

		return nil, result.Error
	}

	return organization.toModel()
}

func (r *OrganizationRepository) FindOrganizationByName(ctx context.Context, _ domain.SystemAdminInterface, name string) (*domain.Organization, error) {
	_, span := tracer.Start(ctx, "organizationRepository.FindOrganizationByName")
	defer span.End()

	var organization organizationEntity
	if result := r.db.WithContext(ctx).Where(organizationEntity{ //nolint:exhaustruct
		Name: name,
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrOrganizationNotFound
		}

		return nil, result.Error
	}

	return organization.toModel()
}

func (r *OrganizationRepository) FindOrganizationByID(ctx context.Context, _ domain.SystemAdminInterface, id *domain.OrganizationID) (*domain.Organization, error) {
	_, span := tracer.Start(ctx, "organizationRepository.FindOrganizationByID")
	defer span.End()

	var organization organizationEntity
	if result := r.db.WithContext(ctx).Where(organizationEntity{ //nolint:exhaustruct
		ID: id.Int(),
	}).First(&organization); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, service.ErrOrganizationNotFound
		}

		return nil, result.Error
	}

	return organization.toModel()
}

func (r *OrganizationRepository) CreateOrganization(ctx context.Context, operator domain.SystemAdminInterface, organizationName string) (*domain.OrganizationID, error) {
	_, span := tracer.Start(ctx, "organizationRepository.CreateOrganization")
	defer span.End()

	organization := organizationEntity{ //nolint:exhaustruct
		BaseModelEntity: BaseModelEntity{ //nolint:exhaustruct
			Version:   1,
			CreatedBy: operator.GetUserID().Int(),
			UpdatedBy: operator.GetUserID().Int(),
		},
		Name: organizationName,
	}

	if result := r.db.WithContext(ctx).Create(&organization); result.Error != nil {
		return nil, fmt.Errorf("create organization: %w", libgateway.ConvertDuplicatedError(result.Error, service.ErrOrganizationAlreadyExists))
	}

	if organization.ID == 0 {
		return nil, fmt.Errorf("organization.ID is 0")
	}

	organizationID, err := domain.NewOrganizationID(organization.ID)
	if err != nil {
		return nil, fmt.Errorf("new organization ID: %w", err)
	}

	return organizationID, nil
}
