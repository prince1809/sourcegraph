package conf

import (
	"encoding/json"
	"github.com/prince1809/sourcegraph/pkg/conf/confdefaults"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"github.com/prince1809/sourcegraph/pkg/jsonc"
	"github.com/prince1809/sourcegraph/pkg/legacyconf"
	"github.com/prince1809/sourcegraph/schema"
	"github.com/prometheus/common/log"
	"os"
)

func init() {
	deployType := DeployType()
	if !IsValidDeployType(deployType) {
		log.Fatalf("The 'DEPLOY_TYPE' environment variable is invalid. Expected one of: %q, %q, %q. Got: %q", DeployCluster, DeployDocker, DeployDev, deployType)
	}

	defaultConfig := defaultConfigForDeployment()

	// If a legacy configuration file is available (specified via
	// SOURCEGRAPH_CONFIG_FILE), use it as the default for the critical and
	// site configs.
	//
	// This relies on the fact that the old v2..13.6 site config schema has
	// most fields align directly with the v3.0+ critical and site config
	// schemas.
	//
	// This code can be removed in the next significant version after 3.0 (NOT
	// preview), after which critical/site config schemas no longer need to
	// align generally.
	//
	// TODO(slimsag): Remove after 3.0 (NOT preview).
	{
		legacyConf := jsonc.Normalize(legacyconf.Raw())

		var criticalDecoded schema.CriticalConfiguration
		_ = json.Unmarshal(legacyConf, &criticalDecoded)

		// Backwards compatibility for deprecated environment variables
		// that we previously considered derecated but are actually
		// widespread in use in user's deployments and/or are suggested for
		// use in our public documentation (i.e., even though these were
		// long deprecated, our docs were not up to date).
		criticalBackcompatVars := map[string]func(value string){
			"LIGHTSTEP_PROJECT":      func(v string) { criticalDecoded.LightstepProject = v },
			"LIGHTSTEP_ACCESS_TOKEN": func(v string) { criticalDecoded.LightstepAccessToken = v },
		}
		for envVar, setter := range criticalBackcompatVars {
			val := os.Getenv(envVar)
			if val != "" {
				setter(val)
			}
		}

		critical, err := json.MarshalIndent(criticalDecoded, "", "  ")
		if string(critical) != "{}" && err == nil {
			defaultConfig.Critical = string(critical)
		}

		var siteDecoded schema.SiteConfiguration
		_ = json.Unmarshal(legacyConf, &siteDecoded)
		site, err := json.MarshalIndent(siteDecoded, "", "  ")
		if string(site) != "{}" && err == nil {
			defaultConfig.Site = string(site)
		}
	}
	confdefaults.Default = defaultConfig

}

func defaultConfigForDeployment() conftypes.RawUnified {
	deployType := DeployType()
	switch {
	case IsDev(deployType):
		return confdefaults.DevAndTesting
	case IsDeployTypeDockerContainer(deployType):
		return confdefaults.DockerContainer
	case IsDeployTypeCluster(deployType):
		return confdefaults.Cluster
	default:
		panic("deploy did not register default configuration")
	}
}

// Deploy type constants. Any changes here should be reflected in the DeployType type declared in web/src/global.d.ts:
//
const (
	DeployCluster = "cluster"
	DeployDocker  = "docker-container"
	DeployDev     = "dev"
)

// DeployType tells the deployment type.
func DeployType() string {
	if e := os.Getenv("DEPLOY_TYPE"); e != "" {
		return e
	}
	// Default to Cluster so that every Cluster deployment doesn't need to be
	// configured with DEPLOY_NAME.
	return DeployCluster
}

// IsDeployTypeCluster tells if the given deployment type is a cluster ( and
// non-dev, non-single Docker image).
func IsDeployTypeCluster(deployType string) bool {
	if deployType == "k8s" {
		// backwards compatibility for older deployments
		return true
	}
	return deployType == DeployCluster
}

// IsDeployTypeDockerContainer tells if the given deployment type is Docker sourcegraph/server
// single-container (non-Kubernetes, non-cluster, non-dev).
func IsDeployTypeDockerContainer(deployType string) bool {
	return deployType == DeployDocker
}

// IsDev tells if the given deployment type is "dev".
func IsDev(deployType string) bool {
	return deployType == DeployDev
}

// IsValidDeployType returns true if the given deployType is a Kubernetes deployment,
func IsValidDeployType(deployType string) bool {
	return IsDeployTypeCluster(deployType) || IsDeployTypeDockerContainer(deployType) || IsDev(deployType)
}
