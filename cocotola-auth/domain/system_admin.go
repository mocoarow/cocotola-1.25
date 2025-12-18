package domain

type SystemAdmin struct {
	UserID *UserID
}

func NewSystemAdmin(_ SystemToken) *SystemAdmin {
	return &SystemAdmin{
		UserID: SystemAdminID,
	}
}

func (m *SystemAdmin) IsSystemAdmin() bool {
	return true
}
func (m *SystemAdmin) GetUserID() *UserID {
	return m.UserID
}
