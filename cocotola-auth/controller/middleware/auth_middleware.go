package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
)

type BearerTokenAuthGateway interface {
	VerifyAccessToken(ctx context.Context, accessToken string) (*domain.User, error)
}

func NewBearerTokenAuthMiddleware(bearerTokenAuthGateway BearerTokenAuthGateway) gin.HandlerFunc {
	logger := slog.Default().With(slog.String(libdomain.LoggerNameKey, domain.AppName+"-BearerTokenAuthMiddleware"))

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		ctx, span := tracer.Start(ctx, "AuthMiddleware")
		defer span.End()

		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			logger.InfoContext(ctx, "invalid header. Bearer not found")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		bearerToken := authorization[len("Bearer "):]
		user, err := bearerTokenAuthGateway.VerifyAccessToken(ctx, bearerToken)
		if err != nil {
			logger.WarnContext(ctx, fmt.Sprintf("verifyAccessToken: %v", err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("AuthorizedUser", user.UserID.Int())
		c.Set("LoginID", user.LoginID)
		c.Set("Username", user.Username)
		c.Set("OrganizationID", user.OrganizationID.Int())
		c.Next()
	}
}
