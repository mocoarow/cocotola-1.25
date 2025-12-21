package gin

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/controller/gin/helper"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/openapi"

	libgin "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
)

type ProfileQueryUsecase interface {
	GetMyProfile(ctx context.Context, operator domain.UserInterface) (*domain.ProfileModel, error)
}

type ProfileHandler struct {
	profileQueryUsecase ProfileQueryUsecase
	logger              *slog.Logger
}

func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	helper.HandleUserFunction(c, func(ctx context.Context, operator domain.UserInterface) error {
		result, err := h.profileQueryUsecase.GetMyProfile(ctx, operator)
		if err != nil {
			return fmt.Errorf("GetMyProfile: %w", err)
		}

		apiResp := openapi.GetMyProfileResponse{
			LoginID:          result.LoginID,
			Username:         result.Username,
			OrganizationID:   int32(result.OrganizationID.Int()),
			OrganizationName: result.OrganizationName,
			PrivateSpaceID:   int32(result.PrivateSpaceID.Int()),
		}
		c.JSON(http.StatusOK, apiResp)

		return nil
	}, h.errorHandle)
}

func NewProfileHandler(profileQueryUsecase ProfileQueryUsecase) *ProfileHandler {
	return &ProfileHandler{
		profileQueryUsecase: profileQueryUsecase,
		logger:              slog.Default().With(slog.String(libdomain.LoggerNameKey, "ProfileHandler")),
	}
}

func (h *ProfileHandler) errorHandle(ctx context.Context, _ *gin.Context, err error) bool {
	h.logger.ErrorContext(ctx, fmt.Sprintf("ProfileHandler. error: %+v", err))

	return false
}

func NewInitProfileRouterFunc(profileQueryUsecase ProfileQueryUsecase) libgin.InitRouterGroupFunc {
	return func(parentRouterGroup gin.IRouter, middleware ...gin.HandlerFunc) {
		profile := parentRouterGroup.Group("profile")
		for _, m := range middleware {
			profile.Use(m)
		}
		profileHandler := NewProfileHandler(profileQueryUsecase)
		profile.GET("me", profileHandler.GetMyProfile)
	}
}
