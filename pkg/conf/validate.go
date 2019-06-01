package conf

import (
	"encoding/json"
	"fmt"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"github.com/sourcegraph/sourcegraph/pkg/jsonc"
	"github.com/sourcegraph/sourcegraph/schema"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

// ignoreLegacyKubernetesFields is the set of field names for which validation occurs should be
// ignored. The validation occur only because deploy-sourcegraph config merged site config
// and Kubernetes cluster-specific config. This is deprecated. Until we have transitioned fully, we
// suppress validation errors on these fields.
var ignoreLegacyKubernetesFields = map[string]struct{}{
	"alertmanagerConfig": {},
}
// Validate validates the configuration against the JSON Schema and other
// custom validation checks.
func Validate(input conftypes.RawUnified) (problems []string, err error) {
	criticalProblems, err := doValidate(input.Critical, schema.CriticalSchemaJSON)
	if err != nil {
		return nil, err
	}
	problems = append(problems, criticalProblems...)

	siteProblems, err := doValidate(input.Site, schema.SiteSchemaJSON)
	if err != nil {
		return nil, err
	}
	problems = append(problems, siteProblems...)

	customProblems, err := validateCustomRaw(conftypes.RawUnified{
		Critical: string(jsonc.Normalize(input.Critical)),
		Site:     string(jsonc.Normalize(input.Site)),
	})
	if err != nil {
		return nil, err
	}
	problems = append(problems, customProblems...)
	return problems, nil
}

func doValidate(inputStr, schema string) (problems []string, err error) {
	input := []byte(jsonc.Normalize(inputStr))

	res, err := validate([]byte(schema), input)
	if err != nil {
		return nil, err
	}
	problems = make([]string, 0, len(res.Errors()))
	for _, e := range res.Errors() {
		if _, ok := ignoreLegacyKubernetesFields[e.Field()]; ok {
			continue
		}

		var keyPath string
		if c := e.Context(); c != nil {
			keyPath = strings.TrimPrefix(e.Context().String("."), "(root).")
		} else {
			keyPath = e.Field()
		}
		problems = append(problems, fmt.Sprintf("%s: %s", keyPath, e.Description()))
	}
	return problems, nil
}

func validate(schema, input []byte) (*gojsonschema.Result, error) {
	if len(input) > 0 {
		var v map[string]interface{}
		if err := json.Unmarshal(input, &v); err != nil {
			return nil, err
		}
		delete(v, "settings")
		var err error
		input, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}

	s, err := gojsonschema.NewSchema(jsonLoader{gojsonschema.NewBytesLoader(schema)})
	if err != nil {
		return nil, err
	}
	return s.Validate(gojsonschema.NewBytesLoader(input))
}

type jsonLoader struct {
	gojsonschema.JSONLoader
}

// MustValidateDefaults should be called after all custom validators have been
// registered. It will panic if any of the default deployment configurations
// are invalid.
func MustValidate() {

}

// mustValidate panics if the configuration does not pass validation.
func mustValidate(name string, cfg conftypes.RawUnified) conftypes.RawUnified {
	problems, err := Validate
}
