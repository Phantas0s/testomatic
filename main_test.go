package main

import (
	"fmt"
	"github.com/Phantas0s/yamlext"
	"github.com/kylelemons/go-gypsy/yaml"
	"testing"
)

var config = yaml.ConfigFile(".testomatic_test.yml")
var root = yamlext.ToMap(config.Root)

func ErrorTest(expected string, got string) string {
	return fmt.Sprintf("Expected '%s' got '%s'", expected, got)
}

func TestIsAbsolute(t *testing.T) {
	if v := IsAbsolute("/home/user/lala"); v != true {
		t.Errorf(ErrorTest("true", fmt.Sprintf("%v", v)))
	}

	if v := IsAbsolute("home/user/lala"); v != false {
		t.Errorf(ErrorTest("false", fmt.Sprintf("%v", v)))
	}
}

func TestCreateRelative(t *testing.T) {
	path := "/home/src/lala.go"
	filepath := "src/"

	if v := CreateRelative(path, filepath); v != "src/lala.go" {
		t.Errorf(ErrorTest("src/lala.go", v))
	}
}

func TestExtractScalar(t *testing.T) {
	if v := ExtractScalar(config, "testomatic.folder"); v != "src/Tests" {
		t.Errorf(ErrorTest("src/Tests", v))
	}
}

func TestOpt(t *testing.T) {
	v := ExtractOpt(root)

	if v[0] != "exec" {
		t.Errorf(ErrorTest("exec", v[0]))
	}

	if v[1] != "-T" {
		t.Errorf(ErrorTest("-T", v[1]))
	}

	if v[2] != "php" {
		t.Errorf(ErrorTest("php", v[2]))
	}
}

func TestExtractExt(t *testing.T) {
	v := ExtractExt(root)

	if v[0] != ".php" {
		t.Errorf(ErrorTest(".php", v[0]))
	}

	if v[1] != ".go" {
		t.Errorf(ErrorTest(".go", v[1]))
	}
}
