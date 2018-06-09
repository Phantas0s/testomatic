package cmd

import (
	"fmt"
	"testing"
)

func ErrorTest(expected string, got string) string {
	return fmt.Sprintf("Expected '%s' got '%s'", expected, got)
}

func TestCreateRelative(t *testing.T) {
	path := "/home/src/tests/tests.go"
	filepath := "src/"

	expected := "src/tests/tests.go"
	if v, _ := CreateRelative(path, filepath, "current"); *v != expected {
		t.Errorf(ErrorTest(expected, *v))
	}

	expected = "src/tests"
	if v, _ := CreateRelative(path, filepath, "dir"); *v != expected {
		t.Errorf(ErrorTest(expected, *v))
	}

	expected = "src/"
	if v, _ := CreateRelative(path, filepath, "all"); *v != expected {
		t.Errorf(ErrorTest(expected, *v))
	}
}
