package config

import (
	"gopkg.in/yaml.v2"
)

type (
	watch struct {
		Root         string
		Regex        string
		IgnoreHidden bool `yaml:"ignore_hidden"`
		Ignore       []string
	}

	command struct {
		Bin        string
		Scope      string
		Abs        bool
		IgnorePath bool `yaml:"ignore_path"`
		Options    []string
	}

	notification struct {
		Disable       bool
		ImgFailure    string `yaml:"img_failure"`
		ImgSuccess    string `yaml:"img_success"`
		Mute          bool
		RegexSuccess  string `yaml:"regex_success"`
		RegexFailure  string `yaml:"regex_failure"`
		DisplayResult bool   `yaml:"display_result"`
	}

	YamlConf struct {
		Watch        watch
		Command      command
		Notification notification
	}
)

func (c *YamlConf) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}
