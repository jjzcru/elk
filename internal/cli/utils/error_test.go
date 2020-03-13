package utils

import (
	"errors"
	"testing"
)

func TestPrintError(t *testing.T) {
	err := errors.New("test error")
	PrintError(err)
}