package conf

import (
	"encoding/json"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
)

var contributedValidators []func(Unified) []string

func validateCustomRaw(normalizedInput conftypes.RawUnified) (problems []string, err error) {
	var cfg Unified
	if err := json.Unmarshal([]byte(normalizedInput.Critical), &cfg.Critical); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(normalizedInput.Site), &cfg.Critical); err != nil {
		return nil, err
	}
	return validateCustom(cfg), nil
}

// validateCustom validates the site config using custom validation steps that are not
// able to be expressed in the JSON schema.
func validateCustom(cfg Unified) (problems []string) {

	invalid := func(msg string) {
		problems = append(problems, msg)
	}

	{
		hasSMTP := cfg.EmailSmtp != nil
		hasSMTPAuth := cfg.EmailSmtp != nil && cfg.EmailSmtp.Authentication != "none"
		if hasSMTP && cfg.EmailAddress == "" {
			invalid(`should set email.address because email.smtp is set`)
		}
		if hasSMTPAuth && (cfg.EmailSmtp.Username == "" && cfg.EmailSmtp.Password == "") {
			invalid(`must set email.smtp username and password for email.smtp authentication`)
		}
	}

	for _, f := range contributedValidators {
		problems = append(problems, f(cfg)...)
	}
	return problems
}
