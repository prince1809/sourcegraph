package jsonc

import (
	"encoding/json"
	"fmt"
	"github.com/sourcegraph/jsonx"
	"strings"
)

// Unmarshal unmarshals the JSON using a fault-tolerant parser that allows comments and trailing
// commas. In any unrecoverable faults are found, an error is returned.
func Unmarshal(text string, v interface{}) error {
	data, err := Parse(text)
	if err != nil {
		return err
	}
	if strings.TrimSpace(text) == "" {
		return nil
	}
	return json.Unmarshal(data, v)
}

// Parse converts JSON with comments, trailing commas, and some types
func Parse(text string) ([]byte, error) {
	data, errs := jsonx.Parse(text, jsonx.ParseOptions{Comments: true, TrailingCommas: true})
	if len(errs) > 0 {
		return data, fmt.Errorf("failed to parse JSON: %v", errs)
	}
	return data, nil
}

// Normalize is like Parse, except it ignores errors and always returns valid JSON, even if that
// JSON is a subset of the input.
func Normalize(input string) []byte {
	output, _ := jsonx.Parse(string(input), jsonx.ParseOptions{Comments: true, TrailingCommas: true})
	if len(output) == 0 {
		return []byte("{}")
	}
	return output
}
