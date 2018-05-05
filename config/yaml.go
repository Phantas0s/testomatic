package config

import (
	"gopkg.in/yaml.v2"
)

type (
	watch struct {
		Folder       string
		Regex        string
		Abs          bool
		IgnoreHidden bool `yaml:"ignore_hidden"`
		IgnorePath   bool `yaml:"ignore_path"`
	}

	command struct {
		Path    string
		Options []string
	}

	notification struct {
		ImgFailure    string `yaml:"img_failure"`
		ImgSuccess    string `yaml:"img_success"`
		SuccessRegex  string `yaml:"success"`
		DisplayResult bool   `yaml:"display_result"`
	}

	YamlConf struct {
		Watch        watch
		Test         string
		Command      command
		Notification notification
	}
)

func (c *YamlConf) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}
