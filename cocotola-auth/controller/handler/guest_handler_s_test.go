//go:build small

package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	libhandler "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller/handler"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/controller/handler"
)

func initGuestRouter(ctx context.Context, t *testing.T, guest handler.GuestUsecase) *gin.Engine {
	t.Helper()

	router := libhandler.InitRootRouterGroup(ctx, &config, "cocotola-auth-test")
	api := router.Group("api")
	v1 := api.Group("v1")

	handler.InitGuestRouter(guest, v1)

	return router
}

func TestGuestHandler_Authenticate_shouldReturn400_whenInvalidRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		requestBody string
	}{
		{
			name:        "request body is empty",
			requestBody: "",
		},
		{
			name:        "organizationName is empty",
			requestBody: `{"organizationName":""}`,
		},
		{
			name:        "organizationName exceeds maxLength of 20",
			requestBody: `{"organizationName":"this-is-a-very-long-organization-name"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// given
			guestUserUsecase := NewMockGuestUsecase(t)
			r := initGuestRouter(ctx, t, guestUserUsecase)
			w := httptest.NewRecorder()

			// when
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/guest/authenticate", bytes.NewBufferString(tt.requestBody))
			require.NoError(t, err)
			r.ServeHTTP(w, req)
			respBytes := readBytes(t, w.Body)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code, "status code should be 400")

			jsonObj := parseJSON(t, respBytes)

			messageExpr := parseExpr(t, "$.message")
			message := messageExpr.Get(jsonObj)
			assert.Equal(t, "Bad Request", message[0], "message should be 'Bad Request'")
		})
	}
}
