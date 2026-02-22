package helper

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mocoarow/cocotola-1.25/cocotola-auth/domain"

	libcontroller "github.com/mocoarow/cocotola-1.25/cocotola-lib/controller"
)

type operator struct {
	userID         *domain.UserID
	loginID        string
	username       string
	organizationID *domain.OrganizationID
}

func (o *operator) GetUserID() *domain.UserID {
	return o.userID
}

func (o *operator) GetLoginID() string {
	return o.loginID
}

func (o *operator) GetUsername() string {
	return o.username
}

func (o *operator) GetOrganizationID() *domain.OrganizationID {
	return o.organizationID
}

func HandleUserFunction(c *gin.Context, fn func(ctx context.Context, operator domain.UserInterface) error, errorHandle func(ctx context.Context, c *gin.Context, err error) bool) {
	ctx := c.Request.Context()
	organizationIDInt := c.GetInt("OrganizationID")
	if organizationIDInt == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})

		return
	}

	organizationID, err := domain.NewOrganizationID(organizationIDInt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})

		return
	}

	userID := c.GetInt("AuthorizedUser")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}

	loginID := c.GetString("LoginID")
	if loginID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}

	username := c.GetString("Username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}

	operatorID, err := domain.NewUserID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}

	// logger.InfoContext(ctx, "", slog.Int("organization_id", organizationID.Int()), slog.Int("operator_id", operatorID.Int()))

	operator := &operator{
		userID:         operatorID,
		loginID:        loginID,
		username:       username,
		organizationID: organizationID,
	}

	if newCtx, err := libcontroller.AddBaggageMembers(ctx, map[string]string{
		"operator_id":     strconv.Itoa(operatorID.Int()),
		"organization_id": strconv.Itoa(organizationID.Int()),
	}); err == nil {
		ctx = newCtx
		// Add baggage members as span attributes
		libcontroller.AddBaggageToCurrentSpan(ctx)
	}

	if err := fn(ctx, operator); err != nil {
		if handled := errorHandle(ctx, c, err); !handled {
			c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
		}
	}
}
