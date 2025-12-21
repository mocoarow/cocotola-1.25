package gateway

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

type UserClaims struct {
	LoginID string `json:"loginId"`
	// TODO: Check if UserID is needed in the token
	UserID           int    `json:"userId"`
	Username         string `json:"username"`
	OrganizationID   int    `json:"organizationId"`
	OrganizationName string `json:"organizationName"`
	// Role             string `json:"role"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}

type AuthTokenManager struct {
	SigningKey     []byte
	SigningMethod  jwt.SigningMethod
	TokenTimeout   time.Duration
	RefreshTimeout time.Duration
	logger         *slog.Logger
}

var _ service.AuthTokenManager = (*AuthTokenManager)(nil)

func NewAuthTokenManager(_ context.Context, signingKey []byte, signingMethod jwt.SigningMethod, tokenTimeout, refreshTimeout time.Duration) *AuthTokenManager {
	return &AuthTokenManager{
		SigningKey:     signingKey,
		SigningMethod:  signingMethod,
		TokenTimeout:   tokenTimeout,
		RefreshTimeout: refreshTimeout,
		logger:         slog.Default().With(slog.String(libdomain.LoggerNameKey, domain.AppName+"-AuthTokenManager")),
	}
}

func (m *AuthTokenManager) CreateTokenSet(ctx context.Context, user domain.UserInterface, organizationID *domain.OrganizationID, organizationName string) (*service.AuthTokenSet, error) {
	if user == nil {
		return nil, errors.New("user is nil")
	}
	accessToken, err := m.createJWT(ctx, user, organizationID, organizationName, m.TokenTimeout, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.createJWT(ctx, user, organizationID, organizationName, m.RefreshTimeout, "refresh")
	if err != nil {
		return nil, err
	}

	return &service.AuthTokenSet{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *AuthTokenManager) createJWT(ctx context.Context, user domain.UserInterface, organizationID *domain.OrganizationID, organizationName string, duration time.Duration, tokenType string) (string, error) {
	if len(m.SigningKey) == 0 {
		return "", fmt.Errorf("m.SigningKey is not set")
	}

	now := time.Now()
	claims := UserClaims{ //nolint:exhaustruct
		LoginID:          user.GetLoginID(),
		Username:         user.GetUsername(),
		OrganizationID:   organizationID.Int(),
		OrganizationName: organizationName,
		TokenType:        tokenType,
		RegisteredClaims: jwt.RegisteredClaims{ //nolint:exhaustruct
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}

	m.logger.DebugContext(ctx, fmt.Sprintf("claims: %+v", claims))

	token := jwt.NewWithClaims(m.SigningMethod, claims)
	signed, err := token.SignedString(m.SigningKey)
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %w", err)
	}

	return signed, nil
}

func (m *AuthTokenManager) GetUserInfo(ctx context.Context, tokenString string) (*service.UserInfo, error) {
	currentClaims, err := m.parseToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("parseToken(%s): %w", err.Error(), service.ErrUnauthenticated)
	}

	return &service.UserInfo{
		// UserID:        currentClaims.UserID,
		LoginID:          currentClaims.LoginID,
		Username:         currentClaims.Username,
		OrganizationID:   currentClaims.OrganizationID,
		OrganizationName: currentClaims.OrganizationName,
	}, nil
}

func (m *AuthTokenManager) parseToken(ctx context.Context, tokenString string) (*UserClaims, error) {
	keyFunc := func(_ *jwt.Token) (interface{}, error) {
		return m.SigningKey, nil
	}

	currentToken, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, keyFunc) //nolint:exhaustruct
	if err != nil {
		m.logger.InfoContext(ctx, fmt.Sprintf("%v", err))
		return nil, fmt.Errorf("ParseWithClaims: %w", err)
	}
	if !currentToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	currentClaims, ok := currentToken.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	v := jwt.NewValidator()
	if err := v.Validate(currentClaims); err != nil {
		return nil, fmt.Errorf("validate claims: %w", err)
	}

	return currentClaims, nil
}
