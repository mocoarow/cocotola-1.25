package gateway

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/mocoarow/cocotola-1.25/moonbeam/user/domain"
	"github.com/mocoarow/cocotola-1.25/moonbeam/user/service"
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

type rbacRepository struct {
	db       *gorm.DB
	conf     string
	enforcer casbin.IEnforcer
}

var _ service.RBACRepository = (*rbacRepository)(nil)

func NewRBACRepository(_ context.Context, db *gorm.DB) (service.RBACRepository, error) {
	if db == nil {
		panic(errors.New("db is nil"))
	}

	a, err := gormadapter.NewAdapterByDB(db)
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

	return &rbacRepository{
		db:       db,
		conf:     rbacConf,
		enforcer: e,
	}, nil
}

func (r *rbacRepository) initEnforcer(_ context.Context) casbin.IEnforcer {
	return r.enforcer
}

func (r *rbacRepository) CreatePolicy(ctx context.Context, domain domain.RBACDomain, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error {
	e := r.initEnforcer(ctx)

	if _, err := e.AddNamedPolicy("p", subject.Subject(), object.Object(), action.Action(), effect.Effect(), domain.Domain()); err != nil {
		return fmt.Errorf("e.AddNamedPolicy: %w", err)
	}

	return nil
}

func (r *rbacRepository) DeletePolicy(ctx context.Context, domain domain.RBACDomain, subject domain.RBACSubject, action domain.RBACAction, object domain.RBACObject, effect domain.RBACEffect) error {
	e := r.initEnforcer(ctx)

	if _, err := e.RemoveNamedPolicy("p", subject.Subject(), object.Object(), action.Action(), effect.Effect(), domain.Domain()); err != nil {
		return fmt.Errorf("e.RemoveNamedPolicy: %w", err)
	}

	return nil
}

func (r *rbacRepository) CreateSubjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, child domain.RBACSubject, parent domain.RBACSubject) error {
	e := r.initEnforcer(ctx)

	if _, err := e.AddNamedGroupingPolicy("g", child.Subject(), parent.Subject(), domain.Domain()); err != nil {
		return fmt.Errorf("e.AddNamedGroupingPolicy: %w", err)
	}

	return nil
}

func (r *rbacRepository) DeleteSubjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, child domain.RBACSubject, parent domain.RBACSubject) error {
	e := r.initEnforcer(ctx)

	if _, err := e.RemoveNamedGroupingPolicy("g", child.Subject(), parent.Subject(), domain.Domain()); err != nil {
		return fmt.Errorf("e.RemoveNamedGroupingPolicy: %w", err)
	}

	return nil
}

func (r *rbacRepository) CreateObjectGroupingPolicy(ctx context.Context, domain domain.RBACDomain, child domain.RBACObject, parent domain.RBACObject) error {
	e := r.initEnforcer(ctx)

	if _, err := e.AddNamedGroupingPolicy("g2", child.Object(), parent.Object(), domain.Domain()); err != nil {
		return fmt.Errorf("e.AddNamedGroupingPolicy: %w", err)
	}

	return nil
}

func (r *rbacRepository) DeleteObjectGroupingPolicy(ctx context.Context, dom domain.RBACDomain, child domain.RBACObject, parent domain.RBACObject) error {
	e := r.initEnforcer(ctx)

	if _, err := e.RemoveNamedGroupingPolicy("g2", child.Object(), parent.Object(), dom.Domain()); err != nil {
		return fmt.Errorf("e.RemoveNamedGroupingPolicy: %w", err)
	}

	return nil
}

func (r *rbacRepository) NewEnforcerWithGroupsAndUsers(_ context.Context, roles []domain.RBACRole, users []domain.RBACUser) (service.RBACEnforcer, error) {
	_ = roles
	_ = users
	if err := r.enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("e.LoadPolicy: %w", err)
	}
	return r.enforcer, nil
}

func (r *rbacRepository) GetEnforcer() service.RBACEnforcer {
	return r.enforcer
}

func (r *rbacRepository) GetGroupsForSubject(ctx context.Context, dom domain.RBACDomain, subject domain.RBACSubject) ([]domain.RBACRole, error) {
	e := r.initEnforcer(ctx)

	if err := e.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("e.LoadPolicy: %w", err)
	}

	roles, err := e.GetImplicitRolesForUser(subject.Subject(), dom.Domain())
	if err != nil {
		return nil, fmt.Errorf("e.GetImplicitRolesForUser: %w", err)
	}

	result := make([]domain.RBACRole, len(roles))
	for i, role := range roles {
		result[i] = domain.NewRBACRole(role)
	}

	return result, nil
}
