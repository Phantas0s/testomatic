package config

type Config interface {
	Parse(data []byte) error
}
