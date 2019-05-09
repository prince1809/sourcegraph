package jsonc

import (
	"fmt"
	"github.com/sourcegraph/jsonx"
)

// Parse converts JSON with comments, trailing commas, and some types
func Parse(text string) ([]byte, error) {
	data, errs := jsonx.Parse(text, jsonx.ParseOptions{Comments: true, TrailingCommas: true})
	if len(errs) > 0 {
		return data, fmt.Errorf("failed to parse JSON: %v", errs)
	}
	return data, nil
}
