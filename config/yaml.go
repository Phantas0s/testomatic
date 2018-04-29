package config

import (
	"gopkg.in/yaml.v2"
)

type (
	watch struct {
		Folder string
		Ext    []string
		Abs    bool
	}

	command struct {
		Path    string
		Options []string
	}

	notification struct {
		ImgFailure   string `yaml:"img_failure"`
		ImgSuccess   string `yaml:"img_success"`
		SuccessRegex string `yaml:"success"`
		DisplayRes   string `yaml:"display_res"`
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
