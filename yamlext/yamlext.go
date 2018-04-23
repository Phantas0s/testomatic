package yamlext

import (
	"github.com/kylelemons/go-gypsy/yaml"
	"log"
)

func ToMap(node yaml.Node) yaml.Map {
	m, ok := node.(yaml.Map)
	if !ok {
		log.Fatalf("%v is not of type map", node)
	}
	return m
}

func ToList(node yaml.Node) yaml.List {
	m, ok := node.(yaml.List)
	if !ok {
		log.Fatalf("%v is not of type list", node)
	}

	return m
}

func ToScalar(node yaml.Node) yaml.Scalar {
	m, ok := node.(yaml.Scalar)
	if !ok {
		log.Fatalf("%v is not of type scalar", node)
	}

	return m
}
