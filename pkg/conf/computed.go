package conf

import (
	"github.com/prince1809/sourcegraph/pkg/conf/confdefaults"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"github.com/prometheus/common/log"
	"os"
)

func init() {
	deployType := DeployType()
	if !IsValidDeployType(deployType) {
		log.Fatalf("The 'DEPLOY_TYPE' environment variable is invalid. Expected one of: %q, %q, %q. Got: %q", DeployCluster, DeployDocker, DeployDev, deployType)
	}

	defaultConfig := defaultConfigForDeployment()
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
