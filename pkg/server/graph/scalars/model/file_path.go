package model

import (
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"io"
	"os"
)

// if the type referenced in .gqlgen.yml is a function that returns a marshaller we can use it to encode and decode
// onto any existing go type.
func MarshalFilePath(f string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, f)
	})
}

// Unmarshal{Typename} is only required if the scalars appears as an input. The raw values have already been decoded
// from json into int/float64/bool/nil/map[string]interface/[]interface
func UnmarshalFilePath(v interface{}) (string, error) {
	if tmpStr, ok := v.(string); ok {
		if len(tmpStr) > 0 {
			info, err := os.Stat(tmpStr)
			if os.IsNotExist(err) {
				return "", fmt.Errorf("file path '%s' do not exist", tmpStr)
			}

			if info.IsDir() {
				return "", fmt.Errorf("file path '%s' is a directory, must be a file", tmpStr)
			}

			return tmpStr, nil
		}
		return "", nil
	}

	return "", errors.New("file path must be a string")
}
