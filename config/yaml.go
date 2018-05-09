package config

import (
	"gopkg.in/yaml.v2"
)

type (
	watch struct {
		Root         string
		Regex        string
		IgnoreHidden bool `yaml:"ignore_hidden"`
		Abs          bool
	}

	command struct {
		Bin        string
		IgnorePath bool `yaml:"ignore_path"`
		Options    []string
	}

	notification struct {
		ImgFailure    string `yaml:"img_failure"`
		ImgSuccess    string `yaml:"img_success"`
		RegexSuccess  string `yaml:"regex_success"`
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
