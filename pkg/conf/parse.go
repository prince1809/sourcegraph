package conf

import (
	"encoding/json"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"github.com/prince1809/sourcegraph/pkg/jsonc"
	"github.com/prince1809/sourcegraph/schema"
)

// parseConfigData parses the provided config string into the given struct
// pointer.
func parseConfigData(data string, cfg interface{}) error {
	if data != "" {
		data, err := jsonc.Parse(data)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, cfg); err != nil {
			return err
		}
	}

	if v, ok := cfg.(*schema.SiteConfiguration); ok {
		// For convenience, make sire this is not nil.
		if v.ExperimentalFeatures == nil {
			v.ExperimentalFeatures = &schema.ExperimentalFeatures{}
		}
	}
	return nil
}

// ParseConfig parses the raw configuration
func ParseConfig(data conftypes.RawUnified) (*Unified, error) {
	cfg := &Unified{
		ServiceConnections: data.ServiceConnections,
	}
	if err := parseConfigData(data.Critical, &cfg.Critical); err != nil {
		return nil, err
	}
	if err := parseConfigData(data.Site, &cfg.SiteConfiguration); err != nil {
		return nil, err
	}
	return cfg, nil
}
