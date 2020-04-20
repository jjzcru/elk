package model

import (
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"io"
	"time"
)

// if the type referenced in .gqlgen.yml is a function that returns a marshaller we can use it to encode and decode
// onto any existing go type.
func MarshalDuration(d time.Duration) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, fmt.Sprintf("\"%s\"", d.String()))
	})
}

// Unmarshal{Typename} is only required if the scalars appears as an input. The raw values have already been decoded
// from json into int/float64/bool/nil/map[string]interface/[]interface
func UnmarshalDuration(v interface{}) (time.Duration, error) {
	if tmpStr, ok := v.(string); ok {
		duration, err := time.ParseDuration(tmpStr)
		if err != nil {
			return 0, errors.New("invalid duration format")
		}

		return duration, nil
	}

	return 0, errors.New("duration needs to be a string")
}
