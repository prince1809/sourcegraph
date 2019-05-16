package shared

import (
	"fmt"
	"github.com/prince1809/sourcegraph/pkg/env"
)

func Main() {
	env.Lock()
	fmt.Println("management-console started")
}
