package domain

type SystemAdminInterface interface {
	GetUserID() *UserID
	IsSystemAdmin() bool
}

type UserInterface interface {
	GetUserID() *UserID
	GetOrganizationID() *OrganizationID
}

var (
	SystemAdminID *UserID
)

func init() {
	systemAdminID := 1
	systemAdminIDTmp, err := NewUserID(systemAdminID)
	if err != nil {
		panic(err)
	}
	SystemAdminID = systemAdminIDTmp
}
