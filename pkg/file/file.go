package file

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetEnvFromFile returns an map with the env variable from a file
func GetEnvFromFile(filePath string) (map[string]string, error) {
	env := make(map[string]string)
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("log path '%s' is a directory", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		parts := strings.SplitAfterN(scanner.Text(), "=", 2)
		key := strings.ReplaceAll(parts[0], "=", "")
		value := parts[1]
		env[key] = value
	}

	return env, nil
}
