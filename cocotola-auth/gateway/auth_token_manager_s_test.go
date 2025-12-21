//go:build small

package gateway_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/gateway"
	"github.com/mocoarow/cocotola-1.25/cocotola-auth/service"
)

func Test_authTokenManager_CreateTokenSet(t *testing.T) {
	t.Parallel()
	organizationID := organizationID(t, 123)
	userID := userID(t, 456)
	type fields struct {
		SigningKey     []byte
		SigningMethod  jwt.SigningMethod
		TokenTimeout   time.Duration
		RefreshTimeout time.Duration
	}
	type args struct {
		user             domain.UserInterface
		organizationID   *domain.OrganizationID
		organizationName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				SigningKey:    []byte("&32KC^L;m'BuH+'ATNhv[qWM:3)2Lt2m"),
				SigningMethod: jwt.SigningMethodHS256,
			},
			args: args{
				user: &user{
					userID:         userID,
					organizationID: organizationID,
					loginID:        "LOGIN_ID",
					username:       "USERNAME",
				},
				organizationID:   organizationID,
				organizationName: "ORG_NAME",
			},
			wantErr: false,
		},
		{
			name: "signingKey is not set",
			fields: fields{
				SigningKey:    []byte(""),
				SigningMethod: jwt.SigningMethodHS256,
			},
			args: args{
				user: &user{
					userID:         userID,
					organizationID: organizationID,
					loginID:        "LOGIN_ID",
					username:       "USERNAME",
				},
				organizationID:   organizationID,
				organizationName: "ORG_NAME",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		m := gateway.NewAuthTokenManager(ctx, tt.fields.SigningKey, tt.fields.SigningMethod, tt.fields.TokenTimeout, tt.fields.RefreshTimeout)

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := m.CreateTokenSet(ctx, tt.args.user, tt.args.organizationID, tt.args.organizationName)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("authTokenManager.CreateTokenSet() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				return
			}
			assert.NotEmpty(t, got.AccessToken)
			assert.NotEmpty(t, got.RefreshToken)
		})
	}
}

func TestAuthTokenManager_GetUserInfo(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	organizationID := organizationID(t, 123)
	userID := userID(t, 456)
	user := &user{
		userID:         userID,
		organizationID: organizationID,
		loginID:        "LOGIN_ID",
		username:       "USERNAME",
	}
	organizationName := "ORG_NAME"

	type fields struct {
		SigningKey     []byte
		SigningMethod  jwt.SigningMethod
		TokenTimeout   time.Duration
		RefreshTimeout time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    *service.UserInfo
		wantErr error
	}{
		{
			name: "valid",
			fields: fields{
				SigningKey:    []byte("&32KC^L;m'BuH+'ATNhv[qWM:3)2Lt2m"),
				SigningMethod: jwt.SigningMethodHS256,
				TokenTimeout:  time.Second,
			},
			want: &service.UserInfo{
				// UserID:        456,
				LoginID:          "LOGIN_ID",
				Username:         "USERNAME",
				OrganizationID:   123,
				OrganizationName: "ORG_NAME",
			},
			wantErr: nil,
		},
		{
			name: "expired",
			fields: fields{
				SigningKey:    []byte("&32KC^L;m'BuH+'ATNhv[qWM:3)2Lt2m"),
				SigningMethod: jwt.SigningMethodHS256,
				TokenTimeout:  -1 * time.Second,
			},
			wantErr: service.ErrUnauthenticated,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := gateway.NewAuthTokenManager(ctx, tt.fields.SigningKey, tt.fields.SigningMethod, tt.fields.TokenTimeout, tt.fields.RefreshTimeout)

			tokenSet, err := m.CreateTokenSet(ctx, user, organizationID, organizationName)
			require.NoError(t, err)
			got, err := m.GetUserInfo(ctx, tokenSet.AccessToken)
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthTokenManager.GetUserInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
