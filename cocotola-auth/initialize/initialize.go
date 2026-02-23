package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	libhandler "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/handler"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/config"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/controller/handler"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/controller/middleware"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/usecase/auth"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/usecase/guest"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/usecase/password"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/usecase/profile"
)

// func newCallbackOnAddUser(cocotolaAuthCallbackClient service.CocotolaAuthCallbackClient, logger *slog.Logger) func(ctx context.Context, obj any) {
// 	return func(ctx context.Context, obj any) {
// 		param, ok := obj.(map[string]int)
// 		if !ok {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid object type: %T", obj))
// 			return
// 		}

// 		organizationIDInt, ok := param["organizationId"]
// 		if !ok {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid organizationId type: %T", param["organizationId"]))
// 			return
// 		}

// 		organizationID, err := domain.NewOrganizationID(organizationIDInt)
// 		if err != nil {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid organizationId: %v", err))
// 			return
// 		}

// 		userIDInt, ok := param["userId"]
// 		if !ok {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid userId type: %T", param["userId"]))
// 			return
// 		}

// 		userID, err := domain.NewUserID(userIDInt)
// 		if err != nil {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid userId: %v", err))
// 			return
// 		}

// 		logger.InfoContext(ctx, fmt.Sprintf("OnAddUser: organizationID=%d, userID=%d", organizationID.Int(), userID.Int()))
// 		if err := cocotolaAuthCallbackClient.OnAddUser(ctx, organizationID, userID); err != nil {
// 			logger.ErrorContext(ctx, fmt.Sprintf("OnAddUser: %v", err))
// 			return
// 		}
// 	}
// }

// func newCallbackOnAddUserSpace(cocotolaCoreCallbackClient service.CocotolaCoreCallbackClient, logger *slog.Logger) func(ctx context.Context, obj any) {
// 	return func(ctx context.Context, obj any) {
// 		param, ok := obj.(map[string]int)
// 		if !ok {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid object type: %T", obj))
// 			return
// 		}

// 		organizationIDInt, ok := param["organizationId"]
// 		if !ok {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid organizationId type: %T", param["organizationId"]))
// 			return
// 		}

// 		organizationID, err := domain.NewOrganizationID(organizationIDInt)
// 		if err != nil {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid organizationId: %v", err))
// 			return
// 		}

// 		userIDInt, ok := param["userId"]
// 		if !ok {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid userId type: %T", param["userId"]))
// 			return
// 		}

// 		userID, err := domain.NewUserID(userIDInt)
// 		if err != nil {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid userId: %v", err))
// 			return
// 		}

// 		spaceIDInt, ok := param["spaceId"]
// 		if !ok {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid spaceID type: %T", param["spaceID"]))
// 			return
// 		}

// 		spaceID, err := domain.NewSpaceID(spaceIDInt)
// 		if err != nil {
// 			logger.ErrorContext(ctx, fmt.Sprintf("invalid spaceID: %v", err))
// 			return
// 		}

// 		logger.InfoContext(ctx, fmt.Sprintf("OnAddUserSpace: organizationID=%d, userID=%d, spaceID:%d", organizationID.Int(), userID.Int(), spaceID.Int()))
// 		if err := cocotolaCoreCallbackClient.OnAddUserSpace(ctx, organizationID, userID, spaceID); err != nil {
// 			logger.ErrorContext(ctx, fmt.Sprintf("OnAddUser: %v", err))
// 			return
// 		}
// 	}
// }

func Initialize(ctx context.Context, systemToken domain.SystemToken, parent gin.IRouter, dbConn *libgateway.DBConnection, logConfig *libhandler.LogConfig, authConfig *config.AuthConfig) error {
	ctx, span := tracer.Start(ctx, "Initialize")
	defer span.End()

	if err := initApp(ctx, systemToken, parent, dbConn, logConfig, authConfig); err != nil {
		return fmt.Errorf("initApp: %w", err)
	}

	return nil
}

func initApp(_ context.Context, systemToken domain.SystemToken, parent gin.IRouter, dbc *libgateway.DBConnection, logConfig *libhandler.LogConfig, authConfig *config.AuthConfig) error {
	ctx := context.Background()
	// logger := slog.Default().With(slog.String(mbliblog.LoggerNameKey, domain.AppName+"initApp"))

	// cocotolaAuthCallbackClient := initCocotolaAuthCallbackClient(authConfig)
	// cocotolaCoreCallbackClient := initCocotolaCoreCallbackClient(authConfig.CoreAPIClient)

	// userEventHandler := mblibservice.ResourceEventHandlerFuncs{ //nolint:exhaustruct
	// 	AddFunc: newCallbackOnAddUser(cocotolaAuthCallbackClient, logger),
	// }
	// spaceEventHandler := mblibservice.ResourceEventHandlerFuncs{ //nolint:exhaustruct
	// 	AddFunc: newCallbackOnAddUserSpace(cocotolaCoreCallbackClient, logger),
	// }
	// resouceEventHandlers := map[domain.ResourceKey]mblibservice.ResourceEventHandler{
	// 	domain.ResourceUser:  userEventHandler,
	// 	domain.RecourceSpace: spaceEventHandler,
	// }
	// rff := func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
	// 	return gateway.NewRepositoryFactory(ctx, dbConn.Dialect, dbConn.DriverName, db, time.UTC)
	// }
	// rf, err := rff(ctx, dbConn.DB)
	// if err != nil {
	// 	return fmt.Errorf("rff: %w", err)
	// }

	orgRepo := gateway.NewOrganizationRepository(dbc)
	userRepo := gateway.NewUserRepository(dbc)
	spaceManager, err := gateway.NewSpaceManager(ctx, dbc)
	if err != nil {
		return fmt.Errorf("NewSpaceManager: %w", err)
	}
	// init auth token manager
	signingKey := []byte(authConfig.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager := gateway.NewAuthTokenManager(ctx, signingKey, signingMethod, time.Duration(authConfig.AccessTokenTTLMin)*time.Minute, time.Duration(authConfig.RefreshTokenTTLHour)*time.Hour)

	// init public and private router group functions
	// publicRouterGroupFuncs, err := ctrlgin.GetPublicRouterGroupFuncs(ctx, systemToken, authConfig, txManager, nonTxManager, authTokenManager)
	// if err != nil {
	// 	return fmt.Errorf("GetPublicRouterGroupFuncs: %w", err)
	// }
	// bearerTokenRouterGroupFuncs := ctrlgin.GetBearerTokenRouterGroupFuncs(ctx, systemToken, txManager, nonTxManager, authTokenManager, rf)
	// basicPrivateRouterGroupFuncs := controller.GetBasicPrivateRouterGroupFuncs(ctx, systemToken, cocotolaCoreCallbackClient)

	// api
	api := libhandler.InitAPIRouterGroup(ctx, parent, logConfig, domain.AppName)

	// v1
	v1 := api.Group("v1")

	// public router
	// libhandler.InitPublicAPIRouterGroup(ctx, v1, publicRouterGroupFuncs)
	{
		passwordUsecase := newPasswordUsecase(systemToken, orgRepo, userRepo, authTokenManager)
		handler.InitPasswordRouter(passwordUsecase, v1)
	}
	{
		guestUsecase := newGuestUsecase(systemToken, orgRepo, userRepo, authTokenManager)
		handler.InitGuestRouter(guestUsecase, v1)
	}
	authUsecase := newVerifyAccessTokenQuery(systemToken, orgRepo, authTokenManager, userRepo)
	bearerTokenAuthMiddleware := middleware.NewBearerTokenAuthMiddleware(authUsecase)
	{
		profileUsecase := newProfileUsecase(orgRepo, userRepo, spaceManager)
		handler.InitProfileRouter(profileUsecase, v1, bearerTokenAuthMiddleware)
	}

	// private router
	// libgin.InitPrivateAPIRouterGroup(ctx, v1, bearerTokenAuthMiddleware, bearerTokenRouterGroupFuncs)

	// libcontroller.InitPrivateAPIRouterGroup(ctx, v1, basicAuthMiddleware, basicPrivateRouterGroupFuncs)

	return nil
}

type VerifyAccessTokenQueryGateway struct {
	*gateway.OrganizationRepository
	*gateway.UserRepository
	*gateway.AuthTokenManager
}

func newVerifyAccessTokenQuery(systemToken domain.SystemToken, orgRepo *gateway.OrganizationRepository, authTokenManager *gateway.AuthTokenManager, userRepo *gateway.UserRepository) *auth.Usecase {
	verifyAccessTokenQueryRepo := &VerifyAccessTokenQueryGateway{
		OrganizationRepository: orgRepo,
		UserRepository:         userRepo,
		AuthTokenManager:       authTokenManager,
	}
	authUsecase := auth.NewUsecase(systemToken, verifyAccessTokenQueryRepo)
	return authUsecase
}

type GuestGateway struct {
	*gateway.UserRepository
	*gateway.OrganizationRepository
	*gateway.AuthTokenManager
}

func newGuestUsecase(systemToken domain.SystemToken, orgRepo *gateway.OrganizationRepository, userRepo *gateway.UserRepository, authTokenManager *gateway.AuthTokenManager) *guest.Usecase {
	guestRepo := &GuestGateway{
		UserRepository:         userRepo,
		OrganizationRepository: orgRepo,
		AuthTokenManager:       authTokenManager,
	}
	guestUsecase := guest.NewUsecase(systemToken, guestRepo)
	return guestUsecase
}

type PasswordRepo struct {
	*gateway.UserRepository
	*gateway.OrganizationRepository
	*gateway.AuthTokenManager
}

func newPasswordUsecase(systemToken domain.SystemToken, orgRepo *gateway.OrganizationRepository, userRepo *gateway.UserRepository, authTokenManager *gateway.AuthTokenManager) *password.Usecase {
	passwordRepo := &PasswordRepo{
		UserRepository:         userRepo,
		OrganizationRepository: orgRepo,
		AuthTokenManager:       authTokenManager,
	}
	passwordUsecase := password.NewPassword(systemToken, passwordRepo)
	return passwordUsecase
}

type ProfileRepo struct {
	*gateway.OrganizationRepository
	*gateway.UserRepository
	*gateway.SpaceManager
}

//	func (r *ProfileRepo) GetOrganization(ctx context.Context, operator domain.UserInterface) (*domain.Organization, error) {
//		return r.organizationRepo.GetOrganization(ctx, operator)
//	}
//
// func
func newProfileUsecase(orgRepo *gateway.OrganizationRepository, userRepo *gateway.UserRepository, spaceManager *gateway.SpaceManager) *profile.Usecase {
	profileRepo := &ProfileRepo{
		OrganizationRepository: orgRepo,
		UserRepository:         userRepo,
		SpaceManager:           spaceManager,
	}
	profileUsecase := profile.NewUsecase(profileRepo)
	return profileUsecase
}

// func initProfileRouter(_ context.Context, _ domain.SystemToken, _, mbNonTxManager service.TransactionManager, _ service.AuthTokenManager, rf service.RepositoryFactory) error {
// 	orgRepo := rf.NewOrganizationRepository(context.Background())
// 	userRepo := rf.NewUserRepository(context.Background())
// 	spaceManager, err := rf.NewSpaceManager(context.Background())
// 	if err != nil {
// 		return fmt.Errorf("NewSpaceManager: %w", err)
// 	}
// 	profileUsecase := newProfileUsecase(mbNonTxManager, orgRepo, userRepo, spaceManager)
// 	return nil
// }

// func initCocotolaAuthCallbackClient(authConfig *config.AuthConfig) service.CocotolaAuthCallbackClient {
// 	httpClient := http.Client{ //nolint:exhaustruct
// 		Timeout:   time.Duration(authConfig.AuthAPIClient.TimeoutSec) * time.Second,
// 		Transport: otelhttp.NewTransport(http.DefaultTransport),
// 	}
// 	authAPIEndpoint, err := url.Parse(authConfig.AuthAPIClient.Endpoint)
// 	if err != nil {
// 		libdomain.CheckError(err)
// 	}

// 	cocotolaAuthCallbackClient := gateway.NewCocotolaAuthCallbackClient(&httpClient, authAPIEndpoint, authConfig.AuthAPIClient.Username, authConfig.AuthAPIClient.Password)

// 	return cocotolaAuthCallbackClient
// }

// func initCocotolaCoreCallbackClient(coreAPIClientConfig *config.CoreAPIClientConfig) service.CocotolaCoreCallbackClient {
// 	httpClient := http.Client{ //nolint:exhaustruct
// 		Timeout:   time.Duration(coreAPIClientConfig.TimeoutSec) * time.Second,
// 		Transport: otelhttp.NewTransport(http.DefaultTransport),
// 	}
// 	coreAPIEndpoint, err := url.Parse(coreAPIClientConfig.Endpoint)
// 	if err != nil {
// 		libdomain.CheckError(err)
// 	}

// 	cocotolaCoreCallbackClient := gateway.NewCocotolaCoreCallbackClient(&httpClient, coreAPIEndpoint, coreAPIClientConfig.Username, coreAPIClientConfig.Password)

// 	return cocotolaCoreCallbackClient
// }

// func initMBTransactionManager(db *gorm.DB, rff func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)) service.TransactionManager {
// 	txManager, err := mblibgateway.NewTransactionManagerT(db, rff)
// 	if err != nil {
// 		libdomain.CheckError(err)
// 	}
// 	return txManager
// }

// func initMBNonTransactionManager(rf service.RepositoryFactory) service.TransactionManager {
// 	nonTxManager, err := mblibgateway.NewNonTransactionManagerT(rf)
// 	if err != nil {
// 		libdomain.CheckError(err)
// 	}
// 	return nonTxManager
// }

// func initTransactionManager(db *gorm.DB, rff func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)) (service.TransactionManager, error) {
// 	txManager, err := libgateway.NewTransactionManagerT(db, rff)
// 	if err != nil {
// 		return nil, fmt.Errorf("NewTransactionManagerT: %w", err)
// 	}
// 	return txManager, nil
// }

// func initNonTransactionManager(rf service.RepositoryFactory) (service.TransactionManager, error) {
// 	nonTxManager, err := libgateway.NewNonTransactionManagerT(rf)
// 	if err != nil {
// 		return nil, fmt.Errorf("NewNonTransactionManagerT: %w", err)
// 	}
// 	return nonTxManager, nil
// }
