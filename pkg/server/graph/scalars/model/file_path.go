package model

import (
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"io"
	"os"
)

// MarshalFilePath marshal the file path
func MarshalFilePath(f string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, f)
	})
}

// UnmarshalFilePath unmarshal the file path
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
