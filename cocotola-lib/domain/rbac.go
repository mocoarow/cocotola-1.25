package domain

type RBACSubject interface {
	Subject() string
}

type RBACUserInterface interface {
	RBACSubject
}

type RBACUser struct {
	value string
}

func NewRBACUser(value string) *RBACUser {
	return &RBACUser{value: value}
}

func (r *RBACUser) Subject() string {
	return r.value
}

type RBACRoleInterface interface {
	RBACSubject
	Role() string
}

type RBACRole struct {
	value string
}

func NewRBACRole(value string) *RBACRole {
	return &RBACRole{value: value}
}

func (r *RBACRole) Subject() string {
	return r.value
}
func (r *RBACRole) Role() string {
	return r.value
}

type RBACDomainInterface interface {
	Domain() string
}

type RBACDomain struct {
	value string
}

func NewRBACDomain(value string) *RBACDomain {
	return &RBACDomain{value: value}
}

func (r *RBACDomain) Domain() string {
	return r.value
}

type RBACObjectInterface interface {
	Object() string
}

type RBACObject struct {
	value string
}

func NewRBACObject(value string) *RBACObject {
	return &RBACObject{value: value}
}

func (r *RBACObject) Object() string {
	return r.value
}

type RBACAction interface {
	Action() string
}

type rbacAction struct {
	value string
}

func NewRBACAction(value string) RBACAction {
	return &rbacAction{value: value}
}

func (r *rbacAction) Action() string {
	return r.value
}

type RBACEffect interface {
	Effect() string
}

type rbacEffect struct {
	value string
}

func NewRBACEffect(value string) RBACEffect {
	return &rbacEffect{value: value}
}

func (r *rbacEffect) Effect() string {
	return r.value
}

type ActionObjectEffect struct {
	Action RBACAction
	Object *RBACObject
	Effect RBACEffect
}
