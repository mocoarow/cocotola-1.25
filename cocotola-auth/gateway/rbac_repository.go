package gateway

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"

	libdomain "github.com/mocoarow/cocotola-1.25/cocotola-lib/domain"
	libgateway "github.com/mocoarow/cocotola-1.25/cocotola-lib/gateway"
)

const rbacConf = `
[request_definition]
r = sub, obj, act, dom

[policy_definition]
p = sub, obj, act, eft, dom

[role_definition]
g = _, _, _
g2 = _, _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub, r.dom) && (keyMatch(r.obj, p.obj) || g2(r.obj, p.obj, r.dom)) && r.act == p.act
`

type RBACRepository struct {
	dbc      *libgateway.DBConnection
	conf     string
	enforcer *casbin.Enforcer
}

// var _ service.RBACRepository = (*RBACRepository)(nil)

func NewRBACRepository(_ context.Context, dbc *libgateway.DBConnection) (*RBACRepository, error) {
	if dbc == nil {
		panic(errors.New("dbc is nil"))
	}

	gormadapter.TurnOffAutoMigrate(dbc.DB)

	a, err := gormadapter.NewAdapterByDB(dbc.DB)
	if err != nil {
		return nil, fmt.Errorf("gormadapter.NewAdapterByDB: %w", err)
	}

	m, err := model.NewModelFromString(rbacConf)
	if err != nil {
		return nil, fmt.Errorf("model.NewModelFromString: %w", err)
	}

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, fmt.Errorf("casbin.NewEnforcer: %w", err)
	}
	e.EnableAutoSave(true)

	return &RBACRepository{
		dbc:      dbc,
		conf:     rbacConf,
		enforcer: e,
	}, nil
}

func (r *RBACRepository) initEnforcer(_ context.Context) *casbin.Enforcer {
	return r.enforcer
}

func (r *RBACRepository) CreatePolicy(ctx context.Context, domain libdomain.RBACDomainInterface, subject libdomain.RBACSubjectInterface, action libdomain.RBACActionInterface, object libdomain.RBACObjectInterface, effect libdomain.RBACEffectInterface) error {
	e := r.initEnforcer(ctx)

	if _, err := e.AddNamedPolicy("p", subject.Subject(), object.Object(), action.Action(), effect.Effect(), domain.Domain()); err != nil {
		return fmt.Errorf("e.AddNamedPolicy: %w", err)
	}

	return nil
}

func (r *RBACRepository) DeletePolicy(ctx context.Context, domain libdomain.RBACDomainInterface, subject libdomain.RBACSubjectInterface, action libdomain.RBACActionInterface, object libdomain.RBACObjectInterface, effect libdomain.RBACEffectInterface) error {
	e := r.initEnforcer(ctx)

	if _, err := e.RemoveNamedPolicy("p", subject.Subject(), object.Object(), action.Action(), effect.Effect(), domain.Domain()); err != nil {
		return fmt.Errorf("e.RemoveNamedPolicy: %w", err)
	}

	return nil
}

func (r *RBACRepository) CreateSubjectGroupingPolicy(ctx context.Context, domain libdomain.RBACDomainInterface, child libdomain.RBACSubjectInterface, parent libdomain.RBACSubjectInterface) error {
	e := r.initEnforcer(ctx)

	if _, err := e.AddNamedGroupingPolicy("g", child.Subject(), parent.Subject(), domain.Domain()); err != nil {
		return fmt.Errorf("e.AddNamedGroupingPolicy: %w", err)
	}

	return nil
}

func (r *RBACRepository) DeleteSubjectGroupingPolicy(ctx context.Context, domain libdomain.RBACDomainInterface, child libdomain.RBACSubjectInterface, parent libdomain.RBACSubjectInterface) error {
	e := r.initEnforcer(ctx)

	if _, err := e.RemoveNamedGroupingPolicy("g", child.Subject(), parent.Subject(), domain.Domain()); err != nil {
		return fmt.Errorf("e.RemoveNamedGroupingPolicy: %w", err)
	}

	return nil
}

func (r *RBACRepository) CreateObjectGroupingPolicy(ctx context.Context, domain libdomain.RBACDomainInterface, child libdomain.RBACObjectInterface, parent libdomain.RBACObjectInterface) error {
	e := r.initEnforcer(ctx)

	if _, err := e.AddNamedGroupingPolicy("g2", child.Object(), parent.Object(), domain.Domain()); err != nil {
		return fmt.Errorf("e.AddNamedGroupingPolicy: %w", err)
	}

	return nil
}

func (r *RBACRepository) DeleteObjectGroupingPolicy(ctx context.Context, dom libdomain.RBACDomainInterface, child libdomain.RBACObjectInterface, parent libdomain.RBACObjectInterface) error {
	e := r.initEnforcer(ctx)

	if _, err := e.RemoveNamedGroupingPolicy("g2", child.Object(), parent.Object(), dom.Domain()); err != nil {
		return fmt.Errorf("e.RemoveNamedGroupingPolicy: %w", err)
	}

	return nil
}

func (r *RBACRepository) NewEnforcerWithGroupsAndUsers(_ context.Context, roles []libdomain.RBACRoleInterface, users []libdomain.RBACUserInterface) (*casbin.Enforcer, error) {
	_ = roles
	_ = users
	if err := r.enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("e.LoadPolicy: %w", err)
	}
	return r.enforcer, nil
}

func (r *RBACRepository) GetEnforcer() *casbin.Enforcer {
	return r.enforcer
}

func (r *RBACRepository) GetGroupsForSubject(ctx context.Context, dom libdomain.RBACDomainInterface, subject libdomain.RBACSubjectInterface) ([]libdomain.RBACRole, error) {
	e := r.initEnforcer(ctx)

	if err := e.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("e.LoadPolicy: %w", err)
	}

	roles, err := e.GetImplicitRolesForUser(subject.Subject(), dom.Domain())
	if err != nil {
		return nil, fmt.Errorf("e.GetImplicitRolesForUser: %w", err)
	}

	result := make([]libdomain.RBACRole, len(roles))
	for i, role := range roles {
		result[i] = *libdomain.NewRBACRole(role)
	}

	return result, nil
}
