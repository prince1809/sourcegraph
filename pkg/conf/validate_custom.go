package conf

import (
	"encoding/json"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
)

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
