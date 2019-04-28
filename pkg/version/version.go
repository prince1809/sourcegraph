package version

const devVersion = "dev" // Version string for unreleased development builds

// Version is configgured at build time via ldflags like this:
// -ldflags "-X github.com/sourcegraph/sourcegraph/pkg/version.version=1.2.3"
