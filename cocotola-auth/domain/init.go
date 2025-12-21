package domain

const AppName = "cocotola-auth"

type SystemAdminInterface interface {
	GetUserID() *UserID
	IsSystemAdmin() bool
}

type UserInterface interface {
	GetUserID() *UserID
	GetOrganizationID() *OrganizationID
	GetLoginID() string
	GetUsername() string
}

type OwnerInterface interface {
	UserInterface
	IsOwner() bool
}

type SystemOwnerInterface interface {
	OwnerInterface
	IsSystemOwner() bool
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
