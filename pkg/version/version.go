package version

const devVersion = "dev" // Version string for unreleased development builds

// Version is configgured at build time via ldflags like this:
// -ldflags "-X github.com/sourcegraph/sourcegraph/pkg/version.version=1.2.3"
var version = devVersion

// Version returns the version string configured at build time.
func Version() string {
	return version
}

// IsDev reports whether the version string is an unreleased development build.
func IsDev(version string) bool {
	return version == devVersion
}

// Mock is used by tests to mocks the results of Version and IsDev.
func Mock(mockVersion string) {
	version = mockVersion
}
