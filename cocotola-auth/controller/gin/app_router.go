package gin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	libgin "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/config"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/controller/gin/middleware"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/usecase"
)

func NewInitTestRouterFunc() libgin.InitRouterGroupFunc {
	return func(parentRouterGroup gin.IRouter, middleware ...gin.HandlerFunc) {
		test := parentRouterGroup.Group("test")
		for _, m := range middleware {
			test.Use(m)
		}
		test.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}
}

// func NewAuthTokenManager(ctx context.Context, authConfig *config.AuthConfig) (service.AuthTokenManager, error) {
// 	signingKey := []byte(authConfig.SigningKey)
// 	signingMethod := jwt.SigningMethodHS256
// 	fireabseAuthClient, err := gateway.NewFirebaseClient(ctx, authConfig.GoogleProjectID)
// 	if err != nil {
// 		return nil, mbliberrors.Errorf("NewFirebaseClient: %w", err)
// 	}
// 	authTokenManager := gateway.NewAuthTokenManager(ctx, fireabseAuthClient, signingKey, signingMethod, time.Duration(authConfig.AccessTokenTTLMin)*time.Minute, time.Duration(authConfig.RefreshTokenTTLHour)*time.Hour)

// 	return authTokenManager, nil
// }

func GetPublicRouterGroupFuncs(_ context.Context, systemToken domain.SystemToken, _ *config.AuthConfig, txManager, nonTxManager service.TransactionManager, authTokenManager service.AuthTokenManager) ([]libgin.InitRouterGroupFunc, error) {
	// // - google
	// httpClient := http.Client{ //nolint:exhaustruct
	// 	Timeout:   time.Duration(authConfig.GoogleAPITimeoutSec) * time.Second,
	// 	Transport: otelhttp.NewTransport(http.DefaultTransport),
	// }

	// googleAuthClient := gateway.NewGoogleAuthClient(&httpClient, authConfig.GoogleClientID, authConfig.GoogleClientSecret, authConfig.GoogleCallbackURL)
	// googleUserUsecase := usecase.NewGoogleUser(systemToken, mbTxManager, mbNonTxManager, txManager, nonTxManager, authTokenManager, googleAuthClient)
	// // - authentication
	// authenticationUsecase := usecase.NewAuthentication(systemToken, mbTxManager, authTokenManager)
	// // &systemOwnerByOrganizationName{})
	// - password
	passwordUsecase := usecase.NewPassword(systemToken, txManager, nonTxManager, authTokenManager)
	// - guest
	guestUsecase := usecase.NewGuest(systemToken, txManager, nonTxManager, authTokenManager)

	// public router
	return []libgin.InitRouterGroupFunc{
		NewInitTestRouterFunc(),
		// public.NewInitAuthRouterFunc(authenticationUsecase),
		// public.NewInitGoogleRouterFunc(googleUserUsecase),
		NewInitPasswordRouterFunc(passwordUsecase),
		NewInitGuestRouterFunc(guestUsecase),
	}, nil
}

// func GetBasicPrivateRouterGroupFuncs(_ context.Context, systemToken libdomain.SystemToken, mbTxManager, mbNonTxManager mbuserservice.TransactionManager, cocotolaCoreCallbackClient service.CocotolaCoreCallbackClient) []libcontroller.InitRouterGroupFunc {
// 	// - rbac
// 	rbacUsecase := usecase.NewRBACUsecase(systemToken, mbTxManager, mbNonTxManager)
// 	// - callback
// 	callbackUsecase := usecase.NewCallbackUsecase(systemToken, mbTxManager, mbNonTxManager, cocotolaCoreCallbackClient)

// 	// private router
// 	return []libcontroller.InitRouterGroupFunc{
// 		private.NewInitRBACRouterFunc(rbacUsecase),
// 		private.NewInitCallbackRouterFunc(callbackUsecase),
// 	}
// }

func GetBearerTokenRouterGroupFuncs(_ context.Context, _ domain.SystemToken, _, mbNonTxManager service.TransactionManager, _ service.AuthTokenManager, _ service.RepositoryFactory) []libgin.InitRouterGroupFunc {
	// - user
	// userUsecase := usecase.NewUserUsecase(systemToken, mbTxManager, mbNonTxManager, authTokenManager)
	// spaceUsecase := usecase.NewSpaceUsecase(mbrf)
	profileUsecase := usecase.NewProfileUsecase(mbNonTxManager)
	return []libgin.InitRouterGroupFunc{
		// public.NewInitUserRouterFunc(userUsecase),
		// private.NewInitSpaceRouterFunc(spaceUsecase),
		NewInitProfileRouterFunc(profileUsecase),
		// NewInitRBACRouterFunc(rbacUsecase),
	}
}

func InitBearerTokenAuthMiddleware(systemToken domain.SystemToken, authTokenManager service.AuthTokenManager, mbNonTxManager service.TransactionManager, mbrf service.RepositoryFactory) (gin.HandlerFunc, error) {
	return middleware.NewBearerTokenAuthMiddleware(systemToken, authTokenManager, mbNonTxManager, mbrf), nil
}
