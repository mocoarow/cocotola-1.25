package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/config"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
	libcontroller "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/gin"
)

func NewInitTestRouterFunc() libcontroller.InitRouterGroupFunc {
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

func GetPublicRouterGroupFuncs(_ context.Context, _ domain.SystemToken, _ *config.AuthConfig, _, _ service.TransactionManager) ([]libcontroller.InitRouterGroupFunc, error) {
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
	// // - password
	// passwordUsecase := usecase.NewPassword(systemToken, mbTxManager, mbNonTxManager, authTokenManager)
	// // - guest
	// guestUsecase := usecase.NewGuest(systemToken, mbTxManager, mbNonTxManager, authTokenManager)

	// public router
	return []libcontroller.InitRouterGroupFunc{
		NewInitTestRouterFunc(),
		// public.NewInitAuthRouterFunc(authenticationUsecase),
		// public.NewInitGoogleRouterFunc(googleUserUsecase),
		// public.NewInitPasswordRouterFunc(passwordUsecase),
		// public.NewInitGuestRouterFunc(guestUsecase),
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

// func GetBearerTokenRouterGroupFuncs(_ context.Context, systemToken libdomain.SystemToken, mbTxManager, mbNonTxManager mbuserservice.TransactionManager, authTokenManager service.AuthTokenManager, mbrf mbuserservice.RepositoryFactory) []libcontroller.InitRouterGroupFunc {
// 	// - user
// 	userUsecase := usecase.NewUserUsecase(systemToken, mbTxManager, mbNonTxManager, authTokenManager)
// 	spaceUsecase := usecase.NewSpaceUsecase(mbrf)
// 	profileUsecase := usecase.NewProfileUsecase(mbNonTxManager)
// 	return []libcontroller.InitRouterGroupFunc{
// 		public.NewInitUserRouterFunc(userUsecase),
// 		private.NewInitSpaceRouterFunc(spaceUsecase),
// 		private.NewInitProfileRouterFunc(profileUsecase),
// 		// NewInitRBACRouterFunc(rbacUsecase),
// 	}
// }

// func InitBearerTokenAuthMiddleware(systemToken libdomain.SystemToken, authTokenManager service.AuthTokenManager, mbNonTxManager mbuserservice.TransactionManager, mbrf mbuserservice.RepositoryFactory) (gin.HandlerFunc, error) {
// 	return middleware.NewAuthMiddleware(systemToken, authTokenManager, mbNonTxManager, mbrf), nil
// }
