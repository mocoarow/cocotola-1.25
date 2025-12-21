package gin

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	libgin "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/openapi"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type PasswordUsecase interface {
	Authenticate(ctx context.Context, loginID, password, organizationName string) (*service.AuthTokenSet, error)
}

type PasswordHandler struct {
	passwordUsecase PasswordUsecase
	logger          *slog.Logger
}

func NewPasswordHandler(passwordUsecase PasswordUsecase) *PasswordHandler {
	return &PasswordHandler{
		passwordUsecase: passwordUsecase,
		logger:          slog.Default().With(slog.String(libdomain.LoggerNameKey, domain.AppName+"-PasswordAuthHandler")),
	}
}

func (h *PasswordHandler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()

	var passwordAuthRequest openapi.PasswordAuthRequest
	if err := c.ShouldBindJSON(&passwordAuthRequest); err != nil {
		h.logger.InfoContext(ctx, fmt.Sprintf("invalid parameter: %+v", err))
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	authResult, err := h.passwordUsecase.Authenticate(ctx, passwordAuthRequest.LoginID, passwordAuthRequest.Password, passwordAuthRequest.OrganizationName)
	if err != nil {
		if errors.Is(err, service.ErrSystemOwnerNotFound) {
			h.logger.InfoContext(ctx, fmt.Sprintf("system owner not found: %+v", err))
			c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		if errors.Is(err, service.ErrUnauthenticated) {
			h.logger.InfoContext(ctx, fmt.Sprintf("invalid parameter: %+v", err))
			c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		h.logger.ErrorContext(ctx, fmt.Sprintf("passwordUsecase.Authenticate: %+v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, openapi.AuthResponse{
		AccessToken:  authResult.AccessToken,
		RefreshToken: authResult.RefreshToken,
	})
}

func NewInitPasswordRouterFunc(password PasswordUsecase) libgin.InitRouterGroupFunc {
	return func(parentRouterGroup gin.IRouter, middleware ...gin.HandlerFunc) {
		auth := parentRouterGroup.Group("password")
		for _, m := range middleware {
			auth.Use(m)
		}

		passwordHandler := NewPasswordHandler(password)
		auth.POST("authenticate", passwordHandler.Authorize)
	}
}
