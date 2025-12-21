package usecase

import (
	"context"
	"fmt"
	"log/slog"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type ProfileUsecase struct {
	nonTxManager service.TransactionManager
	logger       *slog.Logger
}

func NewProfileUsecase(nonTxManager service.TransactionManager) *ProfileUsecase {
	return &ProfileUsecase{
		nonTxManager: nonTxManager,
		logger:       slog.Default().With(slog.String(libdomain.LoggerNameKey, "ProfileUsecase")),
	}
}

func (u *ProfileUsecase) GetMyProfile(ctx context.Context, operator domain.UserInterface) (*domain.ProfileModel, error) {
	command := NewGetMyProfileQuery(u.nonTxManager)
	profile, err := command.Execute(ctx, operator)
	if err != nil {
		return nil, fmt.Errorf("GetMyProfileQuery.Execute: %w", err)
	}
	return profile, err
}
