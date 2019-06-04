package db

import (
	"context"
	"fmt"
	"github.com/keegancsmith/sqlf"
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/cmd/frontend/types"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"github.com/prince1809/sourcegraph/pkg/db/dbconn"
	"github.com/prince1809/sourcegraph/schema"
	"github.com/xeipuuv/gojsonschema"
)

// An ExternalServicesStore stores external services and their configuration.
// Before updating or creating a new external service, validation is performed.
// The enterprise code registers additional validators at run-time and sets the
// global instance in stores.go
type ExternalServicesStore struct {
	GithubValidators []func(*schema.GitHubConnection) error
	GitLabValidators []func(*schema.GitLabConnection, []schema.AuthProviders)
}

// ExternalServiceKinds contains a map of all supported kinds of
// external services.
var ExternalServiceKinds = map[string]ExternalServiceKind{

}

// ExternalServiceKind describes a kind of external service.
type ExternalServiceKind struct {
	// True if the external service can host repositories.
	CodeHost bool

	JSONSchema string // JSON schema for the external service's configuration.
}

// ExternalServicesListOptions contains options for listing external services.
type ExternalServicesListOptions struct {
	Kinds []string
	*LimitOffset
}

func (o ExternalServicesListOptions) sqlConditions() []*sqlf.Query {
	conds := []*sqlf.Query{sqlf.Sprintf("deleted_at IS NULL")}
	if len(o.Kinds) > 0 {
		kinds := []*sqlf.Query{}
		for _, kind := range o.Kinds {
			kinds = append(kinds, sqlf.Sprintf("%s", kind))
		}
		conds = append(conds, sqlf.Sprintf("kind IN (%s)", sqlf.Join(kinds, ", ")))
	}
	return conds
}

// ValidateConfig validates the given external service configuration.
func (e *ExternalServicesStore) ValidateConfig(kind, config string, ps []schema.AuthProviders) error {
	ext, ok := ExternalServiceKinds[kind]
	if !ok {
		return fmt.Errorf("invalid external service kind: %s", kind)
	}

	// All configs must be valid JSON.
	// If this requirement is ever changed, you will need to update
	// serveExternalServiceConfigs to handle this case.
	sl := gojsonschema.NewSchemaLoader()
	sc, err := sl.Compile(gojsonschema.NewStringLoader(ext.JSONSchema))
	if err != nil {
		return errors.Wrapf(err, "failed to compile schema for external service of kind %q", kind)
	}
}
func

// Create creates a external service.
//
// 🚨 SECURITY: The caller must ensure that the actor is a site admin.
func (c *ExternalServicesStore) Create(ctx context.Context, externalService *types.ExternalService) error {
	ps := conf.Get().Critical.AuthProviders
	if err := c.ValidateConfig(externalService.Kind)
}

// List returns all external services.
//
// 🚨 SECURITY: The caller must ensure that the actor is a site admin.
func (c *ExternalServicesStore) List(ctx context.Context, opt ExternalServicesListOptions) ([]*types.ExternalService, error) {
	if Mocks.ExternalServices.List != nil {
		return Mocks.ExternalServices.List(opt)
	}
	return c.list(ctx, opt.sqlConditions(), opt.LimitOffset)
}

func (c *ExternalServicesStore) list(ctx context.Context, conds []*sqlf.Query, limitOffset *LimitOffset) ([]*types.ExternalService, error) {
	q := sqlf.Sprintf(`
		SELECT id, kind, display_name, config, created_at, updated_at
		FROM external_services
		WHERE (%s)
		ORDER BY id DESC
		%s`,
		sqlf.Join(conds, ") AND ("),
		limitOffset.SQL(),
	)

	rows, err := dbconn.Global.QueryContext(ctx, q.Query(sqlf.PostgresBindVar), q.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*types.ExternalService
	for rows.Next() {
		var h types.ExternalService
		if err := rows.Scan(&h.ID, &h.Kind, &h.DisplayName, &h.Config, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, &h)
	}
	return results, nil
}

// MockExternalServices mocks the external services store.
type MockExternalServices struct {
	GetByID func(id int64) (*types.ExternalService, error)
	List    func(opt ExternalServicesListOptions) ([]*types.ExternalService, error)
}
