package confdefaults

import "github.com/prince1809/sourcegraph/pkg/conf/conftypes"

// Default is the default for *this* deployment type. It is populated by
// pkg/conf at init time.
//
// In the case of a migration from an old Sourcegraph version to 3.0, this is
// not strictly one of the declared defaults in this package but rather may be
// defaults from a user's old configuration.

var Default conftypes.RawUnified
